package model

type SysRole struct {
	Base
	Name   string    `gorm:"size:32;not null;comment:角色名称" json:"name"`
	Code   string    `gorm:"size:32;not null;uniqueIndex;comment:角色标识" json:"code"`
	Sort   int       `gorm:"not null;default:0;comment:显示排序" json:"sort"`
	Status int8      `gorm:"not null;default:1;comment:1正常 0停用" json:"status"`
	Remark string    `gorm:"size:255;comment:备注" json:"remark"`
	Menus  []SysMenu `gorm:"many2many:sys_role_menu;" json:"menus,omitempty"`
}

func (SysRole) TableName() string { return "sys_role" }
