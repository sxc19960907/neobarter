package model

import (
	"time"

	"github.com/shopspring/decimal"
)

type TradeRequest struct {
	ID               int64           `json:"id" gorm:"primaryKey"`
	InitiatorID      int64           `json:"initiator_id" gorm:"index;not null"`
	TargetUserID     int64           `json:"target_user_id" gorm:"index;not null"`
	TargetItemID     int64           `json:"target_item_id" gorm:"not null"`
	OfferedItemID    *int64          `json:"offered_item_id"`
	BarterCoinAmount decimal.Decimal `json:"barter_coin_amount" gorm:"type:decimal(10,2);default:0"`
	Status           string          `json:"status" gorm:"size:20;not null;default:pending;index"`
	Message          string          `json:"message"`
	RejectReason     string          `json:"reject_reason"`
	ExpiredAt        *time.Time      `json:"expired_at"`
	CompletedAt      *time.Time      `json:"completed_at"`
	CreatedAt        time.Time       `json:"created_at"`
	UpdatedAt        time.Time       `json:"updated_at"`

	// 关联
	Initiator   *User `json:"initiator,omitempty" gorm:"foreignKey:InitiatorID"`
	TargetUser  *User `json:"target_user,omitempty" gorm:"foreignKey:TargetUserID"`
	TargetItem  *Item `json:"target_item,omitempty" gorm:"foreignKey:TargetItemID"`
	OfferedItem *Item `json:"offered_item,omitempty" gorm:"foreignKey:OfferedItemID"`
}

func (TradeRequest) TableName() string {
	return "trade_requests"
}

// 交易状态常量
const (
	TradeStatusPending   = "pending"
	TradeStatusAccepted  = "accepted"
	TradeStatusRejected  = "rejected"
	TradeStatusCompleted = "completed"
	TradeStatusCancelled = "cancelled"
	TradeStatusExpired   = "expired"
)
