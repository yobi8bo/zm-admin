package v1

import (
	"github.com/gin-gonic/gin"
	"zhanxu-admin/backend/internal/repository"
	"zhanxu-admin/backend/pkg/pagination"
	"zhanxu-admin/backend/pkg/response"
)

type LogHandler struct{ logRepo *repository.LogRepo }

func NewLogHandler(logRepo *repository.LogRepo) *LogHandler {
	return &LogHandler{logRepo: logRepo}
}

func (h *LogHandler) OperationList(c *gin.Context) {
	var q pagination.Query
	_ = c.ShouldBindQuery(&q)
	q.Normalize()
	list, total, err := h.logRepo.ListOperationLog(q.Page, q.PageSize, map[string]any{})
	if err != nil {
		response.ServerError(c)
		return
	}
	response.SuccessPage(c, list, total, q.Page, q.PageSize)
}

func (h *LogHandler) ClearOperationLog(c *gin.Context) {
	if err := h.logRepo.ClearOperationLog(); err != nil {
		response.ServerError(c)
		return
	}
	response.Success(c, nil)
}

func (h *LogHandler) LoginList(c *gin.Context) {
	var q pagination.Query
	_ = c.ShouldBindQuery(&q)
	q.Normalize()
	list, total, err := h.logRepo.ListLoginLog(q.Page, q.PageSize, map[string]any{})
	if err != nil {
		response.ServerError(c)
		return
	}
	response.SuccessPage(c, list, total, q.Page, q.PageSize)
}

func (h *LogHandler) ClearLoginLog(c *gin.Context) {
	if err := h.logRepo.ClearLoginLog(); err != nil {
		response.ServerError(c)
		return
	}
	response.Success(c, nil)
}
