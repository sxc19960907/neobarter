package repository

import (
	"github.com/neobarter/server/internal/model"
	"gorm.io/gorm"
)

type ReviewRepository struct {
	db *gorm.DB
}

func NewReviewRepository(db *gorm.DB) *ReviewRepository {
	return &ReviewRepository{db: db}
}

func (r *ReviewRepository) Create(review *model.Review) error {
	return r.db.Create(review).Error
}

func (r *ReviewRepository) GetByTradeAndReviewer(tradeID, reviewerID int64) (*model.Review, error) {
	var review model.Review
	err := r.db.Where("trade_request_id = ? AND reviewer_id = ?", tradeID, reviewerID).First(&review).Error
	if err != nil {
		return nil, err
	}
	return &review, nil
}

func (r *ReviewRepository) ListByReviewee(revieweeID int64, page, pageSize int) ([]model.Review, int64, error) {
	var reviews []model.Review
	var total int64

	query := r.db.Where("reviewee_id = ?", revieweeID)
	query.Model(&model.Review{}).Count(&total)

	err := query.Preload("Reviewer").
		Order("created_at DESC").
		Offset((page - 1) * pageSize).Limit(pageSize).
		Find(&reviews).Error

	return reviews, total, err
}

func (r *ReviewRepository) GetAverageRating(revieweeID int64) (float64, int64, error) {
	var result struct {
		Avg   float64
		Count int64
	}
	err := r.db.Model(&model.Review{}).
		Where("reviewee_id = ?", revieweeID).
		Select("COALESCE(AVG(rating), 0) as avg, COUNT(*) as count").
		Scan(&result).Error
	return result.Avg, result.Count, err
}

type NotificationRepository struct {
	db *gorm.DB
}

func NewNotificationRepository(db *gorm.DB) *NotificationRepository {
	return &NotificationRepository{db: db}
}

func (r *NotificationRepository) Create(n *model.Notification) error {
	return r.db.Create(n).Error
}

func (r *NotificationRepository) List(userID int64, page, pageSize int) ([]model.Notification, int64, error) {
	var notifications []model.Notification
	var total int64

	query := r.db.Where("user_id = ?", userID)
	query.Model(&model.Notification{}).Count(&total)

	err := query.Order("created_at DESC").
		Offset((page - 1) * pageSize).Limit(pageSize).
		Find(&notifications).Error

	return notifications, total, err
}

func (r *NotificationRepository) UnreadCount(userID int64) (int64, error) {
	var count int64
	err := r.db.Model(&model.Notification{}).
		Where("user_id = ? AND is_read = ?", userID, false).
		Count(&count).Error
	return count, err
}

func (r *NotificationRepository) MarkRead(id, userID int64) error {
	return r.db.Model(&model.Notification{}).
		Where("id = ? AND user_id = ?", id, userID).
		Update("is_read", true).Error
}

func (r *NotificationRepository) MarkAllRead(userID int64) error {
	return r.db.Model(&model.Notification{}).
		Where("user_id = ? AND is_read = ?", userID, false).
		Update("is_read", true).Error
}
