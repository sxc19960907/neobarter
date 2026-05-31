package storage

import (
	"errors"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
)

// Storage 文件存储接口
type Storage interface {
	// Upload 上传文件，返回可访问的 URL
	Upload(data []byte, filename string) (string, error)
	// Delete 删除文件
	Delete(url string) error
}

// 允许的图片类型（扩展名 -> MIME）
var allowedImageTypes = map[string]string{
	".jpg":  "image/jpeg",
	".jpeg": "image/jpeg",
	".png":  "image/png",
	".webp": "image/webp",
	".gif":  "image/gif",
}

const MaxImageSize = 5 * 1024 * 1024 // 5MB

var (
	ErrFileTooLarge   = errors.New("文件大小超过限制（最大5MB）")
	ErrInvalidType    = errors.New("不支持的文件类型，仅支持 jpg/png/webp/gif")
	ErrTypeMismatch   = errors.New("文件内容与扩展名不匹配")
	ErrEmptyFile      = errors.New("文件为空")
)

// ValidateImage 校验图片文件（扩展名 + 大小 + 内容类型）
func ValidateImage(filename string, data []byte, detectedContentType string) error {
	if len(data) == 0 {
		return ErrEmptyFile
	}
	if len(data) > MaxImageSize {
		return ErrFileTooLarge
	}

	ext := strings.ToLower(filepath.Ext(filename))
	expectedMIME, ok := allowedImageTypes[ext]
	if !ok {
		return ErrInvalidType
	}

	// 内容类型校验（http.DetectContentType 嗅探的结果）
	// gif/png/jpeg/webp 嗅探基本可靠；放宽 jpeg 的多种表述
	if !matchContentType(expectedMIME, detectedContentType) {
		return ErrTypeMismatch
	}

	return nil
}

func matchContentType(expected, detected string) bool {
	if expected == detected {
		return true
	}
	// jpeg 嗅探可能返回 image/jpeg，扩展名 jpg/jpeg 都映射到 image/jpeg
	return false
}

// GenerateFilename 生成存储文件名：YYYYMM/uuid.ext
// 不信任客户端原始文件名，只取扩展名，避免路径穿越
func GenerateFilename(originalName string) string {
	ext := strings.ToLower(filepath.Ext(originalName))
	yearMonth := time.Now().Format("200601")
	return fmt.Sprintf("%s/%s%s", yearMonth, uuid.NewString(), ext)
}
