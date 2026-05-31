package storage

import (
	"log"

	"github.com/neobarter/server/internal/config"
)

// New 根据配置创建存储 provider
// provider = "aliyun" 时使用 OSS，否则使用本地存储
func New(cfg config.OSSConfig) Storage {
	if cfg.Provider == "aliyun" && cfg.Endpoint != "" && cfg.AccessKeyID != "" {
		ossProvider, err := NewOSSProvider(cfg.Endpoint, cfg.AccessKeyID, cfg.AccessKeySecret, cfg.Bucket)
		if err != nil {
			log.Printf("WARNING: OSS init failed, falling back to local storage: %v", err)
			return NewLocalProvider("./uploads", "/uploads")
		}
		log.Println("Using Aliyun OSS storage")
		return ossProvider
	}

	log.Println("Using local file storage (./uploads)")
	return NewLocalProvider("./uploads", "/uploads")
}
