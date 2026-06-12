package v1

import (
	"github.com/gin-gonic/gin"
	"zm-project/backend/internal/dto"
	"zm-project/backend/internal/service"
	"zm-project/backend/pkg/response"
)

type DeptHandler struct{ deptSvc *service.DeptService }

func NewDeptHandler(deptSvc *service.DeptService) *DeptHandler {
	return &DeptHandler{deptSvc: deptSvc}
}

func (h *DeptHandler) List(c *gin.Context) {
	list, err := h.deptSvc.List()
	if err != nil {
		response.ServerError(c)
		return
	}
	response.Success(c, list)
}

func (h *DeptHandler) Get(c *gin.Context) {
	var param dto.IDParam
	if err := c.ShouldBindUri(&param); err != nil {
		response.FailWithMsg(c, response.CodeBadRequest, err.Error())
		return
	}
	dept, err := h.deptSvc.Get(param.ID)
	if err != nil {
		handleBizError(c, err)
		return
	}
	response.Success(c, dept)
}

func (h *DeptHandler) Create(c *gin.Context) {
	var req dto.CreateDeptReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMsg(c, response.CodeBadRequest, err.Error())
		return
	}
	if err := h.deptSvc.Create(&req); err != nil {
		handleBizError(c, err)
		return
	}
	response.Success(c, nil)
}

func (h *DeptHandler) Update(c *gin.Context) {
	var param dto.IDParam
	if err := c.ShouldBindUri(&param); err != nil {
		response.FailWithMsg(c, response.CodeBadRequest, err.Error())
		return
	}
	var req dto.UpdateDeptReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMsg(c, response.CodeBadRequest, err.Error())
		return
	}
	if err := h.deptSvc.Update(param.ID, &req); err != nil {
		handleBizError(c, err)
		return
	}
	response.Success(c, nil)
}

func (h *DeptHandler) Delete(c *gin.Context) {
	var param dto.IDParam
	if err := c.ShouldBindUri(&param); err != nil {
		response.FailWithMsg(c, response.CodeBadRequest, err.Error())
		return
	}
	if err := h.deptSvc.Delete(param.ID); err != nil {
		handleBizError(c, err)
		return
	}
	response.Success(c, nil)
}
