package repository

import (
	"github.com/neobarter/server/internal/model"
	"gorm.io/gorm"
)

type MessageRepository struct {
	db *gorm.DB
}

func NewMessageRepository(db *gorm.DB) *MessageRepository {
	return &MessageRepository{db: db}
}

func (r *MessageRepository) CreateConversation(conv *model.Conversation) error {
	return r.db.Create(conv).Error
}

func (r *MessageRepository) AddParticipant(p *model.ConversationParticipant) error {
	return r.db.Create(p).Error
}

func (r *MessageRepository) GetConversation(id int64) (*model.Conversation, error) {
	var conv model.Conversation
	err := r.db.First(&conv, id).Error
	if err != nil {
		return nil, err
	}
	return &conv, nil
}

// FindPrivateConversation 查找两个用户之间的私聊会话
func (r *MessageRepository) FindPrivateConversation(userID1, userID2 int64) (*model.Conversation, error) {
	var conv model.Conversation
	err := r.db.Raw(`
		SELECT c.* FROM conversations c
		JOIN conversation_participants cp1 ON cp1.conversation_id = c.id AND cp1.user_id = ?
		JOIN conversation_participants cp2 ON cp2.conversation_id = c.id AND cp2.user_id = ?
		WHERE c.type = 'private'
		LIMIT 1
	`, userID1, userID2).Scan(&conv).Error
	if err != nil {
		return nil, err
	}
	if conv.ID == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	return &conv, nil
}

func (r *MessageRepository) ListConversations(userID int64) ([]model.ConversationParticipant, error) {
	var participants []model.ConversationParticipant
	err := r.db.Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&participants).Error
	return participants, err
}

func (r *MessageRepository) CreateMessage(msg *model.Message) error {
	return r.db.Create(msg).Error
}

func (r *MessageRepository) ListMessages(conversationID int64, page, pageSize int) ([]model.Message, int64, error) {
	var messages []model.Message
	var total int64

	query := r.db.Where("conversation_id = ?", conversationID)
	query.Model(&model.Message{}).Count(&total)

	err := query.Preload("Sender").
		Order("created_at DESC").
		Offset((page - 1) * pageSize).Limit(pageSize).
		Find(&messages).Error

	return messages, total, err
}

func (r *MessageRepository) UpdateConversationLastMessage(convID int64, msgID int64) error {
	return r.db.Model(&model.Conversation{}).Where("id = ?", convID).
		Updates(map[string]interface{}{
			"last_message_id": msgID,
			"last_message_at": gorm.Expr("NOW()"),
		}).Error
}

func (r *MessageRepository) IncrUnreadCount(convID, userID int64) error {
	return r.db.Model(&model.ConversationParticipant{}).
		Where("conversation_id = ? AND user_id != ?", convID, userID).
		UpdateColumn("unread_count", gorm.Expr("unread_count + 1")).Error
}

func (r *MessageRepository) ClearUnreadCount(convID, userID int64) error {
	return r.db.Model(&model.ConversationParticipant{}).
		Where("conversation_id = ? AND user_id = ?", convID, userID).
		Updates(map[string]interface{}{
			"unread_count": 0,
			"last_read_at": gorm.Expr("NOW()"),
		}).Error
}

func (r *MessageRepository) IsParticipant(convID, userID int64) bool {
	var count int64
	r.db.Model(&model.ConversationParticipant{}).
		Where("conversation_id = ? AND user_id = ?", convID, userID).
		Count(&count)
	return count > 0
}

// ParticipantIDsExcept 返回会话参与者的 userID 列表，排除指定用户（通常是消息发送者）。
// 用于 WebSocket 精确推送。
func (r *MessageRepository) ParticipantIDsExcept(convID, excludeUserID int64) ([]int64, error) {
	var ids []int64
	err := r.db.Model(&model.ConversationParticipant{}).
		Where("conversation_id = ? AND user_id != ?", convID, excludeUserID).
		Pluck("user_id", &ids).Error
	return ids, err
}
