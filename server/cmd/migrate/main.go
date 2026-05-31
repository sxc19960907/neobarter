package main

import (
	"fmt"
	"log"

	"github.com/neobarter/server/internal/config"
	"github.com/neobarter/server/internal/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	cfg, err := config.Load("config.yaml")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Database.Host, cfg.Database.Port, cfg.Database.User,
		cfg.Database.Password, cfg.Database.DBName, cfg.Database.SSLMode)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect database: %v", err)
	}

	log.Println("Running migrations...")

	err = db.AutoMigrate(
		&model.User{},
		&model.UserAddress{},
		&model.Wallet{},
		&model.WalletTransaction{},
		&model.Category{},
		&model.Item{},
		&model.TradeRequest{},
		&model.Conversation{},
		&model.ConversationParticipant{},
		&model.Message{},
		&model.Review{},
		&model.Notification{},
	)
	if err != nil {
		log.Fatalf("Migration failed: %v", err)
	}

	// 初始化分类数据
	var count int64
	db.Model(&model.Category{}).Count(&count)
	if count == 0 {
		categories := []model.Category{
			{Name: "数码电子", Icon: "laptop", SortOrder: 1},
			{Name: "家用电器", Icon: "home", SortOrder: 2},
			{Name: "服饰鞋包", Icon: "shopping", SortOrder: 3},
			{Name: "图书教材", Icon: "book", SortOrder: 4},
			{Name: "美妆护肤", Icon: "heart", SortOrder: 5},
			{Name: "运动户外", Icon: "sports", SortOrder: 6},
			{Name: "家居家具", Icon: "shop", SortOrder: 7},
			{Name: "母婴用品", Icon: "gift", SortOrder: 8},
			{Name: "食品饮料", Icon: "coffee", SortOrder: 9},
			{Name: "其他", Icon: "ellipsis", SortOrder: 99},
		}
		db.Create(&categories)
		log.Println("Seeded categories")
	}

	log.Println("Migration completed successfully!")
}
