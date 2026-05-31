package storage

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenerateFilename(t *testing.T) {
	name := GenerateFilename("photo.JPG")
	// 应为 YYYYMM/uuid.jpg 格式，扩展名小写
	assert.True(t, strings.HasSuffix(name, ".jpg"))
	assert.Contains(t, name, "/")
	// 不包含原始文件名（防路径穿越）
	assert.NotContains(t, name, "photo")
}

func TestGenerateFilename_PathTraversal(t *testing.T) {
	// 恶意文件名应被剥离，只保留扩展名
	name := GenerateFilename("../../etc/passwd.png")
	assert.True(t, strings.HasSuffix(name, ".png"))
	assert.NotContains(t, name, "..")
	assert.NotContains(t, name, "passwd")
}

func TestValidateImage_Valid(t *testing.T) {
	// PNG 文件头
	pngData := []byte{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A}
	pngData = append(pngData, make([]byte, 100)...)

	err := ValidateImage("test.png", pngData, "image/png")
	assert.NoError(t, err)
}

func TestValidateImage_EmptyFile(t *testing.T) {
	err := ValidateImage("test.png", []byte{}, "image/png")
	assert.ErrorIs(t, err, ErrEmptyFile)
}

func TestValidateImage_TooLarge(t *testing.T) {
	bigData := make([]byte, MaxImageSize+1)
	err := ValidateImage("test.png", bigData, "image/png")
	assert.ErrorIs(t, err, ErrFileTooLarge)
}

func TestValidateImage_InvalidType(t *testing.T) {
	err := ValidateImage("test.exe", []byte{1, 2, 3}, "application/octet-stream")
	assert.ErrorIs(t, err, ErrInvalidType)
}

func TestValidateImage_TypeMismatch(t *testing.T) {
	// 扩展名是 png 但内容是 jpeg
	err := ValidateImage("fake.png", []byte{1, 2, 3, 4}, "image/jpeg")
	assert.ErrorIs(t, err, ErrTypeMismatch)
}

func TestLocalProvider_UploadAndDelete(t *testing.T) {
	tmpDir := t.TempDir()
	provider := NewLocalProvider(tmpDir, "/uploads")

	data := []byte("fake image content")
	url, err := provider.Upload(data, "202605/test.png")
	require.NoError(t, err)
	assert.Equal(t, "/uploads/202605/test.png", url)

	// 验证文件确实写入
	savedPath := filepath.Join(tmpDir, "202605/test.png")
	content, err := os.ReadFile(savedPath)
	require.NoError(t, err)
	assert.Equal(t, data, content)

	// 删除
	err = provider.Delete(url)
	require.NoError(t, err)
	_, err = os.Stat(savedPath)
	assert.True(t, os.IsNotExist(err))
}

func TestLocalProvider_DeleteNonExistent(t *testing.T) {
	tmpDir := t.TempDir()
	provider := NewLocalProvider(tmpDir, "/uploads")

	// 删除不存在的文件不应报错
	err := provider.Delete("/uploads/202605/nonexistent.png")
	assert.NoError(t, err)
}

func TestLocalProvider_PathTraversalProtection(t *testing.T) {
	tmpDir := t.TempDir()
	provider := NewLocalProvider(tmpDir, "/uploads")

	// 尝试删除目录外的文件
	err := provider.Delete("/uploads/../../../etc/passwd")
	assert.Error(t, err)
}
