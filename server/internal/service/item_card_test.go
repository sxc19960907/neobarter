package service

import (
	"encoding/json"
	"testing"

	"github.com/neobarter/server/internal/model"
	"github.com/neobarter/server/internal/repository"
	"github.com/neobarter/server/internal/ws"
	"github.com/lib/pq"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupMsgDB(t *testing.T) (*MessageService, *gorm.DB) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)
	require.NoError(t, db.AutoMigrate(
		&model.User{}, &model.Item{}, &model.Conversation{},
		&model.ConversationParticipant{}, &model.Message{},
	))
	svc := NewMessageService(
		repository.NewMessageRepository(db),
		repository.NewItemRepository(db),
		ws.NewHub(),
	)
	return svc, db
}

func TestSendItemCard(t *testing.T) {
	svc, db := setupMsgDB(t)

	// 两个用户 + 一件物品
	require.NoError(t, db.Create(&model.User{Phone: "13800000001"}).Error) // id=1 sender
	require.NoError(t, db.Create(&model.User{Phone: "13800000002"}).Error) // id=2 receiver
	item := &model.Item{
		UserID:         2,
		Title:          "九成新相机",
		Condition:      "like_new",
		EstimatedValue: decimal.NewFromInt(80),
		Images:         pq.StringArray{"https://oss/cam.jpg"},
		Status:         model.ItemStatusActive,
	}
	require.NoError(t, db.Create(item).Error)

	// 用户1 给用户2 发物品卡片
	msg, err := svc.SendItemCard(1, 0, 2, item.ID)
	require.NoError(t, err)
	assert.Equal(t, model.MsgTypeItemCard, msg.MessageType)
	require.NotNil(t, msg.ExtraData)

	// extra_data 应是后端组装的可信卡片
	var card ItemCard
	require.NoError(t, json.Unmarshal([]byte(*msg.ExtraData), &card))
	assert.Equal(t, item.ID, card.ItemID)
	assert.Equal(t, "九成新相机", card.Title)
	assert.Equal(t, "https://oss/cam.jpg", card.Image)
	assert.Equal(t, "like_new", card.Condition)
	assert.Equal(t, "80", card.EstimatedValue)
}

func TestSendItemCard_ItemNotFound(t *testing.T) {
	svc, db := setupMsgDB(t)
	require.NoError(t, db.Create(&model.User{Phone: "13800000001"}).Error)
	require.NoError(t, db.Create(&model.User{Phone: "13800000002"}).Error)

	_, err := svc.SendItemCard(1, 0, 2, 999) // 不存在的物品
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "物品不存在")
}
