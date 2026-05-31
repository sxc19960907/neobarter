package service

import (
	"testing"

	"github.com/neobarter/server/internal/model"
	"github.com/neobarter/server/internal/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupItemTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	err = db.AutoMigrate(&model.User{}, &model.Item{}, &model.Category{})
	require.NoError(t, err)

	// 创建测试分类
	db.Create(&model.Category{Name: "数码电子", Icon: "laptop", SortOrder: 1})

	return db
}

func TestItemService_Create(t *testing.T) {
	db := setupItemTestDB(t)
	itemRepo := repository.NewItemRepository(db)
	svc := NewItemService(itemRepo, nil) // nil publisher = no MQ

	user := &model.User{Phone: "13800000001"}
	db.Create(user)

	categoryID := 1
	item := &model.Item{
		UserID:     user.ID,
		Title:      "iPhone 15 Pro",
		Condition:  "like_new",
		CategoryID: &categoryID,
		Status:     model.ItemStatusActive,
	}

	err := svc.Create(item)
	require.NoError(t, err)
	assert.NotZero(t, item.ID)
}

func TestItemService_Delete_OwnerOnly(t *testing.T) {
	db := setupItemTestDB(t)
	itemRepo := repository.NewItemRepository(db)
	svc := NewItemService(itemRepo, nil)

	user1 := &model.User{Phone: "13800000001"}
	user2 := &model.User{Phone: "13800000002"}
	db.Create(user1)
	db.Create(user2)

	item := &model.Item{UserID: user1.ID, Title: "Test Item", Condition: "good", Status: model.ItemStatusActive}
	db.Create(item)

	// user2 尝试删除 user1 的物品
	err := svc.Delete(item.ID, user2.ID)
	assert.ErrorIs(t, err, ErrForbidden)

	// user1 可以删除自己的物品
	err = svc.Delete(item.ID, user1.ID)
	assert.NoError(t, err)

	// 验证状态变为 deleted
	var updated model.Item
	db.First(&updated, item.ID)
	assert.Equal(t, model.ItemStatusDeleted, updated.Status)
}

func TestItemService_UpdateStatus(t *testing.T) {
	db := setupItemTestDB(t)
	itemRepo := repository.NewItemRepository(db)
	svc := NewItemService(itemRepo, nil)

	user := &model.User{Phone: "13800000001"}
	db.Create(user)

	item := &model.Item{UserID: user.ID, Title: "Test", Condition: "good", Status: model.ItemStatusActive}
	db.Create(item)

	// 下架
	err := svc.UpdateStatus(item.ID, user.ID, model.ItemStatusInactive)
	require.NoError(t, err)

	var updated model.Item
	db.First(&updated, item.ID)
	assert.Equal(t, model.ItemStatusInactive, updated.Status)

	// 重新上架
	err = svc.UpdateStatus(item.ID, user.ID, model.ItemStatusActive)
	require.NoError(t, err)

	db.First(&updated, item.ID)
	assert.Equal(t, model.ItemStatusActive, updated.Status)
}

func TestItemService_List(t *testing.T) {
	db := setupItemTestDB(t)
	itemRepo := repository.NewItemRepository(db)
	svc := NewItemService(itemRepo, nil)

	user := &model.User{Phone: "13800000001"}
	db.Create(user)

	// 创建多个物品
	for i := 0; i < 5; i++ {
		db.Create(&model.Item{
			UserID:    user.ID,
			Title:     "Item " + string(rune('A'+i)),
			Condition: "good",
			Status:    model.ItemStatusActive,
		})
	}
	// 创建一个已删除的
	db.Create(&model.Item{UserID: user.ID, Title: "Deleted", Condition: "good", Status: model.ItemStatusDeleted})

	// 默认只返回 active
	items, total, err := svc.List(repository.ItemQuery{Page: 1, PageSize: 10})
	require.NoError(t, err)
	assert.Equal(t, int64(5), total)
	assert.Len(t, items, 5)
}
