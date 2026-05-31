package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/neobarter/server/internal/pkg/response"
	"github.com/neobarter/server/internal/repository"
	"github.com/neobarter/server/internal/service"
)

type SearchHandler struct {
	searchSvc *service.SearchService
}

func NewSearchHandler(searchSvc *service.SearchService) *SearchHandler {
	return &SearchHandler{searchSvc: searchSvc}
}

// Search 搜索物品
func (h *SearchHandler) Search(c *gin.Context) {
	page, pageSize := parsePage(c)
	categoryID, _ := strconv.Atoi(c.Query("category_id"))
	minValue, _ := strconv.ParseFloat(c.Query("min_value"), 64)
	maxValue, _ := strconv.ParseFloat(c.Query("max_value"), 64)

	q := repository.SearchQuery{
		Keyword:    c.Query("keyword"),
		CategoryID: categoryID,
		Condition:  c.Query("condition"),
		MinValue:   minValue,
		MaxValue:   maxValue,
		Location:   c.Query("location"),
		SortBy:     c.Query("sort_by"),
		Page:       page,
		PageSize:   pageSize,
	}

	result, err := h.searchSvc.Search(q)
	if err != nil {
		response.ServerError(c, "搜索失败")
		return
	}

	response.Success(c, gin.H{
		"items": result.Items,
		"total": result.Total,
		"page":  page,
		"page_size": pageSize,
	})
}

// Suggest 搜索建议
func (h *SearchHandler) Suggest(c *gin.Context) {
	prefix := c.Query("q")
	suggestions, err := h.searchSvc.Suggest(prefix)
	if err != nil {
		response.ServerError(c, "获取建议失败")
		return
	}
	response.Success(c, suggestions)
}
