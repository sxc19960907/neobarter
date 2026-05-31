package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/neobarter/server/internal/config"
	"github.com/neobarter/server/internal/handler"
	"github.com/neobarter/server/internal/middleware"
	jwtPkg "github.com/neobarter/server/internal/pkg/jwt"
	"github.com/neobarter/server/internal/pkg/sms"
	"github.com/neobarter/server/internal/repository"
	"github.com/neobarter/server/internal/service"
	"github.com/neobarter/server/internal/ws"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func main() {
	// 加载配置
	cfg, err := config.Load("config.yaml")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// 初始化数据库
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Database.Host, cfg.Database.Port, cfg.Database.User,
		cfg.Database.Password, cfg.Database.DBName, cfg.Database.SSLMode)

	var gormLogger logger.Interface
	if cfg.Server.Mode == "debug" {
		gormLogger = logger.Default.LogMode(logger.Info)
	} else {
		gormLogger = logger.Default.LogMode(logger.Silent)
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{Logger: gormLogger})
	if err != nil {
		log.Fatalf("Failed to connect database: %v", err)
	}

	sqlDB, _ := db.DB()
	sqlDB.SetMaxOpenConns(cfg.Database.MaxOpenConns)
	sqlDB.SetMaxIdleConns(cfg.Database.MaxIdleConns)

	// 初始化 Redis
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Redis.Host, cfg.Redis.Port),
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})

	// 初始化组件
	jwtManager := jwtPkg.NewManager(cfg.JWT.Secret, cfg.JWT.ExpireHours)
	smsProvider := sms.NewMockProvider()
	wsHub := ws.NewHub()
	go wsHub.Run()

	// 初始化 Repository
	userRepo := repository.NewUserRepository(db)
	walletRepo := repository.NewWalletRepository(db)
	itemRepo := repository.NewItemRepository(db)
	tradeRepo := repository.NewTradeRepository(db)
	messageRepo := repository.NewMessageRepository(db)
	reviewRepo := repository.NewReviewRepository(db)
	notificationRepo := repository.NewNotificationRepository(db)

	// 初始化 Service
	walletSvc := service.NewWalletService(walletRepo, cfg.Wallet.InitialReward)
	authSvc := service.NewAuthService(userRepo, walletSvc, rdb, jwtManager, smsProvider)
	userSvc := service.NewUserService(userRepo)
	itemSvc := service.NewItemService(itemRepo)
	tradeSvc := service.NewTradeService(tradeRepo, itemRepo, walletSvc, notificationRepo)
	messageSvc := service.NewMessageService(messageRepo, wsHub)
	reviewSvc := service.NewReviewService(reviewRepo, userRepo)
	notificationSvc := service.NewNotificationService(notificationRepo)

	// 初始化 Handler
	authHandler := handler.NewAuthHandler(authSvc)
	userHandler := handler.NewUserHandler(userSvc)
	walletHandler := handler.NewWalletHandler(walletSvc)
	itemHandler := handler.NewItemHandler(itemSvc)
	tradeHandler := handler.NewTradeHandler(tradeSvc)
	messageHandler := handler.NewMessageHandler(messageSvc)
	reviewHandler := handler.NewReviewHandler(reviewSvc)
	notificationHandler := handler.NewNotificationHandler(notificationSvc)

	// 设置路由
	if cfg.Server.Mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()
	r.Use(middleware.CORS())

	// API v1
	v1 := r.Group("/v1")
	{
		// 认证（无需登录）
		auth := v1.Group("/auth")
		{
			auth.POST("/send-code", authHandler.SendCode)
			auth.POST("/login", authHandler.Login)
		}

		// 需要登录的接口
		authorized := v1.Group("")
		authorized.Use(middleware.Auth(jwtManager))
		{
			// 用户
			users := authorized.Group("/users")
			{
				users.GET("/me", userHandler.GetMe)
				users.PUT("/me", userHandler.UpdateMe)
				users.GET("/:id", userHandler.GetUser)
				users.GET("/me/addresses", userHandler.ListAddresses)
				users.POST("/me/addresses", userHandler.CreateAddress)
				users.PUT("/me/addresses/:id", userHandler.UpdateAddress)
				users.DELETE("/me/addresses/:id", userHandler.DeleteAddress)
			}

			// 钱包
			wallet := authorized.Group("/wallet")
			{
				wallet.GET("", walletHandler.GetWallet)
				wallet.GET("/transactions", walletHandler.ListTransactions)
			}

			// 物品
			items := authorized.Group("/items")
			{
				items.POST("", itemHandler.Create)
				items.GET("", itemHandler.List)
				items.GET("/:id", itemHandler.Get)
				items.PUT("/:id", itemHandler.Update)
				items.DELETE("/:id", itemHandler.Delete)
				items.PUT("/:id/status", itemHandler.UpdateStatus)
			}

			// 分类（公开）
			v1.GET("/categories", itemHandler.ListCategories)

			// 交易
			trades := authorized.Group("/trades")
			{
				trades.POST("", tradeHandler.Create)
				trades.GET("", tradeHandler.List)
				trades.GET("/:id", tradeHandler.Get)
				trades.PUT("/:id/accept", tradeHandler.Accept)
				trades.PUT("/:id/reject", tradeHandler.Reject)
				trades.PUT("/:id/complete", tradeHandler.Complete)
				trades.PUT("/:id/cancel", tradeHandler.Cancel)
			}

			// 消息
			messages := authorized.Group("/messages")
			{
				messages.GET("/conversations", messageHandler.ListConversations)
				messages.GET("/conversations/:id", messageHandler.GetMessages)
				messages.POST("", messageHandler.Send)
				messages.PUT("/conversations/:id/read", messageHandler.MarkRead)
			}

			// 评价
			reviews := authorized.Group("/reviews")
			{
				reviews.POST("", reviewHandler.Create)
				reviews.GET("/user/:id", reviewHandler.ListByUser)
			}

			// 通知
			notifications := authorized.Group("/notifications")
			{
				notifications.GET("", notificationHandler.List)
				notifications.GET("/unread-count", notificationHandler.UnreadCount)
				notifications.PUT("/:id/read", notificationHandler.MarkRead)
				notifications.PUT("/read-all", notificationHandler.MarkAllRead)
			}
		}

		// WebSocket
		v1.GET("/ws", middleware.Auth(jwtManager), func(c *gin.Context) {
			ws.ServeWS(wsHub, c)
		})
	}

	// 启动服务
	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	log.Printf("NeoBarter server starting on %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
