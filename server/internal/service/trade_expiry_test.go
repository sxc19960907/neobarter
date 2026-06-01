package service

import (
	"testing"
	"time"

	"github.com/neobarter/server/internal/model"
	"github.com/neobarter/server/internal/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTradeTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)
	require.NoError(t, db.AutoMigrate(
		&model.User{}, &model.Item{}, &model.TradeRequest{}, &model.Notification{},
	))
	return db
}

func newTradeService(db *gorm.DB) *TradeService {
	return NewTradeService(
		repository.NewTradeRepository(db),
		repository.NewItemRepository(db),
		nil, // walletSvc 不参与过期逻辑
		repository.NewNotificationRepository(db),
	)
}

func TestExpireStale(t *testing.T) {
	db := setupTradeTestDB(t)
	svc := newTradeService(db)

	past := time.Now().Add(-1 * time.Hour)
	future := time.Now().Add(1 * time.Hour)

	// 1) 已超时的 pending —— 应被过期
	expiredPast := &model.TradeRequest{InitiatorID: 1, TargetUserID: 2, TargetItemID: 1, Status: model.TradeStatusPending, ExpiredAt: &past}
	// 2) 未超时的 pending —— 不动
	stillValid := &model.TradeRequest{InitiatorID: 1, TargetUserID: 2, TargetItemID: 2, Status: model.TradeStatusPending, ExpiredAt: &future}
	// 3) 已 accepted 但 expired_at 过期 —— 不应被影响
	acceptedPast := &model.TradeRequest{InitiatorID: 1, TargetUserID: 2, TargetItemID: 3, Status: model.TradeStatusAccepted, ExpiredAt: &past}
	require.NoError(t, db.Create(expiredPast).Error)
	require.NoError(t, db.Create(stillValid).Error)
	require.NoError(t, db.Create(acceptedPast).Error)

	n, err := svc.ExpireStale()
	require.NoError(t, err)
	assert.Equal(t, 1, n, "只应过期 1 笔")

	var t1, t2, t3 model.TradeRequest
	db.First(&t1, expiredPast.ID)
	db.First(&t2, stillValid.ID)
	db.First(&t3, acceptedPast.ID)
	assert.Equal(t, model.TradeStatusExpired, t1.Status, "超时 pending 应 expired")
	assert.Equal(t, model.TradeStatusPending, t2.Status, "未超时 pending 不变")
	assert.Equal(t, model.TradeStatusAccepted, t3.Status, "accepted 不受影响")

	// 发起方应收到过期通知
	var notifCount int64
	db.Model(&model.Notification{}).Where("user_id = ? AND reference_id = ?", 1, expiredPast.ID).Count(&notifCount)
	assert.Equal(t, int64(1), notifCount, "发起方应收到过期通知")
}

func TestExpireStale_NothingToExpire(t *testing.T) {
	db := setupTradeTestDB(t)
	svc := newTradeService(db)

	future := time.Now().Add(1 * time.Hour)
	db.Create(&model.TradeRequest{InitiatorID: 1, TargetUserID: 2, TargetItemID: 1, Status: model.TradeStatusPending, ExpiredAt: &future})

	n, err := svc.ExpireStale()
	require.NoError(t, err)
	assert.Equal(t, 0, n)
}
