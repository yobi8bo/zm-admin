package model

import "time"

type SysUser struct {
	Base
	DeptID    uint      `gorm:"not null;default:0;comment:所属部门ID" json:"dept_id"`
	Username  string    `gorm:"size:32;not null;uniqueIndex;comment:登录用户名" json:"username"`
	Nickname  string    `gorm:"size:32;not null;default:'';comment:用户昵称" json:"nickname"`
	Password  string    `gorm:"size:128;not null;comment:bcrypt加密密码" json:"-"`
	Avatar    string    `gorm:"size:255;comment:头像URL" json:"avatar"`
	Email     string    `gorm:"size:64;comment:邮箱" json:"email"`
	Phone     string    `gorm:"size:16;comment:联系电话" json:"phone"`
	Gender    int8      `gorm:"not null;default:0;comment:0未知 1男 2女" json:"gender"`
	Status    int8      `gorm:"not null;default:1;comment:1正常 0禁用" json:"status"`
	LastLogin *time.Time `gorm:"comment:最后登录时间" json:"last_login"`
	Remark    string    `gorm:"size:255;comment:备注" json:"remark"`

	Dept  *SysDept  `gorm:"foreignKey:DeptID" json:"dept,omitempty"`
	Roles []SysRole `gorm:"many2many:sys_user_role;" json:"roles,omitempty"`
}

func (SysUser) TableName() string { return "sys_user" }
