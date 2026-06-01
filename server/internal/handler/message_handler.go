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

// ListConversations 获取会话列表
// @Summary      获取会话列表
// @Tags         消息
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  response.Response{data=[]model.ConversationParticipant}
// @Router       /messages/conversations [get]
func (h *MessageHandler) ListConversations(c *gin.Context) {
	userID := middleware.GetUserID(c)
	conversations, err := h.messageSvc.ListConversations(userID)
	if err != nil {
		response.ServerError(c, "获取会话列表失败")
		return
	}
	response.Success(c, conversations)
}

// GetMessages 获取会话消息历史
// @Summary      获取会话消息
// @Tags         消息
// @Produce      json
// @Security     BearerAuth
// @Param        id         path      int  true   "会话ID"
// @Param        page       query     int  false  "页码"
// @Param        page_size  query     int  false  "每页数量"
// @Success      200  {object}  response.Response{data=response.PageData}
// @Failure      403  {object}  response.Response
// @Router       /messages/conversations/{id} [get]
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

// Send 发送消息
// @Summary      发送消息
// @Description  指定 conversation_id 或 receiver_id（自动创建会话）
// @Tags         消息
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        body  body      SendMessageReq  true  "消息内容"
// @Success      200   {object}  response.Response{data=model.Message}
// @Failure      400   {object}  response.Response
// @Router       /messages [post]
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

type SendItemCardReq struct {
	ConversationID int64 `json:"conversation_id"`
	ReceiverID     int64 `json:"receiver_id"`
	ItemID         int64 `json:"item_id" binding:"required"`
}

// SendItemCard 发送物品卡片消息
// @Summary      发送物品卡片
// @Description  分享物品卡片到会话，卡片数据由后端根据 item_id 组装
// @Tags         消息
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        body  body      SendItemCardReq  true  "物品卡片"
// @Success      200   {object}  response.Response{data=model.Message}
// @Failure      400   {object}  response.Response
// @Router       /messages/item-card [post]
func (h *MessageHandler) SendItemCard(c *gin.Context) {
	userID := middleware.GetUserID(c)
	var req SendItemCardReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请指定物品")
		return
	}

	msg, err := h.messageSvc.SendItemCard(userID, req.ConversationID, req.ReceiverID, req.ItemID)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Success(c, msg)
}

// MarkRead 标记会话已读
// @Summary      标记会话已读
// @Tags         消息
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      int  true  "会话ID"
// @Success      200  {object}  response.Response
// @Router       /messages/conversations/{id}/read [put]
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
