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

// SendCode 发送短信验证码
// @Summary      发送验证码
// @Description  向指定手机号发送登录/注册验证码（60秒内不可重复发送）
// @Tags         认证
// @Accept       json
// @Produce      json
// @Param        body  body      SendCodeReq  true  "手机号"
// @Success      200   {object}  response.Response
// @Failure      400   {object}  response.Response
// @Router       /auth/send-code [post]
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

// Login 登录/注册
// @Summary      登录或注册
// @Description  手机号+验证码登录，未注册的手机号自动创建账户并赠送初始巴特币
// @Tags         认证
// @Accept       json
// @Produce      json
// @Param        body  body      LoginReq  true  "登录信息"
// @Success      200   {object}  response.Response  "返回 token 和用户信息"
// @Failure      400   {object}  response.Response
// @Router       /auth/login [post]
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

// GetMe 获取当前用户信息
// @Summary      获取当前用户信息
// @Tags         用户
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  response.Response{data=model.User}
// @Failure      404  {object}  response.Response
// @Router       /users/me [get]
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

// UpdateMe 更新当前用户信息
// @Summary      更新当前用户信息
// @Tags         用户
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        body  body      UpdateUserReq  true  "用户信息"
// @Success      200   {object}  response.Response{data=model.User}
// @Failure      400   {object}  response.Response
// @Router       /users/me [put]
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

// GetUser 获取指定用户公开信息
// @Summary      获取用户公开信息
// @Tags         用户
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      int  true  "用户ID"
// @Success      200  {object}  response.Response{data=model.User}
// @Failure      404  {object}  response.Response
// @Router       /users/{id} [get]
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
