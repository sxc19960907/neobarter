package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/neobarter/server/docs"
	"github.com/neobarter/server/internal/config"
	"github.com/neobarter/server/internal/handler"
	"github.com/neobarter/server/internal/middleware"
	"github.com/neobarter/server/internal/pkg/es"
	jwtPkg "github.com/neobarter/server/internal/pkg/jwt"
	"github.com/neobarter/server/internal/pkg/mq"
	"github.com/neobarter/server/internal/pkg/sms"
	"github.com/neobarter/server/internal/pkg/storage"
	"github.com/neobarter/server/internal/repository"
	"github.com/neobarter/server/internal/scheduler"
	"github.com/neobarter/server/internal/service"
	"github.com/neobarter/server/internal/ws"
	"github.com/redis/go-redis/v9"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// @title           Easy Barter API
// @version         1.0
// @description     现代易货交易平台后端 API。支持 C2C/B2B 双模式、巴特币结算、智能匹配。
// @termsOfService  https://github.com/sxc19960907/neobarter

// @contact.name   NeoBarter Team
// @contact.url    https://github.com/sxc19960907/neobarter

// @license.name  MIT

// @host      localhost:8080
// @BasePath  /v1

// @securityDefinitions.apikey  BearerAuth
// @in                          header
// @name                        Authorization
// @description                 输入格式：Bearer {token}
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
	storageProvider := storage.New(cfg.OSS)
	wsHub := ws.NewHub()
	go wsHub.Run()

	// 初始化 Elasticsearch（可选，连接失败不阻塞启动）
	var esClient *es.Client
	if len(cfg.Elasticsearch.Addresses) > 0 {
		var err error
		esClient, err = es.NewClient(cfg.Elasticsearch.Addresses)
		if err != nil {
			log.Printf("WARNING: Elasticsearch not available: %v", err)
		} else {
			esClient.EnsureIndex()
		}
	}

	// 初始化 RabbitMQ Publisher（可选）
	var mqPublisher *mq.Publisher
	if cfg.RabbitMQ.URL != "" {
		var err error
		mqPublisher, err = mq.NewPublisher(cfg.RabbitMQ.URL)
		if err != nil {
			log.Printf("WARNING: RabbitMQ not available: %v", err)
		}
	}

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
	itemSvc := service.NewItemService(itemRepo, mqPublisher)
	tradeSvc := service.NewTradeService(tradeRepo, itemRepo, walletSvc, notificationRepo)
	messageSvc := service.NewMessageService(messageRepo, itemRepo, wsHub)
	reviewSvc := service.NewReviewService(reviewRepo, userRepo)
	notificationSvc := service.NewNotificationService(notificationRepo)
	uploadSvc := service.NewUploadService(storageProvider)

	// 后台定时任务（交易超时过期等）
	sched := scheduler.New(tradeSvc)
	schedCtx, schedCancel := context.WithCancel(context.Background())
	defer schedCancel()
	sched.Start(schedCtx)

	// 搜索服务（依赖 ES）
	var searchHandler *handler.SearchHandler
	if esClient != nil {
		searchRepo := repository.NewSearchRepository(esClient)
		searchSvc := service.NewSearchService(searchRepo)
		searchHandler = handler.NewSearchHandler(searchSvc)
	}

	// 初始化 Handler
	authHandler := handler.NewAuthHandler(authSvc)
	userHandler := handler.NewUserHandler(userSvc)
	walletHandler := handler.NewWalletHandler(walletSvc)
	itemHandler := handler.NewItemHandler(itemSvc)
	tradeHandler := handler.NewTradeHandler(tradeSvc)
	messageHandler := handler.NewMessageHandler(messageSvc)
	reviewHandler := handler.NewReviewHandler(reviewSvc)
	notificationHandler := handler.NewNotificationHandler(notificationSvc)
	uploadHandler := handler.NewUploadHandler(uploadSvc)

	// 设置路由
	if cfg.Server.Mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()
	r.Use(middleware.CORS())

	// Swagger API 文档（非 release 模式开放）
	if cfg.Server.Mode != "release" {
		r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}

	// 本地存储时挂载静态文件服务
	if local, ok := storageProvider.(*storage.LocalProvider); ok {
		r.Static(local.URLPrefix(), local.BaseDir())
	}

	// API v1
	v1 := r.Group("/v1")
	{
		// 认证（无需登录）
		auth := v1.Group("/auth")
		{
			auth.POST("/send-code", authHandler.SendCode)
			auth.POST("/login", authHandler.Login)
		}

		// 物品分类（公开，首页未登录也可浏览）
		v1.GET("/categories", itemHandler.ListCategories)

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
				users.POST("/me/verify-realname", userHandler.VerifyRealName)
				users.POST("/me/verify-enterprise", userHandler.VerifyEnterprise)
			}

			// 钱包
			wallet := authorized.Group("/wallet")
			{
				wallet.GET("", walletHandler.GetWallet)
				wallet.GET("/transactions", walletHandler.ListTransactions)
			}

			// 上传（限流：每分钟最多30次）
			authorized.POST("/upload/image",
				middleware.RateLimit(30, time.Minute),
				uploadHandler.UploadImage)

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

			// 搜索（依赖 ES）
			if searchHandler != nil {
				search := authorized.Group("/search")
				{
					search.GET("/items", searchHandler.Search)
					search.GET("/suggest", searchHandler.Suggest)
				}
			}

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
				trades.PUT("/:id/counter", tradeHandler.Counter)
				trades.PUT("/:id/counter/accept", tradeHandler.AcceptCounter)
				trades.PUT("/:id/counter/reject", tradeHandler.RejectCounter)
			}

			// 消息
			messages := authorized.Group("/messages")
			{
				messages.GET("/conversations", messageHandler.ListConversations)
				messages.GET("/conversations/:id", messageHandler.GetMessages)
				messages.POST("", messageHandler.Send)
				messages.POST("/item-card", messageHandler.SendItemCard)
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
