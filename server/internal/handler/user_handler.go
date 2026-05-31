package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/neobarter/server/internal/middleware"
	"github.com/neobarter/server/internal/model"
	"github.com/neobarter/server/internal/pkg/response"
	"github.com/neobarter/server/internal/service"
)

// Address handlers on UserHandler

func (h *UserHandler) ListAddresses(c *gin.Context) {
	userID := middleware.GetUserID(c)
	addresses, err := h.userSvc.ListAddresses(userID)
	if err != nil {
		response.ServerError(c, "获取地址列表失败")
		return
	}
	response.Success(c, addresses)
}

type CreateAddressReq struct {
	Name      string `json:"name" binding:"required"`
	Phone     string `json:"phone" binding:"required"`
	Province  string `json:"province" binding:"required"`
	City      string `json:"city" binding:"required"`
	District  string `json:"district" binding:"required"`
	Detail    string `json:"detail" binding:"required"`
	IsDefault bool   `json:"is_default"`
}

func (h *UserHandler) CreateAddress(c *gin.Context) {
	userID := middleware.GetUserID(c)
	var req CreateAddressReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	addr := &model.UserAddress{
		UserID:    userID,
		Name:      req.Name,
		Phone:     req.Phone,
		Province:  req.Province,
		City:      req.City,
		District:  req.District,
		Detail:    req.Detail,
		IsDefault: req.IsDefault,
	}

	if err := h.userSvc.CreateAddress(addr); err != nil {
		response.ServerError(c, "创建地址失败")
		return
	}

	response.Success(c, addr)
}

func (h *UserHandler) UpdateAddress(c *gin.Context) {
	userID := middleware.GetUserID(c)
	id := parseID(c, "id")
	if id == 0 {
		return
	}

	addr, err := h.userSvc.GetAddress(id, userID)
	if err != nil {
		response.NotFound(c, "地址不存在")
		return
	}

	var req CreateAddressReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	addr.Name = req.Name
	addr.Phone = req.Phone
	addr.Province = req.Province
	addr.City = req.City
	addr.District = req.District
	addr.Detail = req.Detail
	addr.IsDefault = req.IsDefault

	if err := h.userSvc.UpdateAddress(addr); err != nil {
		response.ServerError(c, "更新地址失败")
		return
	}

	response.Success(c, addr)
}

func (h *UserHandler) DeleteAddress(c *gin.Context) {
	userID := middleware.GetUserID(c)
	id := parseID(c, "id")
	if id == 0 {
		return
	}

	if err := h.userSvc.DeleteAddress(id, userID); err != nil {
		response.ServerError(c, "删除地址失败")
		return
	}

	response.Success(c, nil)
}

// WalletHandler

type WalletHandler struct {
	walletSvc *service.WalletService
}

func NewWalletHandler(walletSvc *service.WalletService) *WalletHandler {
	return &WalletHandler{walletSvc: walletSvc}
}

func (h *WalletHandler) GetWallet(c *gin.Context) {
	userID := middleware.GetUserID(c)
	wallet, err := h.walletSvc.GetWallet(userID)
	if err != nil {
		response.NotFound(c, "钱包不存在")
		return
	}
	response.Success(c, wallet)
}

func (h *WalletHandler) ListTransactions(c *gin.Context) {
	userID := middleware.GetUserID(c)
	page, pageSize := parsePage(c)

	transactions, total, err := h.walletSvc.ListTransactions(userID, page, pageSize)
	if err != nil {
		response.ServerError(c, "获取流水失败")
		return
	}

	response.SuccessPage(c, transactions, total, page, pageSize)
}
