package repository

import (
	"time"

	"github.com/neobarter/server/internal/model"
	"gorm.io/gorm"
)

type TradeRepository struct {
	db *gorm.DB
}

func NewTradeRepository(db *gorm.DB) *TradeRepository {
	return &TradeRepository{db: db}
}

func (r *TradeRepository) Create(trade *model.TradeRequest) error {
	return r.db.Create(trade).Error
}

func (r *TradeRepository) GetByID(id int64) (*model.TradeRequest, error) {
	var trade model.TradeRequest
	err := r.db.Preload("Initiator").Preload("TargetUser").
		Preload("TargetItem").Preload("OfferedItem").
		First(&trade, id).Error
	if err != nil {
		return nil, err
	}
	return &trade, nil
}

func (r *TradeRepository) Update(trade *model.TradeRequest) error {
	return r.db.Save(trade).Error
}

func (r *TradeRepository) ListByUser(userID int64, status string, page, pageSize int) ([]model.TradeRequest, int64, error) {
	var trades []model.TradeRequest
	var total int64

	query := r.db.Where("initiator_id = ? OR target_user_id = ?", userID, userID)
	if status != "" {
		query = query.Where("status = ?", status)
	}

	query.Model(&model.TradeRequest{}).Count(&total)

	err := query.Preload("Initiator").Preload("TargetUser").
		Preload("TargetItem").Preload("OfferedItem").
		Order("created_at DESC").
		Offset((page-1) * pageSize).Limit(pageSize).
		Find(&trades).Error

	return trades, total, err
}

func (r *TradeRepository) DB() *gorm.DB {
	return r.db
}

// FindExpiredPending 查出已超时但仍为 pending 的交易（用于通知发起方）。
func (r *TradeRepository) FindExpiredPending() ([]model.TradeRequest, error) {
	var trades []model.TradeRequest
	err := r.db.
		Where("status = ? AND expired_at IS NOT NULL AND expired_at < ?",
			model.TradeStatusPending, time.Now()).
		Find(&trades).Error
	return trades, err
}

// ExpirePending 批量把超时的 pending 交易置为 expired，返回受影响行数。
func (r *TradeRepository) ExpirePending() (int64, error) {
	res := r.db.Model(&model.TradeRequest{}).
		Where("status = ? AND expired_at IS NOT NULL AND expired_at < ?",
			model.TradeStatusPending, time.Now()).
		Update("status", model.TradeStatusExpired)
	return res.RowsAffected, res.Error
}
