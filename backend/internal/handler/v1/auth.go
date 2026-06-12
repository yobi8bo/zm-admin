package v1

import (
	"errors"

	"github.com/gin-gonic/gin"
	"zm-project/backend/internal/dto"
	"zm-project/backend/internal/middleware"
	"zm-project/backend/internal/service"
	"zm-project/backend/pkg/response"
)

type AuthHandler struct {
	authSvc *service.AuthService
}

func NewAuthHandler(authSvc *service.AuthService) *AuthHandler {
	return &AuthHandler{authSvc: authSvc}
}

func (h *AuthHandler) GetCaptcha(c *gin.Context) {
	resp, err := h.authSvc.GetCaptcha()
	if err != nil {
		response.ServerError(c)
		return
	}
	response.Success(c, resp)
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMsg(c, response.CodeBadRequest, err.Error())
		return
	}

	resp, err := h.authSvc.Login(&req, c.ClientIP(), c.GetHeader("User-Agent"))
	if err != nil {
		var bizErr *service.BizError
		if errors.As(err, &bizErr) {
			response.Fail(c, bizErr.Code)
			return
		}
		response.FailWithMsg(c, response.CodeBadRequest, err.Error())
		return
	}
	response.Success(c, resp)
}

func (h *AuthHandler) Logout(c *gin.Context) {
	userID := middleware.GetUserID(c)
	token := extractBearerToken(c)
	_ = h.authSvc.Logout(userID, token)
	response.Success(c, nil)
}

func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req dto.RefreshTokenReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMsg(c, response.CodeBadRequest, err.Error())
		return
	}
	resp, err := h.authSvc.RefreshToken(&req)
	if err != nil {
		var bizErr *service.BizError
		if errors.As(err, &bizErr) {
			response.Fail(c, bizErr.Code)
			return
		}
		response.ServerError(c)
		return
	}
	response.Success(c, resp)
}

func extractBearerToken(c *gin.Context) string {
	auth := c.GetHeader("Authorization")
	if len(auth) > 7 && auth[:7] == "Bearer " {
		return auth[7:]
	}
	return ""
}
