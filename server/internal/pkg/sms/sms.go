package sms

import (
	"fmt"
	"math/rand"
)

// Provider 短信服务接口
type Provider interface {
	SendCode(phone string, code string) error
}

// MockProvider 开发环境模拟短信
type MockProvider struct{}

func NewMockProvider() *MockProvider {
	return &MockProvider{}
}

func (p *MockProvider) SendCode(phone string, code string) error {
	fmt.Printf("[SMS Mock] Send code %s to %s\n", code, phone)
	return nil
}

// GenerateCode 生成6位验证码
func GenerateCode() string {
	return fmt.Sprintf("%06d", rand.Intn(1000000))
}
