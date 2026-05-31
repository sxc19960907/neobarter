package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	jwtPkg "github.com/neobarter/server/internal/pkg/jwt"
	"github.com/stretchr/testify/assert"
)

func TestAuth_NoHeader(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	jwtManager := jwtPkg.NewManager("secret", 24)
	r.GET("/test", Auth(jwtManager), func(c *gin.Context) {
		c.JSON(200, gin.H{"ok": true})
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestAuth_InvalidFormat(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	jwtManager := jwtPkg.NewManager("secret", 24)
	r.GET("/test", Auth(jwtManager), func(c *gin.Context) {
		c.JSON(200, gin.H{"ok": true})
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "InvalidToken")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestAuth_ValidToken(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	jwtManager := jwtPkg.NewManager("secret", 24)

	var gotUserID int64
	r.GET("/test", Auth(jwtManager), func(c *gin.Context) {
		gotUserID = GetUserID(c)
		c.JSON(200, gin.H{"ok": true})
	})

	token, _ := jwtManager.GenerateToken(42, "13800138000", "personal")

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, int64(42), gotUserID)
}

func TestAuth_ExpiredToken(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	jwtManager := jwtPkg.NewManager("secret", 0) // 0 hours = expired immediately

	r.GET("/test", Auth(jwtManager), func(c *gin.Context) {
		c.JSON(200, gin.H{"ok": true})
	})

	token, _ := jwtManager.GenerateToken(1, "13800138000", "personal")

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}
