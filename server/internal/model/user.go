package model

import "time"

type User struct {
	ID                   int64      `json:"id" gorm:"primaryKey"`
	Phone                string     `json:"phone" gorm:"uniqueIndex;size:20;not null"`
	Nickname             string     `json:"nickname" gorm:"size:50"`
	AvatarURL            string     `json:"avatar_url" gorm:"size:255"`
	UserType             string     `json:"user_type" gorm:"size:10;not null;default:personal"`
	Status               string     `json:"status" gorm:"size:20;not null;default:active"`
	CreditScore          int        `json:"credit_score" gorm:"not null;default:100"`
	RealName             string     `json:"-" gorm:"size:50"`
	IDCard               string     `json:"-" gorm:"size:30"`
	RealNameVerified     bool       `json:"real_name_verified" gorm:"not null;default:false"`
	EnterpriseName       string     `json:"enterprise_name" gorm:"size:100"`
	EnterpriseLicenseURL string     `json:"-" gorm:"size:255"`
	EnterpriseVerified   bool       `json:"enterprise_verified" gorm:"not null;default:false"`
	Location             string     `json:"location" gorm:"size:100"`
	Bio                  string     `json:"bio"`
	LastLoginAt          *time.Time `json:"last_login_at"`
	CreatedAt            time.Time  `json:"created_at"`
	UpdatedAt            time.Time  `json:"updated_at"`
}

func (User) TableName() string {
	return "users"
}

type UserAddress struct {
	ID        int64     `json:"id" gorm:"primaryKey"`
	UserID    int64     `json:"user_id" gorm:"index;not null"`
	Name      string    `json:"name" gorm:"size:50;not null"`
	Phone     string    `json:"phone" gorm:"size:20;not null"`
	Province  string    `json:"province" gorm:"size:30;not null"`
	City      string    `json:"city" gorm:"size:30;not null"`
	District  string    `json:"district" gorm:"size:30;not null"`
	Detail    string    `json:"detail" gorm:"size:200;not null"`
	IsDefault bool      `json:"is_default" gorm:"not null;default:false"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (UserAddress) TableName() string {
	return "user_addresses"
}
