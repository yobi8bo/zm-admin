package v1

import (
	"github.com/gin-gonic/gin"
	"zm-project/backend/internal/dto"
	"zm-project/backend/internal/service"
	"zm-project/backend/pkg/response"
)

type RoleHandler struct{ roleSvc *service.RoleService }

func NewRoleHandler(roleSvc *service.RoleService) *RoleHandler {
	return &RoleHandler{roleSvc: roleSvc}
}

func (h *RoleHandler) List(c *gin.Context) {
	var req dto.RoleListReq
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
	list, total, err := h.roleSvc.List(&req)
	if err != nil {
		response.ServerError(c)
		return
	}
	response.SuccessPage(c, list, total, req.Page, req.PageSize)
}

func (h *RoleHandler) All(c *gin.Context) {
	list, err := h.roleSvc.All()
	if err != nil {
		response.ServerError(c)
		return
	}
	response.Success(c, list)
}

func (h *RoleHandler) Get(c *gin.Context) {
	var param dto.IDParam
	if err := c.ShouldBindUri(&param); err != nil {
		response.FailWithMsg(c, response.CodeBadRequest, err.Error())
		return
	}
	role, err := h.roleSvc.Get(param.ID)
	if err != nil {
		handleBizError(c, err)
		return
	}
	response.Success(c, role)
}

func (h *RoleHandler) Create(c *gin.Context) {
	var req dto.CreateRoleReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMsg(c, response.CodeBadRequest, err.Error())
		return
	}
	if err := h.roleSvc.Create(&req); err != nil {
		handleBizError(c, err)
		return
	}
	response.Success(c, nil)
}

func (h *RoleHandler) Update(c *gin.Context) {
	var param dto.IDParam
	if err := c.ShouldBindUri(&param); err != nil {
		response.FailWithMsg(c, response.CodeBadRequest, err.Error())
		return
	}
	var req dto.UpdateRoleReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMsg(c, response.CodeBadRequest, err.Error())
		return
	}
	if err := h.roleSvc.Update(param.ID, &req); err != nil {
		handleBizError(c, err)
		return
	}
	response.Success(c, nil)
}

func (h *RoleHandler) Delete(c *gin.Context) {
	var param dto.IDParam
	if err := c.ShouldBindUri(&param); err != nil {
		response.FailWithMsg(c, response.CodeBadRequest, err.Error())
		return
	}
	if err := h.roleSvc.Delete(param.ID); err != nil {
		handleBizError(c, err)
		return
	}
	response.Success(c, nil)
}

func (h *RoleHandler) UpdateStatus(c *gin.Context) {
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
	if err := h.roleSvc.UpdateStatus(param.ID, &req); err != nil {
		handleBizError(c, err)
		return
	}
	response.Success(c, nil)
}

func (h *RoleHandler) AssignMenus(c *gin.Context) {
	var param dto.IDParam
	if err := c.ShouldBindUri(&param); err != nil {
		response.FailWithMsg(c, response.CodeBadRequest, err.Error())
		return
	}
	var req dto.AssignMenusReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMsg(c, response.CodeBadRequest, err.Error())
		return
	}
	if err := h.roleSvc.AssignMenus(param.ID, &req); err != nil {
		handleBizError(c, err)
		return
	}
	response.Success(c, nil)
}

func (h *RoleHandler) GetMenuIDs(c *gin.Context) {
	var param dto.IDParam
	if err := c.ShouldBindUri(&param); err != nil {
		response.FailWithMsg(c, response.CodeBadRequest, err.Error())
		return
	}
	ids, err := h.roleSvc.GetMenuIDs(param.ID)
	if err != nil {
		handleBizError(c, err)
		return
	}
	response.Success(c, ids)
}
