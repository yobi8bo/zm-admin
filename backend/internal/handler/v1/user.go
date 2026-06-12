package v1

import (
	"errors"

	"github.com/gin-gonic/gin"
	"zm-project/backend/internal/dto"
	"zm-project/backend/internal/middleware"
	"zm-project/backend/internal/service"
	"zm-project/backend/pkg/response"
)

type UserHandler struct {
	userSvc *service.UserService
}

func NewUserHandler(userSvc *service.UserService) *UserHandler {
	return &UserHandler{userSvc: userSvc}
}

func (h *UserHandler) List(c *gin.Context) {
	var req dto.UserListReq
	if err := c.ShouldBindQuery(&req); err != nil {
		response.FailWithMsg(c, response.CodeBadRequest, err.Error())
		return
	}
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 20
	}
	list, total, err := h.userSvc.List(&req)
	if err != nil {
		response.ServerError(c)
		return
	}
	response.SuccessPage(c, list, total, req.Page, req.PageSize)
}

func (h *UserHandler) Get(c *gin.Context) {
	var param dto.IDParam
	if err := c.ShouldBindUri(&param); err != nil {
		response.FailWithMsg(c, response.CodeBadRequest, err.Error())
		return
	}
	user, err := h.userSvc.Get(param.ID)
	if err != nil {
		handleBizError(c, err)
		return
	}
	response.Success(c, user)
}

func (h *UserHandler) Create(c *gin.Context) {
	var req dto.CreateUserReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMsg(c, response.CodeBadRequest, err.Error())
		return
	}
	if err := h.userSvc.Create(&req); err != nil {
		handleBizError(c, err)
		return
	}
	response.Success(c, nil)
}

func (h *UserHandler) Update(c *gin.Context) {
	var param dto.IDParam
	if err := c.ShouldBindUri(&param); err != nil {
		response.FailWithMsg(c, response.CodeBadRequest, err.Error())
		return
	}
	var req dto.UpdateUserReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMsg(c, response.CodeBadRequest, err.Error())
		return
	}
	if err := h.userSvc.Update(param.ID, &req); err != nil {
		handleBizError(c, err)
		return
	}
	response.Success(c, nil)
}

func (h *UserHandler) Delete(c *gin.Context) {
	var param dto.IDParam
	if err := c.ShouldBindUri(&param); err != nil {
		response.FailWithMsg(c, response.CodeBadRequest, err.Error())
		return
	}
	if err := h.userSvc.Delete(param.ID); err != nil {
		handleBizError(c, err)
		return
	}
	response.Success(c, nil)
}

func (h *UserHandler) UpdateStatus(c *gin.Context) {
	var param dto.IDParam
	if err := c.ShouldBindUri(&param); err != nil {
		response.FailWithMsg(c, response.CodeBadRequest, err.Error())
		return
	}
	var req dto.UpdateStatusReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMsg(c, response.CodeBadRequest, err.Error())
		return
	}
	if err := h.userSvc.UpdateStatus(param.ID, &req); err != nil {
		handleBizError(c, err)
		return
	}
	response.Success(c, nil)
}

func (h *UserHandler) ResetPassword(c *gin.Context) {
	var param dto.IDParam
	if err := c.ShouldBindUri(&param); err != nil {
		response.FailWithMsg(c, response.CodeBadRequest, err.Error())
		return
	}
	var req dto.ResetPasswordReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMsg(c, response.CodeBadRequest, err.Error())
		return
	}
	if err := h.userSvc.ResetPassword(param.ID, &req); err != nil {
		handleBizError(c, err)
		return
	}
	response.Success(c, nil)
}

func (h *UserHandler) AssignRoles(c *gin.Context) {
	var param dto.IDParam
	if err := c.ShouldBindUri(&param); err != nil {
		response.FailWithMsg(c, response.CodeBadRequest, err.Error())
		return
	}
	var req dto.AssignRolesReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMsg(c, response.CodeBadRequest, err.Error())
		return
	}
	if err := h.userSvc.AssignRoles(param.ID, &req); err != nil {
		handleBizError(c, err)
		return
	}
	response.Success(c, nil)
}

func (h *UserHandler) GetMe(c *gin.Context) {
	userID := middleware.GetUserID(c)
	user, err := h.userSvc.GetMe(userID)
	if err != nil {
		handleBizError(c, err)
		return
	}
	response.Success(c, user)
}

func (h *UserHandler) UpdateMe(c *gin.Context) {
	userID := middleware.GetUserID(c)
	var req dto.UpdateMeReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMsg(c, response.CodeBadRequest, err.Error())
		return
	}
	if err := h.userSvc.UpdateMe(userID, &req); err != nil {
		handleBizError(c, err)
		return
	}
	response.Success(c, nil)
}

func (h *UserHandler) UpdateMyPassword(c *gin.Context) {
	userID := middleware.GetUserID(c)
	var req dto.UpdateMyPasswordReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMsg(c, response.CodeBadRequest, err.Error())
		return
	}
	if err := h.userSvc.UpdateMyPassword(userID, &req); err != nil {
		handleBizError(c, err)
		return
	}
	response.Success(c, nil)
}

func (h *UserHandler) GetMyMenus(c *gin.Context) {
	userID := middleware.GetUserID(c)
	menus, err := h.userSvc.GetMyMenus(userID)
	if err != nil {
		response.ServerError(c)
		return
	}
	response.Success(c, menus)
}

func (h *UserHandler) GetMyPermissions(c *gin.Context) {
	userID := middleware.GetUserID(c)
	perms, err := h.userSvc.GetMyPermissions(userID)
	if err != nil {
		response.ServerError(c)
		return
	}
	response.Success(c, perms)
}

func handleBizError(c *gin.Context, err error) {
	var bizErr *service.BizError
	if errors.As(err, &bizErr) {
		response.Fail(c, bizErr.Code)
		return
	}
	response.ServerError(c)
}
