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

func setupReviewTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	err = db.AutoMigrate(&model.User{}, &model.Review{}, &model.TradeRequest{}, &model.Item{})
	require.NoError(t, err)

	return db
}

func TestReviewService_Create(t *testing.T) {
	db := setupReviewTestDB(t)
	reviewRepo := repository.NewReviewRepository(db)
	userRepo := repository.NewUserRepository(db)
	svc := NewReviewService(reviewRepo, userRepo)

	// 创建用户
	reviewer := &model.User{Phone: "13800000001", CreditScore: 100}
	reviewee := &model.User{Phone: "13800000002", CreditScore: 100}
	db.Create(reviewer)
	db.Create(reviewee)

	// 创建交易
	trade := &model.TradeRequest{InitiatorID: reviewer.ID, TargetUserID: reviewee.ID, TargetItemID: 1, Status: "completed"}
	db.Create(trade)

	// 提交好评
	review := &model.Review{
		TradeRequestID: trade.ID,
		ReviewerID:     reviewer.ID,
		RevieweeID:     reviewee.ID,
		Rating:         5,
		Comment:        "很好的交易体验",
	}
	err := svc.Create(review)
	require.NoError(t, err)

	// 验证信用积分增加
	var updatedUser model.User
	db.First(&updatedUser, reviewee.ID)
	assert.Equal(t, 102, updatedUser.CreditScore) // 好评 +2
}

func TestReviewService_Create_BadRating(t *testing.T) {
	db := setupReviewTestDB(t)
	reviewRepo := repository.NewReviewRepository(db)
	userRepo := repository.NewUserRepository(db)
	svc := NewReviewService(reviewRepo, userRepo)

	reviewer := &model.User{Phone: "13800000001", CreditScore: 100}
	reviewee := &model.User{Phone: "13800000002", CreditScore: 100}
	db.Create(reviewer)
	db.Create(reviewee)

	trade := &model.TradeRequest{InitiatorID: reviewer.ID, TargetUserID: reviewee.ID, TargetItemID: 1, Status: "completed"}
	db.Create(trade)

	// 提交差评
	review := &model.Review{
		TradeRequestID: trade.ID,
		ReviewerID:     reviewer.ID,
		RevieweeID:     reviewee.ID,
		Rating:         1,
		Comment:        "物品与描述不符",
	}
	err := svc.Create(review)
	require.NoError(t, err)

	// 验证信用积分减少
	var updatedUser model.User
	db.First(&updatedUser, reviewee.ID)
	assert.Equal(t, 98, updatedUser.CreditScore) // 差评 -2
}

func TestReviewService_Create_Duplicate(t *testing.T) {
	db := setupReviewTestDB(t)
	reviewRepo := repository.NewReviewRepository(db)
	userRepo := repository.NewUserRepository(db)
	svc := NewReviewService(reviewRepo, userRepo)

	reviewer := &model.User{Phone: "13800000001", CreditScore: 100}
	reviewee := &model.User{Phone: "13800000002", CreditScore: 100}
	db.Create(reviewer)
	db.Create(reviewee)

	trade := &model.TradeRequest{InitiatorID: reviewer.ID, TargetUserID: reviewee.ID, TargetItemID: 1, Status: "completed"}
	db.Create(trade)

	review := &model.Review{
		TradeRequestID: trade.ID,
		ReviewerID:     reviewer.ID,
		RevieweeID:     reviewee.ID,
		Rating:         5,
		Comment:        "好",
	}
	err := svc.Create(review)
	require.NoError(t, err)

	// 重复评价应该失败
	review2 := &model.Review{
		TradeRequestID: trade.ID,
		ReviewerID:     reviewer.ID,
		RevieweeID:     reviewee.ID,
		Rating:         4,
		Comment:        "还行",
	}
	err = svc.Create(review2)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "已经评价过了")
}
