package service

import (
	"errors"

	"github.com/neobarter/server/internal/model"
	"github.com/neobarter/server/internal/repository"
	"github.com/neobarter/server/internal/ws"
	"gorm.io/gorm"
)

type MessageService struct {
	messageRepo *repository.MessageRepository
	wsHub       *ws.Hub
}

func NewMessageService(messageRepo *repository.MessageRepository, wsHub *ws.Hub) *MessageService {
	return &MessageService{
		messageRepo: messageRepo,
		wsHub:       wsHub,
	}
}

// GetOrCreateConversation 获取或创建私聊会话
func (s *MessageService) GetOrCreateConversation(userID1, userID2 int64) (*model.Conversation, error) {
	conv, err := s.messageRepo.FindPrivateConversation(userID1, userID2)
	if err == nil {
		return conv, nil
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	// 创建新会话
	conv = &model.Conversation{Type: "private"}
	if err := s.messageRepo.CreateConversation(conv); err != nil {
		return nil, err
	}

	// 添加参与者
	s.messageRepo.AddParticipant(&model.ConversationParticipant{
		ConversationID: conv.ID,
		UserID:         userID1,
	})
	s.messageRepo.AddParticipant(&model.ConversationParticipant{
		ConversationID: conv.ID,
		UserID:         userID2,
	})

	return conv, nil
}

// Send 发送消息
func (s *MessageService) Send(senderID int64, conversationID int64, content, msgType string, extraData *string) (*model.Message, error) {
	// 验证是否为会话参与者
	if !s.messageRepo.IsParticipant(conversationID, senderID) {
		return nil, ErrForbidden
	}

	msg := &model.Message{
		ConversationID: conversationID,
		SenderID:       senderID,
		Content:        content,
		MessageType:    msgType,
		ExtraData:      extraData,
	}

	if err := s.messageRepo.CreateMessage(msg); err != nil {
		return nil, err
	}

	// 更新会话最后消息
	s.messageRepo.UpdateConversationLastMessage(conversationID, msg.ID)

	// 增加对方未读数
	s.messageRepo.IncrUnreadCount(conversationID, senderID)

	// 通过 WebSocket 精确推送给会话其他参与者（排除发送者）
	if recipients, err := s.messageRepo.ParticipantIDsExcept(conversationID, senderID); err == nil {
		s.wsHub.SendToUsers(recipients, "new_message", msg)
	}

	return msg, nil
}

// SendToUser 发送消息给指定用户（自动创建/获取会话）
func (s *MessageService) SendToUser(senderID, receiverID int64, content, msgType string) (*model.Message, error) {
	conv, err := s.GetOrCreateConversation(senderID, receiverID)
	if err != nil {
		return nil, err
	}
	return s.Send(senderID, conv.ID, content, msgType, nil)
}

// ListConversations 获取用户会话列表
func (s *MessageService) ListConversations(userID int64) ([]model.ConversationParticipant, error) {
	return s.messageRepo.ListConversations(userID)
}

// GetMessages 获取会话消息
func (s *MessageService) GetMessages(conversationID, userID int64, page, pageSize int) ([]model.Message, int64, error) {
	if !s.messageRepo.IsParticipant(conversationID, userID) {
		return nil, 0, ErrForbidden
	}
	return s.messageRepo.ListMessages(conversationID, page, pageSize)
}

// MarkRead 标记会话已读
func (s *MessageService) MarkRead(conversationID, userID int64) error {
	return s.messageRepo.ClearUnreadCount(conversationID, userID)
}
