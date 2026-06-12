package middleware

import (
	"fmt"

	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
	"zhanxu-admin/backend/pkg/response"
)

func Casbin(enforcer *casbin.Enforcer) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := GetUserID(c)
		if userID == 0 {
			response.Fail(c, response.CodeUnauthorized)
			c.Abort()
			return
		}

		sub := fmt.Sprintf("%d", userID)
		obj := c.Request.URL.Path
		act := c.Request.Method

		allowed, err := enforcer.Enforce(sub, obj, act)
		if err != nil || !allowed {
			response.Fail(c, response.CodeForbidden)
			c.Abort()
			return
		}
		c.Next()
	}
}
