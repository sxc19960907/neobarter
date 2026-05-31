package service

import (
	"errors"

	"github.com/neobarter/server/internal/model"
	"github.com/neobarter/server/internal/repository"
	"gorm.io/gorm"
)

type ReviewService struct {
	reviewRepo *repository.ReviewRepository
	userRepo   *repository.UserRepository
}

func NewReviewService(reviewRepo *repository.ReviewRepository, userRepo *repository.UserRepository) *ReviewService {
	return &ReviewService{
		reviewRepo: reviewRepo,
		userRepo:   userRepo,
	}
}

// Create 创建评价
func (s *ReviewService) Create(review *model.Review) error {
	// 检查是否已评价
	_, err := s.reviewRepo.GetByTradeAndReviewer(review.TradeRequestID, review.ReviewerID)
	if err == nil {
		return errors.New("已经评价过了")
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	if err := s.reviewRepo.Create(review); err != nil {
		return err
	}

	// 更新被评价者信用积分
	if review.Rating >= 4 {
		s.userRepo.UpdateCreditScore(review.RevieweeID, 2) // 好评 +2
	} else if review.Rating <= 2 {
		s.userRepo.UpdateCreditScore(review.RevieweeID, -2) // 差评 -2
	}

	return nil
}

// ListByUser 获取用户收到的评价
func (s *ReviewService) ListByUser(userID int64, page, pageSize int) ([]model.Review, int64, error) {
	return s.reviewRepo.ListByReviewee(userID, page, pageSize)
}

type NotificationService struct {
	notificationRepo *repository.NotificationRepository
}

func NewNotificationService(notificationRepo *repository.NotificationRepository) *NotificationService {
	return &NotificationService{notificationRepo: notificationRepo}
}

func (s *NotificationService) List(userID int64, page, pageSize int) ([]model.Notification, int64, error) {
	return s.notificationRepo.List(userID, page, pageSize)
}

func (s *NotificationService) UnreadCount(userID int64) (int64, error) {
	return s.notificationRepo.UnreadCount(userID)
}

func (s *NotificationService) MarkRead(id, userID int64) error {
	return s.notificationRepo.MarkRead(id, userID)
}

func (s *NotificationService) MarkAllRead(userID int64) error {
	return s.notificationRepo.MarkAllRead(userID)
}
