package bootstrap

import (
	"fmt"

	"github.com/spf13/viper"
	"zhanxu-admin/backend/config"
)

func InitConfig(cfgFile string) (*config.Config, error) {
	viper.SetConfigFile(cfgFile)
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("读取配置文件失败: %w", err)
	}

	var cfg config.Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("解析配置文件失败: %w", err)
	}
	return &cfg, nil
}
