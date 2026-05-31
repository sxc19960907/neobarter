package repository

import (
	"github.com/neobarter/server/internal/model"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type WalletRepository struct {
	db *gorm.DB
}

func NewWalletRepository(db *gorm.DB) *WalletRepository {
	return &WalletRepository{db: db}
}

func (r *WalletRepository) Create(wallet *model.Wallet) error {
	return r.db.Create(wallet).Error
}

func (r *WalletRepository) GetByUserID(userID int64) (*model.Wallet, error) {
	var wallet model.Wallet
	err := r.db.Where("user_id = ?", userID).First(&wallet).Error
	if err != nil {
		return nil, err
	}
	return &wallet, nil
}

// GetByUserIDForUpdate 加行锁查询钱包（用于转账等需要原子操作的场景）
func (r *WalletRepository) GetByUserIDForUpdate(tx *gorm.DB, userID int64) (*model.Wallet, error) {
	var wallet model.Wallet
	err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("user_id = ?", userID).First(&wallet).Error
	if err != nil {
		return nil, err
	}
	return &wallet, nil
}

func (r *WalletRepository) UpdateBalance(tx *gorm.DB, walletID int64, balance, frozen decimal.Decimal) error {
	return tx.Model(&model.Wallet{}).Where("id = ?", walletID).
		Updates(map[string]interface{}{
			"balance":        balance,
			"frozen_balance": frozen,
		}).Error
}

func (r *WalletRepository) AddIncome(tx *gorm.DB, walletID int64, amount decimal.Decimal) error {
	return tx.Model(&model.Wallet{}).Where("id = ?", walletID).
		UpdateColumn("total_income", gorm.Expr("total_income + ?", amount)).Error
}

func (r *WalletRepository) AddExpense(tx *gorm.DB, walletID int64, amount decimal.Decimal) error {
	return tx.Model(&model.Wallet{}).Where("id = ?", walletID).
		UpdateColumn("total_expense", gorm.Expr("total_expense + ?", amount)).Error
}

func (r *WalletRepository) CreateTransaction(tx *gorm.DB, txn *model.WalletTransaction) error {
	return tx.Create(txn).Error
}

func (r *WalletRepository) ListTransactions(walletID int64, page, pageSize int) ([]model.WalletTransaction, int64, error) {
	var transactions []model.WalletTransaction
	var total int64

	query := r.db.Where("wallet_id = ?", walletID)
	query.Model(&model.WalletTransaction{}).Count(&total)

	err := query.Order("created_at DESC").
		Offset((page - 1) * pageSize).Limit(pageSize).
		Find(&transactions).Error

	return transactions, total, err
}

func (r *WalletRepository) DB() *gorm.DB {
	return r.db
}
