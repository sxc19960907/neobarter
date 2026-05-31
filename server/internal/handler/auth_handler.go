package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/neobarter/server/internal/middleware"
	"github.com/neobarter/server/internal/pkg/response"
	"github.com/neobarter/server/internal/service"
)

type AuthHandler struct {
	authSvc *service.AuthService
}

func NewAuthHandler(authSvc *service.AuthService) *AuthHandler {
	return &AuthHandler{authSvc: authSvc}
}

type SendCodeReq struct {
	Phone string `json:"phone" binding:"required"`
}

func (h *AuthHandler) SendCode(c *gin.Context) {
	var req SendCodeReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "手机号不能为空")
		return
	}

	if len(req.Phone) != 11 {
		response.BadRequest(c, "手机号格式错误")
		return
	}

	if err := h.authSvc.SendCode(req.Phone); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Success(c, nil)
}

type LoginReq struct {
	Phone    string `json:"phone" binding:"required"`
	Code     string `json:"code" binding:"required"`
	UserType string `json:"user_type"` // personal / enterprise
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	if req.UserType == "" {
		req.UserType = "personal"
	}

	token, user, err := h.authSvc.Login(req.Phone, req.Code, req.UserType)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Success(c, gin.H{
		"token": token,
		"user":  user,
	})
}

// UserHandler

type UserHandler struct {
	userSvc *service.UserService
}

func NewUserHandler(userSvc *service.UserService) *UserHandler {
	return &UserHandler{userSvc: userSvc}
}

func (h *UserHandler) GetMe(c *gin.Context) {
	userID := middleware.GetUserID(c)
	user, err := h.userSvc.GetByID(userID)
	if err != nil {
		response.NotFound(c, "用户不存在")
		return
	}
	response.Success(c, user)
}

type UpdateUserReq struct {
	Nickname string `json:"nickname"`
	AvatarURL string `json:"avatar_url"`
	Bio      string `json:"bio"`
	Location string `json:"location"`
}

func (h *UserHandler) UpdateMe(c *gin.Context) {
	userID := middleware.GetUserID(c)
	user, err := h.userSvc.GetByID(userID)
	if err != nil {
		response.NotFound(c, "用户不存在")
		return
	}

	var req UpdateUserReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	if req.Nickname != "" {
		user.Nickname = req.Nickname
	}
	if req.AvatarURL != "" {
		user.AvatarURL = req.AvatarURL
	}
	if req.Bio != "" {
		user.Bio = req.Bio
	}
	if req.Location != "" {
		user.Location = req.Location
	}

	if err := h.userSvc.Update(user); err != nil {
		response.ServerError(c, "更新失败")
		return
	}

	response.Success(c, user)
}

func (h *UserHandler) GetUser(c *gin.Context) {
	id := parseID(c, "id")
	if id == 0 {
		return
	}

	user, err := h.userSvc.GetByID(id)
	if err != nil {
		response.NotFound(c, "用户不存在")
		return
	}
	response.Success(c, user)
}
