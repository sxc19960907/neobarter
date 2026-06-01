package service

import (
	"testing"

	"github.com/neobarter/server/internal/model"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// 复用 setupTradeTestDB / newTradeService（定义在 trade_expiry_test.go）

func TestCounter_Flow(t *testing.T) {
	db := setupTradeTestDB(t)
	svc := newTradeService(db)

	// A(1) 向 B(2) 发起 pending 交易
	tr := &model.TradeRequest{InitiatorID: 1, TargetUserID: 2, TargetItemID: 1, Status: model.TradeStatusPending}
	require.NoError(t, db.Create(tr).Error)

	// B 反向提议：要求 50 巴特币
	counterItem := int64(9)
	err := svc.Counter(tr.ID, 2, &counterItem, decimal.NewFromInt(50), "加点巴特币吧")
	require.NoError(t, err)

	var after model.TradeRequest
	db.First(&after, tr.ID)
	assert.Equal(t, model.TradeStatusCountered, after.Status)
	require.NotNil(t, after.CounterItemID)
	assert.Equal(t, int64(9), *after.CounterItemID)
	assert.True(t, decimal.NewFromInt(50).Equal(after.CounterCoinAmount))

	// A 接受反向提议 -> 还价条件落地 + accepted
	err = svc.AcceptCounter(tr.ID, 1)
	require.NoError(t, err)

	db.First(&after, tr.ID)
	assert.Equal(t, model.TradeStatusAccepted, after.Status)
	require.NotNil(t, after.OfferedItemID)
	assert.Equal(t, int64(9), *after.OfferedItemID, "还价物品应落到 offered_item_id")
	assert.True(t, decimal.NewFromInt(50).Equal(after.BarterCoinAmount), "还价巴特币应生效")
}

func TestCounter_OnlyTargetCanCounter(t *testing.T) {
	db := setupTradeTestDB(t)
	svc := newTradeService(db)
	tr := &model.TradeRequest{InitiatorID: 1, TargetUserID: 2, TargetItemID: 1, Status: model.TradeStatusPending}
	require.NoError(t, db.Create(tr).Error)

	// 发起方 A(1) 不能反向提议自己的请求
	err := svc.Counter(tr.ID, 1, nil, decimal.Zero, "")
	assert.ErrorIs(t, err, ErrForbidden)
}

func TestAcceptCounter_OnlyInitiator(t *testing.T) {
	db := setupTradeTestDB(t)
	svc := newTradeService(db)
	tr := &model.TradeRequest{InitiatorID: 1, TargetUserID: 2, TargetItemID: 1, Status: model.TradeStatusCountered}
	require.NoError(t, db.Create(tr).Error)

	// 目标方 B(2) 不能接受自己提出的反向提议
	err := svc.AcceptCounter(tr.ID, 2)
	assert.ErrorIs(t, err, ErrForbidden)
}

func TestRejectCounter(t *testing.T) {
	db := setupTradeTestDB(t)
	svc := newTradeService(db)
	tr := &model.TradeRequest{InitiatorID: 1, TargetUserID: 2, TargetItemID: 1, Status: model.TradeStatusCountered}
	require.NoError(t, db.Create(tr).Error)

	err := svc.RejectCounter(tr.ID, 1, "条件不合适")
	require.NoError(t, err)

	var after model.TradeRequest
	db.First(&after, tr.ID)
	assert.Equal(t, model.TradeStatusRejected, after.Status)
	assert.Equal(t, "条件不合适", after.RejectReason)
}

func TestCounter_WrongStatus(t *testing.T) {
	db := setupTradeTestDB(t)
	svc := newTradeService(db)
	// accepted 状态不能再反向提议
	tr := &model.TradeRequest{InitiatorID: 1, TargetUserID: 2, TargetItemID: 1, Status: model.TradeStatusAccepted}
	require.NoError(t, db.Create(tr).Error)

	err := svc.Counter(tr.ID, 2, nil, decimal.Zero, "")
	assert.Error(t, err)
}
