package service

import (
	"errors"

	"github.com/neobarter/server/internal/model"
	"github.com/neobarter/server/internal/repository"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type WalletService struct {
	walletRepo    *repository.WalletRepository
	initialReward float64
}

func NewWalletService(walletRepo *repository.WalletRepository, initialReward float64) *WalletService {
	return &WalletService{
		walletRepo:    walletRepo,
		initialReward: initialReward,
	}
}

// CreateWalletWithReward 创建钱包并赠送初始巴特币
func (s *WalletService) CreateWalletWithReward(userID int64) error {
	reward := decimal.NewFromFloat(s.initialReward)

	wallet := &model.Wallet{
		UserID:      userID,
		Balance:     reward,
		TotalIncome: reward,
	}
	if err := s.walletRepo.Create(wallet); err != nil {
		return err
	}

	// 记录流水
	tx := s.walletRepo.DB().Begin()
	txn := &model.WalletTransaction{
		WalletID:      wallet.ID,
		Type:          model.TxTypeReward,
		Amount:        reward,
		BalanceAfter:  reward,
		ReferenceType: "system",
		Description:   "注册赠送巴特币",
	}
	if err := s.walletRepo.CreateTransaction(tx, txn); err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

// GetWallet 获取用户钱包
func (s *WalletService) GetWallet(userID int64) (*model.Wallet, error) {
	return s.walletRepo.GetByUserID(userID)
}

// ListTransactions 获取流水列表
func (s *WalletService) ListTransactions(userID int64, page, pageSize int) ([]model.WalletTransaction, int64, error) {
	wallet, err := s.walletRepo.GetByUserID(userID)
	if err != nil {
		return nil, 0, err
	}
	return s.walletRepo.ListTransactions(wallet.ID, page, pageSize)
}

// Transfer 内部转账（交易结算）
func (s *WalletService) Transfer(fromUserID, toUserID int64, amount decimal.Decimal, refType string, refID int64, desc string) error {
	if amount.LessThanOrEqual(decimal.Zero) {
		return errors.New("转账金额必须大于0")
	}

	db := s.walletRepo.DB()
	return db.Transaction(func(tx *gorm.DB) error {
		// 锁定付款方钱包
		fromWallet, err := s.walletRepo.GetByUserIDForUpdate(tx, fromUserID)
		if err != nil {
			return errors.New("付款方钱包不存在")
		}

		if fromWallet.Balance.LessThan(amount) {
			return errors.New("巴特币余额不足")
		}

		// 锁定收款方钱包
		toWallet, err := s.walletRepo.GetByUserIDForUpdate(tx, toUserID)
		if err != nil {
			return errors.New("收款方钱包不存在")
		}

		// 扣减付款方
		newFromBalance := fromWallet.Balance.Sub(amount)
		if err := s.walletRepo.UpdateBalance(tx, fromWallet.ID, newFromBalance, fromWallet.FrozenBalance); err != nil {
			return err
		}
		if err := s.walletRepo.AddExpense(tx, fromWallet.ID, amount); err != nil {
			return err
		}

		// 增加收款方
		newToBalance := toWallet.Balance.Add(amount)
		if err := s.walletRepo.UpdateBalance(tx, toWallet.ID, newToBalance, toWallet.FrozenBalance); err != nil {
			return err
		}
		if err := s.walletRepo.AddIncome(tx, toWallet.ID, amount); err != nil {
			return err
		}

		// 记录付款方流水
		fromTxn := &model.WalletTransaction{
			WalletID:      fromWallet.ID,
			Type:          model.TxTypeTradeOut,
			Amount:        amount.Neg(),
			BalanceAfter:  newFromBalance,
			ReferenceType: refType,
			ReferenceID:   refID,
			Description:   desc,
		}
		if err := s.walletRepo.CreateTransaction(tx, fromTxn); err != nil {
			return err
		}

		// 记录收款方流水
		toTxn := &model.WalletTransaction{
			WalletID:      toWallet.ID,
			Type:          model.TxTypeTradeIn,
			Amount:        amount,
			BalanceAfter:  newToBalance,
			ReferenceType: refType,
			ReferenceID:   refID,
			Description:   desc,
		}
		return s.walletRepo.CreateTransaction(tx, toTxn)
	})
}
