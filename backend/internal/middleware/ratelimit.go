package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
	"zhanxu-admin/backend/config"
	"zhanxu-admin/backend/pkg/response"
)

func RateLimit(cfg config.RateLimitConfig) gin.HandlerFunc {
	limiter := rate.NewLimiter(rate.Limit(cfg.Rate), cfg.Burst)
	return func(c *gin.Context) {
		if !limiter.Allow() {
			c.AbortWithStatusJSON(http.StatusOK, gin.H{
				"code": response.CodeTooManyReqs,
				"msg":  response.GetMsg(response.CodeTooManyReqs),
				"data": nil,
			})
			return
		}
		c.Next()
	}
}
