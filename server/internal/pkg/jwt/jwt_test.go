package jwt

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenerateAndParseToken(t *testing.T) {
	manager := NewManager("test-secret-key", 24)

	token, err := manager.GenerateToken(1, "13800138000", "personal")
	require.NoError(t, err)
	assert.NotEmpty(t, token)

	claims, err := manager.ParseToken(token)
	require.NoError(t, err)
	assert.Equal(t, int64(1), claims.UserID)
	assert.Equal(t, "13800138000", claims.Phone)
	assert.Equal(t, "personal", claims.UserType)
}

func TestParseToken_Invalid(t *testing.T) {
	manager := NewManager("test-secret-key", 24)

	_, err := manager.ParseToken("invalid-token")
	assert.ErrorIs(t, err, ErrTokenInvalid)
}

func TestParseToken_Expired(t *testing.T) {
	// 创建一个过期时间为 0 小时的 manager（立即过期）
	manager := NewManager("test-secret-key", 0)

	token, err := manager.GenerateToken(1, "13800138000", "personal")
	require.NoError(t, err)

	// 等待 token 过期
	time.Sleep(time.Second)

	_, err = manager.ParseToken(token)
	assert.ErrorIs(t, err, ErrTokenExpired)
}

func TestParseToken_WrongSecret(t *testing.T) {
	manager1 := NewManager("secret-1", 24)
	manager2 := NewManager("secret-2", 24)

	token, err := manager1.GenerateToken(1, "13800138000", "personal")
	require.NoError(t, err)

	_, err = manager2.ParseToken(token)
	assert.ErrorIs(t, err, ErrTokenInvalid)
}
