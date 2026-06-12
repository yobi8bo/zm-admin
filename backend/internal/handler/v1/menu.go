package v1

import (
	"github.com/gin-gonic/gin"
	"zhanxu-admin/backend/internal/dto"
	"zhanxu-admin/backend/internal/service"
	"zhanxu-admin/backend/pkg/response"
)

type MenuHandler struct{ menuSvc *service.MenuService }

func NewMenuHandler(menuSvc *service.MenuService) *MenuHandler {
	return &MenuHandler{menuSvc: menuSvc}
}

func (h *MenuHandler) List(c *gin.Context) {
	list, err := h.menuSvc.List()
	if err != nil {
		response.ServerError(c)
		return
	}
	response.Success(c, list)
}

func (h *MenuHandler) Get(c *gin.Context) {
	var param dto.IDParam
	if err := c.ShouldBindUri(&param); err != nil {
		response.FailWithMsg(c, response.CodeBadRequest, err.Error())
		return
	}
	menu, err := h.menuSvc.Get(param.ID)
	if err != nil {
		handleBizError(c, err)
		return
	}
	response.Success(c, menu)
}

func (h *MenuHandler) Create(c *gin.Context) {
	var req dto.CreateMenuReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMsg(c, response.CodeBadRequest, err.Error())
		return
	}
	if err := h.menuSvc.Create(&req); err != nil {
		handleBizError(c, err)
		return
	}
	response.Success(c, nil)
}

func (h *MenuHandler) Update(c *gin.Context) {
	var param dto.IDParam
	if err := c.ShouldBindUri(&param); err != nil {
		response.FailWithMsg(c, response.CodeBadRequest, err.Error())
		return
	}
	var req dto.UpdateMenuReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMsg(c, response.CodeBadRequest, err.Error())
		return
	}
	if err := h.menuSvc.Update(param.ID, &req); err != nil {
		handleBizError(c, err)
		return
	}
	response.Success(c, nil)
}

func (h *MenuHandler) Delete(c *gin.Context) {
	var param dto.IDParam
	if err := c.ShouldBindUri(&param); err != nil {
		response.FailWithMsg(c, response.CodeBadRequest, err.Error())
		return
	}
	if err := h.menuSvc.Delete(param.ID); err != nil {
		handleBizError(c, err)
		return
	}
	response.Success(c, nil)
}
