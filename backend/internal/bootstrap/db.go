package bootstrap

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
	"zhanxu-admin/backend/config"
	"zhanxu-admin/backend/internal/model"
	"zhanxu-admin/backend/pkg/logger"
)

var DB *gorm.DB

func InitDB(cfg config.DatabaseConfig) error {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local",
		cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.DBName, cfg.Charset,
	)

	logLevel := gormlogger.Info
	switch cfg.LogLevel {
	case "silent":
		logLevel = gormlogger.Silent
	case "error":
		logLevel = gormlogger.Error
	case "warn":
		logLevel = gormlogger.Warn
	}

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: gormlogger.Default.LogMode(logLevel),
	})
	if err != nil {
		return fmt.Errorf("连接数据库失败: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return err
	}
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)

	// 自动迁移表结构
	if err = db.AutoMigrate(
		&model.SysDept{},
		&model.SysUser{},
		&model.SysRole{},
		&model.SysMenu{},
		&model.SysOperationLog{},
		&model.SysLoginLog{},
	); err != nil {
		return fmt.Errorf("数据库迁移失败: %w", err)
	}

	DB = db
	logger.Info("数据库连接成功")
	return nil
}
