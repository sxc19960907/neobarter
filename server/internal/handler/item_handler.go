package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	"github.com/neobarter/server/internal/middleware"
	"github.com/neobarter/server/internal/model"
	"github.com/neobarter/server/internal/pkg/response"
	"github.com/neobarter/server/internal/repository"
	"github.com/neobarter/server/internal/service"
	"github.com/shopspring/decimal"
)

type ItemHandler struct {
	itemSvc *service.ItemService
}

func NewItemHandler(itemSvc *service.ItemService) *ItemHandler {
	return &ItemHandler{itemSvc: itemSvc}
}

type CreateItemReq struct {
	Title          string   `json:"title" binding:"required"`
	Description    string   `json:"description"`
	CategoryID     *int     `json:"category_id"`
	EstimatedValue float64  `json:"estimated_value"`
	Condition      string   `json:"condition" binding:"required"`
	Images         []string `json:"images"`
	VideoURL       string   `json:"video_url"`
	Location       string   `json:"location"`
	WantItems      []string `json:"want_items"`
}

// Create 发布物品
// @Summary      发布物品
// @Description  发布一个新的可交换物品
// @Tags         物品
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        body  body      CreateItemReq  true  "物品信息"
// @Success      200   {object}  response.Response{data=model.Item}
// @Failure      400   {object}  response.Response
// @Router       /items [post]
func (h *ItemHandler) Create(c *gin.Context) {
	userID := middleware.GetUserID(c)
	var req CreateItemReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	item := &model.Item{
		UserID:         userID,
		Title:          req.Title,
		Description:    req.Description,
		CategoryID:     req.CategoryID,
		EstimatedValue: decimal.NewFromFloat(req.EstimatedValue),
		Condition:      req.Condition,
		Images:         pq.StringArray(req.Images),
		VideoURL:       req.VideoURL,
		Location:       req.Location,
		WantItems:      pq.StringArray(req.WantItems),
	}

	if err := h.itemSvc.Create(item); err != nil {
		response.ServerError(c, "发布物品失败")
		return
	}

	response.Success(c, item)
}

// Get 获取物品详情
// @Summary      获取物品详情
// @Tags         物品
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      int  true  "物品ID"
// @Success      200  {object}  response.Response{data=model.Item}
// @Failure      404  {object}  response.Response
// @Router       /items/{id} [get]
func (h *ItemHandler) Get(c *gin.Context) {
	id := parseID(c, "id")
	if id == 0 {
		return
	}

	item, err := h.itemSvc.GetByID(id)
	if err != nil {
		response.NotFound(c, "物品不存在")
		return
	}

	response.Success(c, item)
}

// List 获取物品列表
// @Summary      获取物品列表
// @Description  支持分类、关键词、地区、成色、估值范围筛选；mine=true 查询自己的物品
// @Tags         物品
// @Produce      json
// @Security     BearerAuth
// @Param        page         query     int     false  "页码"
// @Param        page_size    query     int     false  "每页数量"
// @Param        category_id  query     int     false  "分类ID"
// @Param        keyword      query     string  false  "关键词"
// @Param        location     query     string  false  "地区"
// @Param        condition    query     string  false  "成色"
// @Param        min_value    query     number  false  "最低估值"
// @Param        max_value    query     number  false  "最高估值"
// @Param        sort_by      query     string  false  "排序：created_at/view_count/estimated_value"
// @Param        mine         query     bool    false  "是否只看自己的物品"
// @Success      200  {object}  response.Response{data=response.PageData}
// @Router       /items [get]
func (h *ItemHandler) List(c *gin.Context) {
	page, pageSize := parsePage(c)
	categoryID, _ := strconv.Atoi(c.Query("category_id"))
	minValue, _ := strconv.ParseFloat(c.Query("min_value"), 64)
	maxValue, _ := strconv.ParseFloat(c.Query("max_value"), 64)

	q := repository.ItemQuery{
		CategoryID: categoryID,
		Status:     c.Query("status"),
		Keyword:    c.Query("keyword"),
		Location:   c.Query("location"),
		Condition:  c.Query("condition"),
		MinValue:   minValue,
		MaxValue:   maxValue,
		SortBy:     c.Query("sort_by"),
		Page:       page,
		PageSize:   pageSize,
	}

	// 如果查询自己的物品
	if c.Query("mine") == "true" {
		q.UserID = middleware.GetUserID(c)
		q.Status = "" // 显示所有状态
	}

	items, total, err := h.itemSvc.List(q)
	if err != nil {
		response.ServerError(c, "获取物品列表失败")
		return
	}

	response.SuccessPage(c, items, total, page, pageSize)
}

