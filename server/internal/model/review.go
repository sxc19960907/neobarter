package model

import "time"

type Review struct {
	ID             int64     `json:"id" gorm:"primaryKey"`
	TradeRequestID int64     `json:"trade_request_id" gorm:"uniqueIndex:idx_trade_reviewer;not null"`
	ReviewerID     int64     `json:"reviewer_id" gorm:"uniqueIndex:idx_trade_reviewer;not null"`
	RevieweeID     int64     `json:"reviewee_id" gorm:"index;not null"`
	Rating         int       `json:"rating" gorm:"not null"`
	Comment        string    `json:"comment"`
	CreatedAt      time.Time `json:"created_at"`

	Reviewer *User `json:"reviewer,omitempty" gorm:"foreignKey:ReviewerID"`
	Reviewee *User `json:"reviewee,omitempty" gorm:"foreignKey:RevieweeID"`
}

func (Review) TableName() string {
	return "reviews"
}

type Notification struct {
	ID            int64     `json:"id" gorm:"primaryKey"`
	UserID        int64     `json:"user_id" gorm:"index;not null"`
	Type          string    `json:"type" gorm:"size:30;not null"`
	Title         string    `json:"title" gorm:"size:100;not null"`
	Content       string    `json:"content"`
	ReferenceType string    `json:"reference_type" gorm:"size:30"`
	ReferenceID   int64     `json:"reference_id"`
	IsRead        bool      `json:"is_read" gorm:"not null;default:false"`
	CreatedAt     time.Time `json:"created_at"`
}

func (Notification) TableName() string {
	return "notifications"
}

// 通知类型常量
const (
	NotifyTradeRequest  = "trade_request"
	NotifyTradeAccepted = "trade_accepted"
	NotifyTradeRejected = "trade_rejected"
	NotifyMessage       = "message"
	NotifySystem        = "system"
)
