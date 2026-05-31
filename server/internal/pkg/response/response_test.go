package response

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupRouter(handler gin.HandlerFunc) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/test", handler)
	return r
}

func TestSuccess(t *testing.T) {
	r := setupRouter(func(c *gin.Context) {
		Success(c, map[string]string{"name": "test"})
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp Response
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, 0, resp.Code)
	assert.Equal(t, "success", resp.Message)
}

func TestSuccessPage(t *testing.T) {
	r := setupRouter(func(c *gin.Context) {
		SuccessPage(c, []string{"a", "b"}, 100, 1, 20)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp struct {
		Code int `json:"code"`
		Data struct {
			List     []string `json:"list"`
			Total    int64    `json:"total"`
			Page     int      `json:"page"`
			PageSize int      `json:"page_size"`
		} `json:"data"`
	}
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, 0, resp.Code)
	assert.Equal(t, int64(100), resp.Data.Total)
	assert.Equal(t, 1, resp.Data.Page)
	assert.Equal(t, 20, resp.Data.PageSize)
	assert.Len(t, resp.Data.List, 2)
}

func TestBadRequest(t *testing.T) {
	r := setupRouter(func(c *gin.Context) {
		BadRequest(c, "参数错误")
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var resp Response
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, 40000, resp.Code)
	assert.Equal(t, "参数错误", resp.Message)
}

func TestUnauthorized(t *testing.T) {
	r := setupRouter(func(c *gin.Context) {
		Unauthorized(c, "未登录")
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestNotFound(t *testing.T) {
	r := setupRouter(func(c *gin.Context) {
		NotFound(c, "不存在")
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}
