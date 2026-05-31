package service

import (
	"net/http"

	"github.com/neobarter/server/internal/pkg/storage"
)

type UploadService struct {
	storage storage.Storage
}

func NewUploadService(s storage.Storage) *UploadService {
	return &UploadService{storage: s}
}

// UploadImage 上传图片：校验后存储，返回可访问 URL
func (s *UploadService) UploadImage(data []byte, originalName string) (string, error) {
	// 嗅探内容类型
	detectedType := http.DetectContentType(data)

	// 校验
	if err := storage.ValidateImage(originalName, data, detectedType); err != nil {
		return "", err
	}

	// 生成安全文件名并上传
	filename := storage.GenerateFilename(originalName)
	return s.storage.Upload(data, filename)
}

// DeleteImage 删除图片
func (s *UploadService) DeleteImage(url string) error {
	return s.storage.Delete(url)
}
