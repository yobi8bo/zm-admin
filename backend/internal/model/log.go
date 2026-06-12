package model

import (
	"encoding/json"

	"zhanxu-admin/backend/pkg/datetime"
)

type SysOperationLog struct {
	ID        uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID    uint   `gorm:"not null;default:0;index" json:"user_id"`
	Username  string `gorm:"size:32;not null;default:''" json:"username"`
	Module    string `gorm:"size:32;not null;default:'';comment:模块名称" json:"module"`
	Action    string `gorm:"size:32;not null;default:'';comment:操作描述" json:"action"`
	Method    string `gorm:"size:16;not null;default:''" json:"method"`
	Path      string `gorm:"size:255;not null;default:''" json:"path"`
	Query     string `gorm:"type:text;comment:Query参数" json:"query"`
	Body      string `gorm:"type:text;comment:请求Body" json:"body"`
	IP        string `gorm:"size:64;not null;default:''" json:"ip"`
	UserAgent string `gorm:"size:255" json:"user_agent"`
	Status    int    `gorm:"not null;default:0;comment:业务响应码或HTTP状态码" json:"status"`
	Latency   int64  `gorm:"not null;default:0;comment:耗时ms" json:"latency"`
	Error     string `gorm:"type:text;comment:错误信息" json:"error"`
	CreatedAt int64  `gorm:"autoCreateTime:milli" json:"created_at"`
}

func (SysOperationLog) TableName() string { return "sys_operation_log" }

func (log SysOperationLog) MarshalJSON() ([]byte, error) {
	type Alias SysOperationLog
	return json.Marshal(struct {
		Alias
		CreatedAt string `json:"created_at"`
	}{
		Alias:     Alias(log),
		CreatedAt: datetime.FormatMillis(log.CreatedAt),
	})
}

type SysLoginLog struct {
	ID        uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID    uint   `gorm:"not null;default:0;index" json:"user_id"`
	Username  string `gorm:"size:32;not null;default:''" json:"username"`
	IP        string `gorm:"size:64;not null;default:''" json:"ip"`
	Location  string `gorm:"size:64;comment:IP归属地" json:"location"`
	Browser   string `gorm:"size:64" json:"browser"`
	OS        string `gorm:"size:64" json:"os"`
	Status    int8   `gorm:"not null;default:1;comment:1成功 0失败" json:"status"`
	Message   string `gorm:"size:255;comment:失败原因" json:"message"`
	CreatedAt int64  `gorm:"autoCreateTime:milli" json:"created_at"`
}

func (SysLoginLog) TableName() string { return "sys_login_log" }

func (log SysLoginLog) MarshalJSON() ([]byte, error) {
	type Alias SysLoginLog
	return json.Marshal(struct {
		Alias
		CreatedAt string `json:"created_at"`
	}{
		Alias:     Alias(log),
		CreatedAt: datetime.FormatMillis(log.CreatedAt),
	})
}
