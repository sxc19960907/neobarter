package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/neobarter/server/internal/middleware"
	"github.com/neobarter/server/internal/model"
	"github.com/neobarter/server/internal/pkg/response"
	"github.com/neobarter/server/internal/service"
	"github.com/shopspring/decimal"
)

type TradeHandler struct {
	tradeSvc *service.TradeService
}

func NewTradeHandler(tradeSvc *service.TradeService) *TradeHandler {
	return &TradeHandler{tradeSvc: tradeSvc}
}

type CreateTradeReq struct {
	TargetItemID     int64   `json:"target_item_id" binding:"required"`
	OfferedItemID    *int64  `json:"offered_item_id"`
	BarterCoinAmount float64 `json:"barter_coin_amount"`
	Message          string  `json:"message"`
}

func (h *TradeHandler) Create(c *gin.Context) {
	userID := middleware.GetUserID(c)
	var req CreateTradeReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	trade := &model.TradeRequest{
		TargetItemID:     req.TargetItemID,
		OfferedItemID:    req.OfferedItemID,
		BarterCoinAmount: decimal.NewFromFloat(req.BarterCoinAmount),
		Message:          req.Message,
	}

	if err := h.tradeSvc.Create(userID, trade); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Success(c, trade)
}

func (h *TradeHandler) List(c *gin.Context) {
	userID := middleware.GetUserID(c)
	page, pageSize := parsePage(c)
	status := c.Query("status")

	trades, total, err := h.tradeSvc.List(userID, status, page, pageSize)
	if err != nil {
		response.ServerError(c, "获取交易列表失败")
		return
	}

	response.SuccessPage(c, trades, total, page, pageSize)
}

func (h *TradeHandler) Get(c *gin.Context) {
	userID := middleware.GetUserID(c)
	id := parseID(c, "id")
	if id == 0 {
		return
	}

	trade, err := h.tradeSvc.Get(id, userID)
	if err != nil {
		response.NotFound(c, err.Error())
		return
	}

	response.Success(c, trade)
}

func (h *TradeHandler) Accept(c *gin.Context) {
	userID := middleware.GetUserID(c)
	id := parseID(c, "id")
	if id == 0 {
		return
	}

	if err := h.tradeSvc.Accept(id, userID); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Success(c, nil)
}

type RejectReq struct {
	Reason string `json:"reason" binding:"required"`
}

func (h *TradeHandler) Reject(c *gin.Context) {
	userID := middleware.GetUserID(c)
	id := parseID(c, "id")
	if id == 0 {
		return
	}

	var req RejectReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请填写拒绝原因")
		return
	}

	if err := h.tradeSvc.Reject(id, userID, req.Reason); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Success(c, nil)
}

func (h *TradeHandler) Complete(c *gin.Context) {
	userID := middleware.GetUserID(c)
	id := parseID(c, "id")
	if id == 0 {
		return
	}

	if err := h.tradeSvc.Complete(id, userID); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Success(c, nil)
}

func (h *TradeHandler) Cancel(c *gin.Context) {
	userID := middleware.GetUserID(c)
	id := parseID(c, "id")
	if id == 0 {
		return
	}

	if err := h.tradeSvc.Cancel(id, userID); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Success(c, nil)
}
