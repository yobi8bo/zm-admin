package bootstrap

import (
	"github.com/casbin/casbin/v2"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"zhanxu-admin/backend/pkg/logger"
)

var Enforcer *casbin.Enforcer

func InitCasbin() error {
	adapter, err := gormadapter.NewAdapterByDB(DB)
	if err != nil {
		return err
	}

	enforcer, err := casbin.NewEnforcer("config/rbac_model.conf", adapter)
	if err != nil {
		return err
	}

	if err = enforcer.LoadPolicy(); err != nil {
		return err
	}

	Enforcer = enforcer
	logger.Info("Casbin初始化成功")
	return nil
}
