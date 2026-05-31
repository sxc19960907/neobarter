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

func (h *NotificationHandler) UnreadCount(c *gin.Context) {
	userID := middleware.GetUserID(c)
	count, err := h.notificationSvc.UnreadCount(userID)
	if err != nil {
		response.ServerError(c, "获取未读数失败")
		return
	}
	response.Success(c, gin.H{"count": count})
}

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

func (h *NotificationHandler) MarkAllRead(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if err := h.notificationSvc.MarkAllRead(userID); err != nil {
		response.ServerError(c, "标记全部已读失败")
		return
	}
	response.Success(c, nil)
}
