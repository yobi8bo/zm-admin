package model

type SysDept struct {
	Base
	ParentID uint      `gorm:"not null;default:0;comment:父部门ID" json:"parent_id"`
	Name     string    `gorm:"size:64;not null;comment:部门名称" json:"name"`
	Sort     int       `gorm:"not null;default:0;comment:显示排序" json:"sort"`
	Leader   string    `gorm:"size:32;comment:负责人" json:"leader"`
	Phone    string    `gorm:"size:16;comment:联系电话" json:"phone"`
	Email    string    `gorm:"size:64;comment:邮箱" json:"email"`
	Status   int8      `gorm:"not null;default:1;comment:1正常 0停用" json:"status"`
	Remark   string    `gorm:"size:255;comment:备注" json:"remark"`
	Children []SysDept `gorm:"-" json:"children,omitempty"`
}

func (SysDept) TableName() string { return "sys_dept" }
