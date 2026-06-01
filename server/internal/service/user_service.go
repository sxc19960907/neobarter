package service

import (
	"errors"
	"regexp"

	"github.com/neobarter/server/internal/model"
	"github.com/neobarter/server/internal/repository"
)

type UserService struct {
	userRepo *repository.UserRepository
}

func NewUserService(userRepo *repository.UserRepository) *UserService {
	return &UserService{userRepo: userRepo}
}

func (s *UserService) GetByID(id int64) (*model.User, error) {
	return s.userRepo.GetByID(id)
}

func (s *UserService) Update(user *model.User) error {
	return s.userRepo.Update(user)
}

// 身份证号：18 位，末位可为 X
var idCardRe = regexp.MustCompile(`^\d{17}[\dXx]$`)

// VerifyRealName 提交个人实名认证。
// 注意：MVP 直接置为已认证，未接入三方核验/人工审核。
func (s *UserService) VerifyRealName(userID int64, realName, idCard string) error {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return ErrNotFound
	}
	if user.RealNameVerified {
		return errors.New("已完成实名认证，无需重复提交")
	}
	if realName == "" {
		return errors.New("请填写真实姓名")
	}
	if !idCardRe.MatchString(idCard) {
		return errors.New("身份证号格式不正确")
	}

	user.RealName = realName
	user.IDCard = idCard
	user.RealNameVerified = true
	return s.userRepo.Update(user)
}

// VerifyEnterprise 提交企业认证。
// 注意：MVP 直接置为已认证，未接入工商核验/人工审核。
func (s *UserService) VerifyEnterprise(userID int64, enterpriseName, licenseURL string) error {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return ErrNotFound
	}
	if user.UserType != "enterprise" {
		return errors.New("仅企业用户可提交企业认证")
	}
	if user.EnterpriseVerified {
		return errors.New("已完成企业认证，无需重复提交")
	}
	if enterpriseName == "" || licenseURL == "" {
		return errors.New("请填写企业名称并上传营业执照")
	}

	user.EnterpriseName = enterpriseName
	user.EnterpriseLicenseURL = licenseURL
	user.EnterpriseVerified = true
	return s.userRepo.Update(user)
}

// Address operations

func (s *UserService) ListAddresses(userID int64) ([]model.UserAddress, error) {
	return s.userRepo.ListAddresses(userID)
}

func (s *UserService) CreateAddress(addr *model.UserAddress) error {
	if addr.IsDefault {
		s.userRepo.ClearDefaultAddress(addr.UserID)
	}
	return s.userRepo.CreateAddress(addr)
}

func (s *UserService) UpdateAddress(addr *model.UserAddress) error {
	if addr.IsDefault {
		s.userRepo.ClearDefaultAddress(addr.UserID)
	}
	return s.userRepo.UpdateAddress(addr)
}

func (s *UserService) DeleteAddress(id, userID int64) error {
	return s.userRepo.DeleteAddress(id, userID)
}

func (s *UserService) GetAddress(id, userID int64) (*model.UserAddress, error) {
	return s.userRepo.GetAddress(id, userID)
}
