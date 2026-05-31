package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/neobarter/server/internal/middleware"
	"github.com/neobarter/server/internal/model"
	"github.com/neobarter/server/internal/pkg/response"
	"github.com/neobarter/server/internal/service"
)

type MessageHandler struct {
	messageSvc *service.MessageService
}

func NewMessageHandler(messageSvc *service.MessageService) *MessageHandler {
	return &MessageHandler{messageSvc: messageSvc}
}

func (h *MessageHandler) ListConversations(c *gin.Context) {
	userID := middleware.GetUserID(c)
	conversations, err := h.messageSvc.ListConversations(userID)
	if err != nil {
		response.ServerError(c, "获取会话列表失败")
		return
	}
	response.Success(c, conversations)
}

func (h *MessageHandler) GetMessages(c *gin.Context) {
	userID := middleware.GetUserID(c)
	convID := parseID(c, "id")
	if convID == 0 {
		return
	}

	page, pageSize := parsePage(c)
	messages, total, err := h.messageSvc.GetMessages(convID, userID, page, pageSize)
	if err != nil {
		response.Forbidden(c, err.Error())
		return
	}

	response.SuccessPage(c, messages, total, page, pageSize)
}

type SendMessageReq struct {
	ConversationID int64  `json:"conversation_id"`
	ReceiverID     int64  `json:"receiver_id"`
	Content        string `json:"content" binding:"required"`
	MessageType    string `json:"message_type"`
}

func (h *MessageHandler) Send(c *gin.Context) {
	userID := middleware.GetUserID(c)
	var req SendMessageReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	if req.MessageType == "" {
		req.MessageType = model.MsgTypeText
	}

	var msg *model.Message
	var err error

	if req.ConversationID > 0 {
		msg, err = h.messageSvc.Send(userID, req.ConversationID, req.Content, req.MessageType, nil)
	} else if req.ReceiverID > 0 {
		msg, err = h.messageSvc.SendToUser(userID, req.ReceiverID, req.Content, req.MessageType)
	} else {
		response.BadRequest(c, "请指定会话或接收者")
		return
	}

	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Success(c, msg)
}

func (h *MessageHandler) MarkRead(c *gin.Context) {
	userID := middleware.GetUserID(c)
	convID := parseID(c, "id")
	if convID == 0 {
		return
	}

	if err := h.messageSvc.MarkRead(convID, userID); err != nil {
		response.ServerError(c, "标记已读失败")
		return
	}

	response.Success(c, nil)
}
