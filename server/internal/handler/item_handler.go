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

func (h *ItemHandler) ListCategories(c *gin.Context) {
	categories, err := h.itemSvc.ListCategories()
	if err != nil {
		response.ServerError(c, "获取分类失败")
		return
	}
	response.Success(c, categories)
}
