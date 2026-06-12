package router

import (
	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
	"zhanxu-admin/backend/config"
	"zhanxu-admin/backend/internal/handler/v1"
	"zhanxu-admin/backend/internal/middleware"
	"zhanxu-admin/backend/internal/repository"
)

func Init(cfg *config.Config, enforcer *casbin.Enforcer, h *v1.Handler, logRepo *repository.LogRepo) *gin.Engine {
	gin.SetMode(cfg.Server.Mode)

	r := gin.New()
	r.Use(
		middleware.Cors(),
		middleware.Recovery(),
		middleware.Logger(),
		middleware.RateLimit(cfg.RateLimit),
	)

	// 公开路由（无需鉴权）
	public := r.Group("/api/v1")
	{
		public.POST("/auth/login",
			middleware.RateLimitByIP("login", cfg.RateLimit.LoginRate, cfg.RateLimit.LoginBurst),
			h.Auth.Login,
		)
		public.POST("/auth/refresh", h.Auth.RefreshToken)
		public.GET("/auth/captcha",
			middleware.RateLimitByIP("captcha", cfg.RateLimit.CaptchaRate, cfg.RateLimit.CaptchaBurst),
			h.Auth.GetCaptcha,
		)
	}

	// 需要鉴权的路由
	authed := r.Group("/api/v1")
	authed.Use(middleware.Auth(cfg.Server.JwtSecret), middleware.OperationLog(logRepo))
	{
		authed.POST("/auth/logout", h.Auth.Logout)

		// 当前用户
		authed.GET("/user/me", h.User.GetMe)
		authed.PUT("/user/me", h.User.UpdateMe)
		authed.PUT("/user/me/password",
			middleware.RateLimitByIP("password", cfg.RateLimit.SensitiveRate, cfg.RateLimit.SensitiveBurst),
			h.User.UpdateMyPassword,
		)
		authed.GET("/user/me/menus", h.User.GetMyMenus)
		authed.GET("/user/me/permissions", h.User.GetMyPermissions)

		// 需要权限校验的路由
		permRoutes := authed.Group("")
		permRoutes.Use(middleware.Casbin(enforcer))
		{
			// 用户管理
			permRoutes.GET("/users", h.User.List)
			permRoutes.POST("/users", h.User.Create)
			permRoutes.GET("/users/:id", h.User.Get)
			permRoutes.PUT("/users/:id", h.User.Update)
			permRoutes.DELETE("/users/:id", h.User.Delete)
			permRoutes.PUT("/users/:id/status", h.User.UpdateStatus)
			permRoutes.PUT("/users/:id/password",
				middleware.RateLimitByIP("password", cfg.RateLimit.SensitiveRate, cfg.RateLimit.SensitiveBurst),
				h.User.ResetPassword,
			)
			permRoutes.PUT("/users/:id/roles", h.User.AssignRoles)

			// 角色管理
			permRoutes.GET("/roles", h.Role.List)
			permRoutes.GET("/roles/all", h.Role.All)
			permRoutes.POST("/roles", h.Role.Create)
			permRoutes.GET("/roles/:id", h.Role.Get)
			permRoutes.PUT("/roles/:id", h.Role.Update)
			permRoutes.DELETE("/roles/:id", h.Role.Delete)
			permRoutes.PUT("/roles/:id/status", h.Role.UpdateStatus)
			permRoutes.PUT("/roles/:id/menus", h.Role.AssignMenus)
			permRoutes.GET("/roles/:id/menus", h.Role.GetMenuIDs)

			// 菜单管理
			permRoutes.GET("/menus", h.Menu.List)
			permRoutes.POST("/menus", h.Menu.Create)
			permRoutes.GET("/menus/:id", h.Menu.Get)
			permRoutes.PUT("/menus/:id", h.Menu.Update)
			permRoutes.DELETE("/menus/:id", h.Menu.Delete)

			// 部门管理
			permRoutes.GET("/depts", h.Dept.List)
			permRoutes.POST("/depts", h.Dept.Create)
			permRoutes.GET("/depts/:id", h.Dept.Get)
			permRoutes.PUT("/depts/:id", h.Dept.Update)
			permRoutes.DELETE("/depts/:id", h.Dept.Delete)

			// 日志
			permRoutes.GET("/logs/operation", h.Log.OperationList)
			permRoutes.DELETE("/logs/operation", h.Log.ClearOperationLog)
			permRoutes.GET("/logs/login", h.Log.LoginList)
			permRoutes.DELETE("/logs/login", h.Log.ClearLoginLog)

			// 文件上传
			permRoutes.POST("/files/upload", h.File.Upload)
			permRoutes.DELETE("/files", h.File.Delete)
		}
	}

	return r
}
