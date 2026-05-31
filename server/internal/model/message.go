package model

import "time"

type Conversation struct {
	ID            int64      `json:"id" gorm:"primaryKey"`
	Type          string     `json:"type" gorm:"size:20;not null;default:private"`
	LastMessageID *int64     `json:"last_message_id"`
	LastMessageAt *time.Time `json:"last_message_at"`
	CreatedAt     time.Time  `json:"created_at"`
}

func (Conversation) TableName() string {
	return "conversations"
}

type ConversationParticipant struct {
	ID             int64      `json:"id" gorm:"primaryKey"`
	ConversationID int64      `json:"conversation_id" gorm:"uniqueIndex:idx_conv_user;not null"`
	UserID         int64      `json:"user_id" gorm:"uniqueIndex:idx_conv_user;index;not null"`
	UnreadCount    int        `json:"unread_count" gorm:"not null;default:0"`
	LastReadAt     *time.Time `json:"last_read_at"`
	CreatedAt      time.Time  `json:"created_at"`

	User *User `json:"user,omitempty" gorm:"foreignKey:UserID"`
}

func (ConversationParticipant) TableName() string {
	return "conversation_participants"
}

type Message struct {
	ID             int64     `json:"id" gorm:"primaryKey"`
	ConversationID int64     `json:"conversation_id" gorm:"index;not null"`
	SenderID       int64     `json:"sender_id" gorm:"index;not null"`
	Content        string    `json:"content" gorm:"not null"`
	MessageType    string    `json:"message_type" gorm:"size:20;not null;default:text"`
	ExtraData      *string   `json:"extra_data" gorm:"type:jsonb"`
	CreatedAt      time.Time `json:"created_at"`

	Sender *User `json:"sender,omitempty" gorm:"foreignKey:SenderID"`
}

func (Message) TableName() string {
	return "messages"
}

// 消息类型常量
const (
	MsgTypeText     = "text"
	MsgTypeImage    = "image"
	MsgTypeVoice    = "voice"
	MsgTypeItemCard = "item_card"
)
