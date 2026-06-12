package v1

import (
	"github.com/gin-gonic/gin"
	"zm-project/backend/pkg/response"
)

type FileHandler struct{}

func NewFileHandler() *FileHandler { return &FileHandler{} }

func (h *FileHandler) Upload(c *gin.Context) {
	// TODO: 接入 rustfs 后实现
	response.FailWithMsg(c, response.CodeServerError, "文件存储暂未配置")
}

func (h *FileHandler) Delete(c *gin.Context) {
	// TODO: 接入 rustfs 后实现
	response.FailWithMsg(c, response.CodeServerError, "文件存储暂未配置")
}
