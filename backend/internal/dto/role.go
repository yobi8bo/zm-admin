package dto

type RoleListReq struct {
	PageQuery
	Name   string `form:"name"`
	Code   string `form:"code"`
	Status *int8  `form:"status"`
}

type CreateRoleReq struct {
	Name   string `json:"name" binding:"required,max=32"`
	Code   string `json:"code" binding:"required,max=32"`
	Sort   int    `json:"sort"`
	Status int8   `json:"status"`
	Remark string `json:"remark"`
}

type UpdateRoleReq struct {
	Name   string `json:"name" binding:"required,max=32"`
	Code   string `json:"code" binding:"required,max=32"`
	Sort   int    `json:"sort"`
	Status int8   `json:"status"`
	Remark string `json:"remark"`
}

type AssignMenusReq struct {
	MenuIDs []uint `json:"menu_ids" binding:"required"`
}

type RoleResp struct {
	ID        uint   `json:"id"`
	Name      string `json:"name"`
	Code      string `json:"code"`
	Sort      int    `json:"sort"`
	Status    int8   `json:"status"`
	Remark    string `json:"remark"`
	CreatedAt string `json:"created_at"`
}
