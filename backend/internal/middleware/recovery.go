package middleware

import (
	"net/http"
	"runtime/debug"

	"github.com/gin-gonic/gin"
	"zhanxu-admin/backend/pkg/logger"
	"zhanxu-admin/backend/pkg/response"
	"go.uber.org/zap"
)

func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				logger.Error("panic recovered",
					zap.Any("error", err),
					zap.String("stack", string(debug.Stack())),
				)
				c.AbortWithStatusJSON(http.StatusOK, gin.H{
					"code": response.CodeServerError,
					"msg":  response.GetMsg(response.CodeServerError),
					"data": nil,
				})
			}
		}()
		c.Next()
	}
}
