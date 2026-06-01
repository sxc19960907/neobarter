package service

import (
	"errors"
	"fmt"
	"time"

	"github.com/neobarter/server/internal/model"
	"github.com/neobarter/server/internal/repository"
	"github.com/shopspring/decimal"
)

type TradeService struct {
	tradeRepo        *repository.TradeRepository
	itemRepo         *repository.ItemRepository
	walletSvc        *WalletService
	notificationRepo *repository.NotificationRepository
}

func NewTradeService(
	tradeRepo *repository.TradeRepository,
	itemRepo *repository.ItemRepository,
	walletSvc *WalletService,
	notificationRepo *repository.NotificationRepository,
) *TradeService {
	return &TradeService{
		tradeRepo:        tradeRepo,
		itemRepo:         itemRepo,
		walletSvc:        walletSvc,
		notificationRepo: notificationRepo,
	}
}

// Create 发起交换请求
func (s *TradeService) Create(initiatorID int64, req *model.TradeRequest) error {
	// 验证目标物品存在且可交易
	targetItem, err := s.itemRepo.GetByID(req.TargetItemID)
	if err != nil {
		return errors.New("目标物品不存在")
	}
	if targetItem.Status != model.ItemStatusActive {
		return errors.New("目标物品不可交易")
	}
	if targetItem.UserID == initiatorID {
		return errors.New("不能与自己交换")
	}

	req.InitiatorID = initiatorID
	req.TargetUserID = targetItem.UserID
	req.Status = model.TradeStatusPending

	// 设置24小时过期
	expiredAt := time.Now().Add(24 * time.Hour)
	req.ExpiredAt = &expiredAt

	if err := s.tradeRepo.Create(req); err != nil {
		return err
	}

	// 发送通知
	s.notificationRepo.Create(&model.Notification{
		UserID:        targetItem.UserID,
		Type:          model.NotifyTradeRequest,
		Title:         "收到新的交换请求",
		Content:       fmt.Sprintf("有人想交换你的「%s」", targetItem.Title),
		ReferenceType: "trade_request",
		ReferenceID:   req.ID,
	})

	return nil
}

// Accept 接受交换
func (s *TradeService) Accept(tradeID, userID int64) error {
	trade, err := s.tradeRepo.GetByID(tradeID)
	if err != nil {
		return ErrNotFound
	}
	if trade.TargetUserID != userID {
		return ErrForbidden
	}
	if trade.Status != model.TradeStatusPending {
		return errors.New("交易状态不允许此操作")
	}

	trade.Status = model.TradeStatusAccepted
	if err := s.tradeRepo.Update(trade); err != nil {
		return err
	}

	// 通知发起方
	s.notificationRepo.Create(&model.Notification{
		UserID:        trade.InitiatorID,
		Type:          model.NotifyTradeAccepted,
		Title:         "交换请求已被接受",
		Content:       "对方已接受你的交换请求，请确认完成交易",
		ReferenceType: "trade_request",
		ReferenceID:   trade.ID,
	})

	return nil
}

// Reject 拒绝交换
func (s *TradeService) Reject(tradeID, userID int64, reason string) error {
	trade, err := s.tradeRepo.GetByID(tradeID)
	if err != nil {
		return ErrNotFound
	}
	if trade.TargetUserID != userID {
		return ErrForbidden
	}
	if trade.Status != model.TradeStatusPending {
		return errors.New("交易状态不允许此操作")
	}

	trade.Status = model.TradeStatusRejected
	trade.RejectReason = reason
	if err := s.tradeRepo.Update(trade); err != nil {
		return err
	}

	// 通知发起方
	s.notificationRepo.Create(&model.Notification{
		UserID:        trade.InitiatorID,
		Type:          model.NotifyTradeRejected,
		Title:         "交换请求被拒绝",
		Content:       fmt.Sprintf("对方拒绝了你的交换请求，原因：%s", reason),
		ReferenceType: "trade_request",
		ReferenceID:   trade.ID,
	})

	return nil
}

// Counter 反向提议：目标用户(B)对 pending 交易还价。
func (s *TradeService) Counter(tradeID, userID int64, counterItemID *int64, coinAmount decimal.Decimal, message string) error {
	trade, err := s.tradeRepo.GetByID(tradeID)
	if err != nil {
		return ErrNotFound
	}
	if trade.TargetUserID != userID {
		return ErrForbidden
	}
	if trade.Status != model.TradeStatusPending {
		return errors.New("只能对待处理的交易发起反向提议")
	}

	trade.Status = model.TradeStatusCountered
	trade.CounterItemID = counterItemID
	trade.CounterCoinAmount = coinAmount
	trade.CounterMessage = message
	if err := s.tradeRepo.Update(trade); err != nil {
		return err
	}

	s.notificationRepo.Create(&model.Notification{
		UserID:        trade.InitiatorID,
		Type:          model.NotifyTradeRequest,
		Title:         "收到反向提议",
		Content:       "对方对你的交换请求提出了新条件，请确认",
		ReferenceType: "trade_request",
		ReferenceID:   trade.ID,
	})
	return nil
}

