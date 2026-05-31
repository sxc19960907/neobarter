package storage

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// LocalProvider 本地文件存储（开发环境）
type LocalProvider struct {
	baseDir   string // 文件存储根目录，如 ./uploads
	urlPrefix string // 访问 URL 前缀，如 /uploads
}

func NewLocalProvider(baseDir, urlPrefix string) *LocalProvider {
	return &LocalProvider{
		baseDir:   baseDir,
		urlPrefix: strings.TrimRight(urlPrefix, "/"),
	}
}

func (p *LocalProvider) Upload(data []byte, filename string) (string, error) {
	fullPath := filepath.Join(p.baseDir, filename)

	// 创建子目录（如 uploads/202605）
	if err := os.MkdirAll(filepath.Dir(fullPath), 0o755); err != nil {
		return "", fmt.Errorf("创建目录失败: %w", err)
	}

	if err := os.WriteFile(fullPath, data, 0o644); err != nil {
		return "", fmt.Errorf("写入文件失败: %w", err)
	}

	return fmt.Sprintf("%s/%s", p.urlPrefix, filename), nil
}

func (p *LocalProvider) Delete(url string) error {
	// 从 URL 还原相对路径
	rel := strings.TrimPrefix(url, p.urlPrefix+"/")
	fullPath := filepath.Join(p.baseDir, rel)

	// 防止路径穿越：确保最终路径仍在 baseDir 内
	absBase, _ := filepath.Abs(p.baseDir)
	absTarget, _ := filepath.Abs(fullPath)
	if !strings.HasPrefix(absTarget, absBase) {
		return fmt.Errorf("非法路径")
	}

	if err := os.Remove(fullPath); err != nil && !os.IsNotExist(err) {
		return err
	}
	return nil
}

// BaseDir 返回存储根目录（用于静态文件服务注册）
func (p *LocalProvider) BaseDir() string {
	return p.baseDir
}

// URLPrefix 返回 URL 前缀
func (p *LocalProvider) URLPrefix() string {
	return p.urlPrefix
}
