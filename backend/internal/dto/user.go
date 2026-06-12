package dto

type UserListReq struct {
	PageQuery
	Username string `form:"username"`
	Phone    string `form:"phone"`
	Status   *int8  `form:"status"`
	DeptID   uint   `form:"dept_id"`
}

type CreateUserReq struct {
	DeptID   uint   `json:"dept_id"`
	Username string `json:"username" binding:"required,min=3,max=32"`
	Nickname string `json:"nickname" binding:"required"`
	Password string `json:"password" binding:"required,min=6"`
	Email    string `json:"email" binding:"omitempty,email"`
	Phone    string `json:"phone"`
	Gender   int8   `json:"gender"`
	Status   int8   `json:"status"`
	Remark   string `json:"remark"`
	RoleIDs  []uint `json:"role_ids"`
}

type UpdateUserReq struct {
	DeptID   uint   `json:"dept_id"`
	Nickname string `json:"nickname" binding:"required"`
	Email    string `json:"email" binding:"omitempty,email"`
	Phone    string `json:"phone"`
	Gender   int8   `json:"gender"`
	Remark   string `json:"remark"`
}

type UpdateStatusReq struct {
	Status int8 `json:"status" binding:"oneof=0 1"`
}

type ResetPasswordReq struct {
	Password string `json:"password" binding:"required,min=6"`
}

type UpdateMyPasswordReq struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=6"`
}

type AssignRolesReq struct {
	RoleIDs []uint `json:"role_ids" binding:"required"`
}

type UpdateMeReq struct {
	Nickname string `json:"nickname"`
	Email    string `json:"email" binding:"omitempty,email"`
	Phone    string `json:"phone"`
	Gender   int8   `json:"gender"`
	Avatar   string `json:"avatar"`
}

type UserResp struct {
	ID        uint    `json:"id"`
	DeptID    uint    `json:"dept_id"`
	DeptName  string  `json:"dept_name"`
	Username  string  `json:"username"`
	Nickname  string  `json:"nickname"`
	Avatar    string  `json:"avatar"`
	Email     string  `json:"email"`
	Phone     string  `json:"phone"`
	Gender    int8    `json:"gender"`
	Status    int8    `json:"status"`
	IsAdmin   bool    `json:"is_admin"`
	RoleIDs   []uint  `json:"role_ids"`
	LastLogin *string `json:"last_login"`
	CreatedAt string  `json:"created_at"`
}
