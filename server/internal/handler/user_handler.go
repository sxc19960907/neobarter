package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/neobarter/server/internal/middleware"
	"github.com/neobarter/server/internal/model"
	"github.com/neobarter/server/internal/pkg/response"
	"github.com/neobarter/server/internal/service"
)

// Address handlers on UserHandler

// ListAddresses 获取收货地址列表
// @Summary      获取收货地址列表
// @Tags         用户
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  response.Response{data=[]model.UserAddress}
// @Router       /users/me/addresses [get]
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

// CreateAddress 新增收货地址
// @Summary      新增收货地址
// @Tags         用户
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        body  body      CreateAddressReq  true  "地址信息"
// @Success      200   {object}  response.Response{data=model.UserAddress}
// @Failure      400   {object}  response.Response
// @Router       /users/me/addresses [post]
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

// UpdateAddress 更新收货地址
// @Summary      更新收货地址
// @Tags         用户
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id    path      int               true  "地址ID"
// @Param        body  body      CreateAddressReq  true  "地址信息"
// @Success      200   {object}  response.Response{data=model.UserAddress}
// @Failure      404   {object}  response.Response
// @Router       /users/me/addresses/{id} [put]
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

// DeleteAddress 删除收货地址
// @Summary      删除收货地址
// @Tags         用户
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      int  true  "地址ID"
// @Success      200  {object}  response.Response
// @Router       /users/me/addresses/{id} [delete]
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

type VerifyRealNameReq struct {
	RealName string `json:"real_name" binding:"required"`
	IDCard   string `json:"id_card" binding:"required"`
}

// VerifyRealName 提交实名认证
// @Summary      提交实名认证
// @Description  提交真实姓名+身份证号（MVP 提交即认证，未接三方核验）
// @Tags         用户
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        body  body      VerifyRealNameReq  true  "实名信息"
// @Success      200   {object}  response.Response
// @Failure      400   {object}  response.Response
// @Router       /users/me/verify-realname [post]
func (h *UserHandler) VerifyRealName(c *gin.Context) {
	userID := middleware.GetUserID(c)
	var req VerifyRealNameReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请填写姓名和身份证号")
		return
	}

	if err := h.userSvc.VerifyRealName(userID, req.RealName, req.IDCard); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Success(c, nil)
}

type VerifyEnterpriseReq struct {
	EnterpriseName string `json:"enterprise_name" binding:"required"`
	LicenseURL     string `json:"license_url" binding:"required"`
}

// VerifyEnterprise 提交企业认证
// @Summary      提交企业认证
// @Description  提交企业名称+营业执照URL（仅企业用户，MVP 提交即认证）
// @Tags         用户
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        body  body      VerifyEnterpriseReq  true  "企业信息"
// @Success      200   {object}  response.Response
// @Failure      400   {object}  response.Response
// @Router       /users/me/verify-enterprise [post]
func (h *UserHandler) VerifyEnterprise(c *gin.Context) {
	userID := middleware.GetUserID(c)
	var req VerifyEnterpriseReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请填写企业名称并上传营业执照")
		return
	}

	if err := h.userSvc.VerifyEnterprise(userID, req.EnterpriseName, req.LicenseURL); err != nil {
		response.BadRequest(c, err.Error())
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

// GetWallet 获取钱包信息
// @Summary      获取巴特币钱包
// @Tags         钱包
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  response.Response{data=model.Wallet}
// @Failure      404  {object}  response.Response
// @Router       /wallet [get]
func (h *WalletHandler) GetWallet(c *gin.Context) {
	userID := middleware.GetUserID(c)
	wallet, err := h.walletSvc.GetWallet(userID)
	if err != nil {
		response.NotFound(c, "钱包不存在")
		return
	}
	response.Success(c, wallet)
}

// ListTransactions 获取钱包流水
// @Summary      获取巴特币流水
// @Tags         钱包
// @Produce      json
// @Security     BearerAuth
// @Param        page       query     int  false  "页码"
// @Param        page_size  query     int  false  "每页数量"
// @Success      200  {object}  response.Response{data=response.PageData}
// @Router       /wallet/transactions [get]
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
