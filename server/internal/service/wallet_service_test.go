package service

import (
	"testing"

	"github.com/neobarter/server/internal/model"
	"github.com/neobarter/server/internal/repository"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	err = db.AutoMigrate(&model.User{}, &model.Wallet{}, &model.WalletTransaction{})
	require.NoError(t, err)

	return db
}

func createTestUser(t *testing.T, db *gorm.DB, phone string) *model.User {
	user := &model.User{Phone: phone, Nickname: "test", UserType: "personal"}
	require.NoError(t, db.Create(user).Error)
	return user
}

func TestWalletService_CreateWalletWithReward(t *testing.T) {
	db := setupTestDB(t)
	walletRepo := repository.NewWalletRepository(db)
	svc := NewWalletService(walletRepo, 100.0)

	user := createTestUser(t, db, "13800000001")

	err := svc.CreateWalletWithReward(user.ID)
	require.NoError(t, err)

	// 验证钱包创建
	wallet, err := svc.GetWallet(user.ID)
	require.NoError(t, err)
	assert.True(t, decimal.NewFromFloat(100.0).Equal(wallet.Balance))
	assert.True(t, decimal.NewFromFloat(100.0).Equal(wallet.TotalIncome))

	// 验证流水记录
	txns, total, err := svc.ListTransactions(user.ID, 1, 10)
	require.NoError(t, err)
	assert.Equal(t, int64(1), total)
	assert.Equal(t, model.TxTypeReward, txns[0].Type)
	assert.Equal(t, "注册赠送巴特币", txns[0].Description)
}

func TestWalletService_Transfer_Success(t *testing.T) {
	db := setupTestDB(t)
	walletRepo := repository.NewWalletRepository(db)
	svc := NewWalletService(walletRepo, 100.0)

	user1 := createTestUser(t, db, "13800000001")
	user2 := createTestUser(t, db, "13800000002")

	svc.CreateWalletWithReward(user1.ID)
	svc.CreateWalletWithReward(user2.ID)

	// user1 转 30 巴特币给 user2
	err := svc.Transfer(user1.ID, user2.ID, decimal.NewFromFloat(30.0), "trade_request", 1, "交易结算")
	require.NoError(t, err)

	// 验证余额
	w1, _ := svc.GetWallet(user1.ID)
	w2, _ := svc.GetWallet(user2.ID)
	assert.True(t, decimal.NewFromFloat(70.0).Equal(w1.Balance))
	assert.True(t, decimal.NewFromFloat(130.0).Equal(w2.Balance))

	// 验证累计收支
	assert.True(t, decimal.NewFromFloat(30.0).Equal(w1.TotalExpense))
	assert.True(t, decimal.NewFromFloat(130.0).Equal(w2.TotalIncome)) // 100 reward + 30 transfer
}

func TestWalletService_Transfer_InsufficientBalance(t *testing.T) {
	db := setupTestDB(t)
	walletRepo := repository.NewWalletRepository(db)
	svc := NewWalletService(walletRepo, 100.0)

	user1 := createTestUser(t, db, "13800000001")
	user2 := createTestUser(t, db, "13800000002")

	svc.CreateWalletWithReward(user1.ID)
	svc.CreateWalletWithReward(user2.ID)

	// 尝试转 200 巴特币（余额只有 100）
	err := svc.Transfer(user1.ID, user2.ID, decimal.NewFromFloat(200.0), "trade_request", 1, "交易结算")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "余额不足")

	// 余额不变
	w1, _ := svc.GetWallet(user1.ID)
	assert.True(t, decimal.NewFromFloat(100.0).Equal(w1.Balance))
}

func TestWalletService_Transfer_ZeroAmount(t *testing.T) {
	db := setupTestDB(t)
	walletRepo := repository.NewWalletRepository(db)
	svc := NewWalletService(walletRepo, 100.0)

	user1 := createTestUser(t, db, "13800000001")
	user2 := createTestUser(t, db, "13800000002")

	svc.CreateWalletWithReward(user1.ID)
	svc.CreateWalletWithReward(user2.ID)

	err := svc.Transfer(user1.ID, user2.ID, decimal.Zero, "trade_request", 1, "test")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "大于0")
}
