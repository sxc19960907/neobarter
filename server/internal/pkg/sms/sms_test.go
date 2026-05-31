package sms

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateCode(t *testing.T) {
	code := GenerateCode()
	assert.Len(t, code, 6)

	// 确保生成的是数字
	for _, c := range code {
		assert.True(t, c >= '0' && c <= '9')
	}
}

func TestGenerateCode_Unique(t *testing.T) {
	codes := make(map[string]bool)
	for i := 0; i < 100; i++ {
		code := GenerateCode()
		codes[code] = true
	}
	// 100次生成应该有多个不同的值
	assert.Greater(t, len(codes), 1)
}

func TestMockProvider_SendCode(t *testing.T) {
	provider := NewMockProvider()
	err := provider.SendCode("13800138000", "123456")
	assert.NoError(t, err)
}
