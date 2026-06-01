package model

import (
	"time"

	"github.com/lib/pq"
	"github.com/shopspring/decimal"
)

type Category struct {
	ID        int    `json:"id" gorm:"primaryKey"`
	Name      string `json:"name" gorm:"size:50;not null"`
	ParentID  *int   `json:"parent_id"`
	Icon      string `json:"icon" gorm:"size:50"`
	SortOrder int    `json:"sort_order" gorm:"default:0"`
}

func (Category) TableName() string {
	return "categories"
}

type Item struct {
	ID             int64           `json:"id" gorm:"primaryKey"`
	UserID         int64           `json:"user_id" gorm:"index;not null"`
	Title          string          `json:"title" gorm:"size:100;not null"`
	Description    string          `json:"description"`
	CategoryID     *int            `json:"category_id" gorm:"index"`
	EstimatedValue decimal.Decimal `json:"estimated_value" gorm:"type:decimal(10,2)" swaggertype:"string"`
	Condition      string          `json:"condition" gorm:"size:20;not null;default:good"`
	Images         pq.StringArray  `json:"images" gorm:"type:text[]" swaggertype:"array,string"`
	VideoURL       string          `json:"video_url" gorm:"size:255"`
	Status         string          `json:"status" gorm:"size:20;not null;default:active;index"`
	Location       string          `json:"location" gorm:"size:100"`
	ViewCount      int             `json:"view_count" gorm:"not null;default:0"`
	WantItems      pq.StringArray  `json:"want_items" gorm:"type:text[]" swaggertype:"array,string"`
	CreatedAt      time.Time       `json:"created_at"`
	UpdatedAt      time.Time       `json:"updated_at"`

	// 关联
	User     *User     `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Category *Category `json:"category,omitempty" gorm:"foreignKey:CategoryID"`
}

func (Item) TableName() string {
	return "items"
}

// 物品状态常量
const (
	ItemStatusActive   = "active"
	ItemStatusInactive = "inactive"
	ItemStatusTraded   = "traded"
	ItemStatusDeleted  = "deleted"
)

// 物品成色常量
const (
	ConditionNew     = "new"
	ConditionLikeNew = "like_new"
	ConditionGood    = "good"
	ConditionFair    = "fair"
)
