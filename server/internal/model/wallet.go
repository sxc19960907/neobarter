package model

import (
	"time"

	"github.com/shopspring/decimal"
)

type Wallet struct {
	ID            int64           `json:"id" gorm:"primaryKey"`
	UserID        int64           `json:"user_id" gorm:"uniqueIndex;not null"`
	Balance       decimal.Decimal `json:"balance" gorm:"type:decimal(12,2);not null;default:0" swaggertype:"string"`
	FrozenBalance decimal.Decimal `json:"frozen_balance" gorm:"type:decimal(12,2);not null;default:0" swaggertype:"string"`
	TotalIncome   decimal.Decimal `json:"total_income" gorm:"type:decimal(12,2);not null;default:0" swaggertype:"string"`
	TotalExpense  decimal.Decimal `json:"total_expense" gorm:"type:decimal(12,2);not null;default:0" swaggertype:"string"`
	CreatedAt     time.Time       `json:"created_at"`
	UpdatedAt     time.Time       `json:"updated_at"`
}

func (Wallet) TableName() string {
	return "wallets"
}

// WalletTransaction 钱包流水
type WalletTransaction struct {
	ID            int64           `json:"id" gorm:"primaryKey"`
	WalletID      int64           `json:"wallet_id" gorm:"index;not null"`
	Type          string          `json:"type" gorm:"size:20;not null"`          // reward / trade_in / trade_out / deposit / withdraw / freeze / unfreeze
	Amount        decimal.Decimal `json:"amount" gorm:"type:decimal(12,2);not null" swaggertype:"string"`
	BalanceAfter  decimal.Decimal `json:"balance_after" gorm:"type:decimal(12,2);not null" swaggertype:"string"`
	ReferenceType string          `json:"reference_type" gorm:"size:30"`         // trade_request / system / usdc
	ReferenceID   int64           `json:"reference_id"`
	Description   string          `json:"description" gorm:"size:200"`
	CreatedAt     time.Time       `json:"created_at"`
}

func (WalletTransaction) TableName() string {
	return "wallet_transactions"
}

// 流水类型常量
const (
	TxTypeReward   = "reward"
	TxTypeTradeIn  = "trade_in"
	TxTypeTradeOut = "trade_out"
	TxTypeDeposit  = "deposit"
	TxTypeWithdraw = "withdraw"
	TxTypeFreeze   = "freeze"
	TxTypeUnfreeze = "unfreeze"
)