// Update 更新物品
// @Summary      更新物品信息
// @Tags         物品
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id    path      int            true  "物品ID"
// @Param        body  body      CreateItemReq  true  "物品信息"
// @Success      200   {object}  response.Response{data=model.Item}
// @Failure      403   {object}  response.Response
// @Failure      404   {object}  response.Response
// @Router       /items/{id} [put]
func (h *ItemHandler) Update(c *gin.Context) {
	userID := middleware.GetUserID(c)
	id := parseID(c, "id")
	if id == 0 {
		return
	}

	item, err := h.itemSvc.GetByID(id)
	if err != nil {
		response.NotFound(c, "物品不存在")
		return
	}
	if item.UserID != userID {
		response.Forbidden(c, "无权修改此物品")
		return
	}

	var req CreateItemReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	item.Title = req.Title
	item.Description = req.Description
	item.CategoryID = req.CategoryID
	item.EstimatedValue = decimal.NewFromFloat(req.EstimatedValue)
	item.Condition = req.Condition
	if req.Images != nil {
		item.Images = pq.StringArray(req.Images)
	}
	item.VideoURL = req.VideoURL
	item.Location = req.Location
	if req.WantItems != nil {
		item.WantItems = pq.StringArray(req.WantItems)
	}

	if err := h.itemSvc.Update(item); err != nil {
		response.ServerError(c, "更新物品失败")
		return
	}

	response.Success(c, item)
}

// Delete 删除物品
// @Summary      删除物品
// @Tags         物品
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      int  true  "物品ID"
// @Success      200  {object}  response.Response
// @Failure      403  {object}  response.Response
// @Router       /items/{id} [delete]
func (h *ItemHandler) Delete(c *gin.Context) {
	userID := middleware.GetUserID(c)
	id := parseID(c, "id")
	if id == 0 {
		return
	}

	if err := h.itemSvc.Delete(id, userID); err != nil {
		response.Forbidden(c, err.Error())
		return
	}

	response.Success(c, nil)
}

type UpdateStatusReq struct {
	Status string `json:"status" binding:"required"`
}

// UpdateStatus 上架/下架物品
// @Summary      修改物品状态
// @Description  active=上架, inactive=下架
// @Tags         物品
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id    path      int              true  "物品ID"
// @Param        body  body      UpdateStatusReq  true  "状态"
// @Success      200   {object}  response.Response
// @Failure      403   {object}  response.Response
// @Router       /items/{id}/status [put]
func (h *ItemHandler) UpdateStatus(c *gin.Context) {
	userID := middleware.GetUserID(c)
	id := parseID(c, "id")
	if id == 0 {
		return
	}

	var req UpdateStatusReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	if err := h.itemSvc.UpdateStatus(id, userID, req.Status); err != nil {
		response.Forbidden(c, err.Error())
		return
	}

	response.Success(c, nil)
}

// ListCategories 获取物品分类列表
// @Summary      获取分类列表
// @Tags         分类
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  response.Response{data=[]model.Category}
// @Router       /categories [get]
func (h *ItemHandler) ListCategories(c *gin.Context) {
	categories, err := h.itemSvc.ListCategories()
	if err != nil {
		response.ServerError(c, "获取分类失败")
		return
	}
	response.Success(c, categories)
}
