.PHONY: docs build run test migrate

# 生成 Swagger API 文档
docs:
	swag init -g cmd/server/main.go -o docs --parseInternal

# 编译
build:
	go build -o bin/server ./cmd/server

# 运行服务
run:
	go run ./cmd/server

# 运行测试
test:
	CGO_ENABLED=1 go test ./...

# 数据库迁移
migrate:
	go run ./cmd/migrate

# 启动消费者（ES 索引同步）
consumer:
	go run ./cmd/consumer
