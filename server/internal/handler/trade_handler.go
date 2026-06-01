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

// Create 发起交换请求
// @Summary      发起交换请求
// @Description  对目标物品发起交换，可附带巴特币补差价
// @Tags         交易
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        body  body      CreateTradeReq  true  "交换请求"
// @Success      200   {object}  response.Response{data=model.TradeRequest}
// @Failure      400   {object}  response.Response
// @Router       /trades [post]
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

// List 获取交易列表
// @Summary      获取我的交易列表
// @Tags         交易
// @Produce      json
// @Security     BearerAuth
// @Param        status     query     string  false  "状态筛选：pending/accepted/rejected/completed/cancelled"
// @Param        page       query     int     false  "页码"
// @Param        page_size  query     int     false  "每页数量"
// @Success      200  {object}  response.Response{data=response.PageData}
// @Router       /trades [get]
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

// Get 获取交易详情
// @Summary      获取交易详情
// @Tags         交易
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      int  true  "交易ID"
// @Success      200  {object}  response.Response{data=model.TradeRequest}
// @Failure      404  {object}  response.Response
// @Router       /trades/{id} [get]
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

// Accept 接受交换
// @Summary      接受交换请求
// @Tags         交易
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      int  true  "交易ID"
// @Success      200  {object}  response.Response
// @Failure      400  {object}  response.Response
// @Router       /trades/{id}/accept [put]
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

// Reject 拒绝交换
// @Summary      拒绝交换请求
// @Tags         交易
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id    path      int        true  "交易ID"
// @Param        body  body      RejectReq  true  "拒绝原因"
// @Success      200   {object}  response.Response
// @Failure      400   {object}  response.Response
// @Router       /trades/{id}/reject [put]
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

// Complete 完成交易
// @Summary      确认完成交易
// @Description  双方确认后结算巴特币，物品标记为已交易
// @Tags         交易
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      int  true  "交易ID"
// @Success      200  {object}  response.Response
// @Failure      400  {object}  response.Response
// @Router       /trades/{id}/complete [put]
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

// Cancel 取消交易
// @Summary      取消交易请求
// @Description  仅发起方可取消待处理的交易
// @Tags         交易
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      int  true  "交易ID"
// @Success      200  {object}  response.Response
// @Failure      400  {object}  response.Response
// @Router       /trades/{id}/cancel [put]
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
