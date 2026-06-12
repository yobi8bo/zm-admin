package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"zm-project/backend/internal/bootstrap"
	v1 "zm-project/backend/internal/handler/v1"
	"zm-project/backend/internal/repository"
	"zm-project/backend/internal/router"
	"zm-project/backend/internal/service"
	"zm-project/backend/pkg/logger"
)

func main() {
	// 1. 加载配置
	cfg, err := bootstrap.InitConfig("config/config.yaml")
	if err != nil {
		fmt.Printf("加载配置失败: %v\n", err)
		os.Exit(1)
	}

	// 2. 初始化日志
	logger.Init(cfg.Log)
	defer logger.Sync()

	// 3. 初始化数据库
	if err = bootstrap.InitDB(cfg.Database); err != nil {
		logger.Fatal("初始化数据库失败")
		os.Exit(1)
	}

	// 4. 初始化 Redis
	if err = bootstrap.InitRedis(cfg.Redis); err != nil {
		logger.Fatal("初始化Redis失败")
		os.Exit(1)
	}

	// 5. 初始化 Casbin
	if err = bootstrap.InitCasbin(); err != nil {
		logger.Fatal("初始化Casbin失败")
		os.Exit(1)
	}

	// 6. 初始化定时任务
	bootstrap.InitCron()
	defer bootstrap.Cron.Stop()

	// 7. 依赖注入
	db := bootstrap.DB
	enforcer := bootstrap.Enforcer

	userRepo := repository.NewUserRepo(db)
	roleRepo := repository.NewRoleRepo(db)
	menuRepo := repository.NewMenuRepo(db)
	deptRepo := repository.NewDeptRepo(db)
	logRepo := repository.NewLogRepo(db)

	authSvc := service.NewAuthService(cfg, userRepo, logRepo)
	userSvc := service.NewUserService(userRepo, roleRepo, menuRepo, enforcer)
	roleSvc := service.NewRoleService(roleRepo, enforcer)
	menuSvc := service.NewMenuService(menuRepo)
	deptSvc := service.NewDeptService(deptRepo)

	h := v1.NewHandler(authSvc, userSvc, roleSvc, menuSvc, deptSvc, logRepo)

	// 8. 注册路由
	r := router.Init(cfg, enforcer, h, logRepo)

	// 9. 启动服务（优雅关闭）
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Server.Port),
		Handler: r,
	}

	go func() {
		logger.Info(fmt.Sprintf("服务启动，监听端口 :%d", cfg.Server.Port))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal(fmt.Sprintf("服务启动失败: %v", err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("正在关闭服务...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatal(fmt.Sprintf("服务关闭失败: %v", err))
	}
	logger.Info("服务已关闭")
}
