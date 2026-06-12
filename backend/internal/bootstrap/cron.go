package bootstrap

import (
	"github.com/robfig/cron/v3"
	"zhanxu-admin/backend/pkg/logger"
)

var Cron *cron.Cron

func InitCron() {
	Cron = cron.New(cron.WithSeconds())
	// 在此注册定时任务，例如：
	// Cron.AddFunc("0 0 * * * *", tasks.CleanExpiredLogs)
	Cron.Start()
	logger.Info("定时任务启动成功")
}
