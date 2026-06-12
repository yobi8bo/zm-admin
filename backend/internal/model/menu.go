package model

type SysMenu struct {
	Base
	ParentID   uint      `gorm:"not null;default:0;comment:父菜单ID" json:"parent_id"`
	Name       string    `gorm:"size:32;not null;comment:菜单名称" json:"name"`
	Type       int8      `gorm:"not null;comment:1目录 2菜单 3按钮" json:"type"`
	Path       string    `gorm:"size:128;comment:路由path" json:"path"`
	Component  string    `gorm:"size:128;comment:前端组件路径" json:"component"`
	Permission string    `gorm:"size:64;comment:权限标识" json:"permission"`
	Icon       string    `gorm:"size:64;comment:菜单图标" json:"icon"`
	Sort       int       `gorm:"not null;default:0;comment:显示排序" json:"sort"`
	Visible    int8      `gorm:"not null;default:1;comment:1显示 0隐藏" json:"visible"`
	Status     int8      `gorm:"not null;default:1;comment:1正常 0停用" json:"status"`
	Children   []SysMenu `gorm:"-" json:"children,omitempty"`
}

func (SysMenu) TableName() string { return "sys_menu" }