// AcceptCounter 发起方(A)接受反向提议：把还价条件落到生效字段，状态转 accepted。
func (s *TradeService) AcceptCounter(tradeID, userID int64) error {
	trade, err := s.tradeRepo.GetByID(tradeID)
	if err != nil {
		return ErrNotFound
	}
	if trade.InitiatorID != userID {
		return ErrForbidden
	}
	if trade.Status != model.TradeStatusCountered {
		return errors.New("交易状态不允许此操作")
	}

	// 套用还价条件，使后续 Complete 结算逻辑无需改动
	trade.OfferedItemID = trade.CounterItemID
	trade.BarterCoinAmount = trade.CounterCoinAmount
	trade.Status = model.TradeStatusAccepted
	if err := s.tradeRepo.Update(trade); err != nil {
		return err
	}

	s.notificationRepo.Create(&model.Notification{
		UserID:        trade.TargetUserID,
		Type:          model.NotifyTradeAccepted,
		Title:         "反向提议已被接受",
		Content:       "对方接受了你的反向提议，请确认完成交易",
		ReferenceType: "trade_request",
		ReferenceID:   trade.ID,
	})
	return nil
}

// RejectCounter 发起方(A)拒绝反向提议。
func (s *TradeService) RejectCounter(tradeID, userID int64, reason string) error {
	trade, err := s.tradeRepo.GetByID(tradeID)
	if err != nil {
		return ErrNotFound
	}
	if trade.InitiatorID != userID {
		return ErrForbidden
	}
	if trade.Status != model.TradeStatusCountered {
		return errors.New("交易状态不允许此操作")
	}

	trade.Status = model.TradeStatusRejected
	trade.RejectReason = reason
	if err := s.tradeRepo.Update(trade); err != nil {
		return err
	}

	s.notificationRepo.Create(&model.Notification{
		UserID:        trade.TargetUserID,
		Type:          model.NotifyTradeRejected,
		Title:         "反向提议被拒绝",
		Content:       fmt.Sprintf("对方拒绝了你的反向提议，原因：%s", reason),
		ReferenceType: "trade_request",
		ReferenceID:   trade.ID,
	})
	return nil
}

// Complete 完成交易（双方确认后结算巴特币）
func (s *TradeService) Complete(tradeID, userID int64) error {
	trade, err := s.tradeRepo.GetByID(tradeID)
	if err != nil {
		return ErrNotFound
	}
	if trade.InitiatorID != userID && trade.TargetUserID != userID {
		return ErrForbidden
	}
	if trade.Status != model.TradeStatusAccepted {
		return errors.New("交易状态不允许此操作")
	}

	// 巴特币结算
	if trade.BarterCoinAmount.GreaterThan(decimal.Zero) {
		desc := fmt.Sprintf("交易结算 #%d", trade.ID)
		if err := s.walletSvc.Transfer(
			trade.InitiatorID, trade.TargetUserID,
			trade.BarterCoinAmount, "trade_request", trade.ID, desc,
		); err != nil {
			return fmt.Errorf("巴特币结算失败: %w", err)
		}
	}

	// 更新交易状态
	now := time.Now()
	trade.Status = model.TradeStatusCompleted
	trade.CompletedAt = &now
	if err := s.tradeRepo.Update(trade); err != nil {
		return err
	}

	// 更新物品状态
	s.itemRepo.UpdateStatus(trade.TargetItemID, model.ItemStatusTraded)
	if trade.OfferedItemID != nil {
		s.itemRepo.UpdateStatus(*trade.OfferedItemID, model.ItemStatusTraded)
	}

	return nil
}

// Cancel 取消交易
func (s *TradeService) Cancel(tradeID, userID int64) error {
	trade, err := s.tradeRepo.GetByID(tradeID)
	if err != nil {
		return ErrNotFound
	}
	if trade.InitiatorID != userID {
		return ErrForbidden
	}
	if trade.Status != model.TradeStatusPending {
		return errors.New("只能取消待处理的交易")
	}

	trade.Status = model.TradeStatusCancelled
	return s.tradeRepo.Update(trade)
}

// Get 获取交易详情
func (s *TradeService) Get(tradeID, userID int64) (*model.TradeRequest, error) {
	trade, err := s.tradeRepo.GetByID(tradeID)
	if err != nil {
		return nil, ErrNotFound
	}
	if trade.InitiatorID != userID && trade.TargetUserID != userID {
		return nil, ErrForbidden
	}
	return trade, nil
}

// List 获取用户交易列表
func (s *TradeService) List(userID int64, status string, page, pageSize int) ([]model.TradeRequest, int64, error) {
	return s.tradeRepo.ListByUser(userID, status, page, pageSize)
}

// ExpireStale 把所有超时的 pending 交易置为 expired，并通知发起方。
// 返回过期的交易数量。由后台定时任务周期调用。
func (s *TradeService) ExpireStale() (int, error) {
	// 先查出待过期的交易（用于通知），再批量更新
	expired, err := s.tradeRepo.FindExpiredPending()
	if err != nil {
		return 0, err
	}
	if len(expired) == 0 {
		return 0, nil
	}

	if _, err := s.tradeRepo.ExpirePending(); err != nil {
		return 0, err
	}

	// 通知发起方交易已过期
	for _, t := range expired {
		s.notificationRepo.Create(&model.Notification{
			UserID:        t.InitiatorID,
			Type:          model.NotifyTradeRejected,
			Title:         "交换请求已过期",
			Content:       "对方24小时内未处理，交换请求已自动关闭",
			ReferenceType: "trade_request",
			ReferenceID:   t.ID,
		})
	}

	return len(expired), nil
}
