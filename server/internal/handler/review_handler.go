package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/neobarter/server/internal/middleware"
	"github.com/neobarter/server/internal/model"
	"github.com/neobarter/server/internal/pkg/response"
	"github.com/neobarter/server/internal/service"
)

type ReviewHandler struct {
	reviewSvc *service.ReviewService
}

func NewReviewHandler(reviewSvc *service.ReviewService) *ReviewHandler {
	return &ReviewHandler{reviewSvc: reviewSvc}
}

type CreateReviewReq struct {
	TradeRequestID int64  `json:"trade_request_id" binding:"required"`
	RevieweeID     int64  `json:"reviewee_id" binding:"required"`
	Rating         int    `json:"rating" binding:"required,min=1,max=5"`
	Comment        string `json:"comment"`
}

// Create 提交评价
// @Summary      提交交易评价
// @Description  交易完成后对对方评分（1-5星），影响信用积分
// @Tags         评价
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        body  body      CreateReviewReq  true  "评价内容"
// @Success      200   {object}  response.Response{data=model.Review}
// @Failure      400   {object}  response.Response
// @Router       /reviews [post]
func (h *ReviewHandler) Create(c *gin.Context) {
	userID := middleware.GetUserID(c)
	var req CreateReviewReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误，评分需在1-5之间")
		return
	}

	review := &model.Review{
		TradeRequestID: req.TradeRequestID,
		ReviewerID:     userID,
		RevieweeID:     req.RevieweeID,
		Rating:         req.Rating,
		Comment:        req.Comment,
	}

	if err := h.reviewSvc.Create(review); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Success(c, review)
}

// ListByUser 获取用户收到的评价
// @Summary      获取用户评价列表
// @Tags         评价
// @Produce      json
// @Security     BearerAuth
// @Param        id         path      int  true   "用户ID"
// @Param        page       query     int  false  "页码"
// @Param        page_size  query     int  false  "每页数量"
// @Success      200  {object}  response.Response{data=response.PageData}
// @Router       /reviews/user/{id} [get]
func (h *ReviewHandler) ListByUser(c *gin.Context) {
	id := parseID(c, "id")
	if id == 0 {
		return
	}

	page, pageSize := parsePage(c)
	reviews, total, err := h.reviewSvc.ListByUser(id, page, pageSize)
	if err != nil {
		response.ServerError(c, "获取评价列表失败")
		return
	}

	response.SuccessPage(c, reviews, total, page, pageSize)
}

// NotificationHandler

type NotificationHandler struct {
	notificationSvc *service.NotificationService
}

func NewNotificationHandler(notificationSvc *service.NotificationService) *NotificationHandler {
	return &NotificationHandler{notificationSvc: notificationSvc}
}

// List 获取通知列表
// @Summary      获取通知列表
// @Tags         通知
// @Produce      json
// @Security     BearerAuth
// @Param        page       query     int  false  "页码"
// @Param        page_size  query     int  false  "每页数量"
// @Success      200  {object}  response.Response{data=response.PageData}
// @Router       /notifications [get]
func (h *NotificationHandler) List(c *gin.Context) {
	userID := middleware.GetUserID(c)
	page, pageSize := parsePage(c)

	notifications, total, err := h.notificationSvc.List(userID, page, pageSize)
	if err != nil {
		response.ServerError(c, "获取通知列表失败")
		return
	}

	response.SuccessPage(c, notifications, total, page, pageSize)
}

// UnreadCount 获取未读通知数
// @Summary      获取未读通知数
// @Tags         通知
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  response.Response{data=object}
// @Router       /notifications/unread-count [get]
func (h *NotificationHandler) UnreadCount(c *gin.Context) {
	userID := middleware.GetUserID(c)
	count, err := h.notificationSvc.UnreadCount(userID)
	if err != nil {
		response.ServerError(c, "获取未读数失败")
		return
	}
	response.Success(c, gin.H{"count": count})
}

// MarkRead 标记通知已读
// @Summary      标记单条通知已读
// @Tags         通知
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      int  true  "通知ID"
// @Success      200  {object}  response.Response
// @Router       /notifications/{id}/read [put]
func (h *NotificationHandler) MarkRead(c *gin.Context) {
	userID := middleware.GetUserID(c)
	id := parseID(c, "id")
	if id == 0 {
		return
	}

	if err := h.notificationSvc.MarkRead(id, userID); err != nil {
		response.ServerError(c, "标记已读失败")
		return
	}
	response.Success(c, nil)
}

// MarkAllRead 标记全部通知已读
// @Summary      标记全部通知已读
// @Tags         通知
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  response.Response
// @Router       /notifications/read-all [put]
func (h *NotificationHandler) MarkAllRead(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if err := h.notificationSvc.MarkAllRead(userID); err != nil {
		response.ServerError(c, "标记全部已读失败")
		return
	}
	response.Success(c, nil)
}
