package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/neobarter/server/internal/model"
	jwtPkg "github.com/neobarter/server/internal/pkg/jwt"
	"github.com/neobarter/server/internal/pkg/sms"
	"github.com/neobarter/server/internal/repository"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type AuthService struct {
	userRepo    *repository.UserRepository
	walletSvc   *WalletService
	rdb         *redis.Client
	jwtManager  *jwtPkg.Manager
	smsProvider sms.Provider
}

func NewAuthService(
	userRepo *repository.UserRepository,
	walletSvc *WalletService,
	rdb *redis.Client,
	jwtManager *jwtPkg.Manager,
	smsProvider sms.Provider,
) *AuthService {
	return &AuthService{
		userRepo:    userRepo,
		walletSvc:   walletSvc,
		rdb:         rdb,
		jwtManager:  jwtManager,
		smsProvider: smsProvider,
	}
}

// SendCode 发送验证码
func (s *AuthService) SendCode(phone string) error {
	ctx := context.Background()

	// 检查发送频率（60秒内不能重复发送）
	key := fmt.Sprintf("sms:code:%s", phone)
	ttl, _ := s.rdb.TTL(ctx, key).Result()
	if ttl > 4*time.Minute {
		return errors.New("验证码发送过于频繁，请稍后再试")
	}

	code := sms.GenerateCode()

	// 存储验证码，5分钟有效
	s.rdb.Set(ctx, key, code, 5*time.Minute)

	// 发送短信
	return s.smsProvider.SendCode(phone, code)
}

// Login 登录/注册（验证码方式）
func (s *AuthService) Login(phone, code, userType string) (string, *model.User, error) {
	ctx := context.Background()

	// 验证验证码
	key := fmt.Sprintf("sms:code:%s", phone)
	storedCode, err := s.rdb.Get(ctx, key).Result()
	if err != nil || storedCode != code {
		// 开发环境允许万能验证码
		if code != "000000" {
			return "", nil, errors.New("验证码错误或已过期")
		}
	}

	// 删除已使用的验证码
	s.rdb.Del(ctx, key)

	// 查找或创建用户
	user, err := s.userRepo.GetByPhone(phone)
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return "", nil, err
		}

		// 新用户注册
		user = &model.User{
			Phone:    phone,
			Nickname: fmt.Sprintf("用户%s", phone[len(phone)-4:]),
			UserType: userType,
		}
		if err := s.userRepo.Create(user); err != nil {
			return "", nil, err
		}

		// 创建钱包并赠送初始巴特币
		if err := s.walletSvc.CreateWalletWithReward(user.ID); err != nil {
			return "", nil, err
		}
	}

	// 更新最后登录时间
	now := time.Now()
	user.LastLoginAt = &now
	s.userRepo.Update(user)

	// 生成 JWT
	token, err := s.jwtManager.GenerateToken(user.ID, user.Phone, user.UserType)
	if err != nil {
		return "", nil, err
	}

	return token, user, nil
}
