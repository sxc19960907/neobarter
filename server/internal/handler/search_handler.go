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
// @Summary      搜索物品（Elasticsearch）
// @Description  全文搜索 + 多条件筛选，返回高亮结果
// @Tags         搜索
// @Produce      json
// @Security     BearerAuth
// @Param        keyword      query     string  false  "关键词"
// @Param        category_id  query     int     false  "分类ID"
// @Param        condition    query     string  false  "成色"
// @Param        min_value    query     number  false  "最低估值"
// @Param        max_value    query     number  false  "最高估值"
// @Param        location     query     string  false  "地区"
// @Param        sort_by      query     string  false  "排序"
// @Param        page         query     int     false  "页码"
// @Param        page_size    query     int     false  "每页数量"
// @Success      200  {object}  response.Response
// @Router       /search/items [get]
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
// @Summary      搜索自动补全
// @Tags         搜索
// @Produce      json
// @Security     BearerAuth
// @Param        q    query     string  true  "前缀关键词"
// @Success      200  {object}  response.Response{data=[]string}
// @Router       /search/suggest [get]
func (h *SearchHandler) Suggest(c *gin.Context) {
	prefix := c.Query("q")
	suggestions, err := h.searchSvc.Suggest(prefix)
	if err != nil {
		response.ServerError(c, "获取建议失败")
		return
	}
	response.Success(c, suggestions)
}
