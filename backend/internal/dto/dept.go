package dto

type CreateDeptReq struct {
	ParentID uint   `json:"parent_id"`
	Name     string `json:"name" binding:"required,max=64"`
	Sort     int    `json:"sort"`
	Leader   string `json:"leader"`
	Phone    string `json:"phone"`
	Email    string `json:"email" binding:"omitempty,email"`
	Status   int8   `json:"status"`
	Remark   string `json:"remark"`
}

type UpdateDeptReq struct {
	ParentID uint   `json:"parent_id"`
	Name     string `json:"name" binding:"required,max=64"`
	Sort     int    `json:"sort"`
	Leader   string `json:"leader"`
	Phone    string `json:"phone"`
	Email    string `json:"email" binding:"omitempty,email"`
	Status   int8   `json:"status"`
	Remark   string `json:"remark"`
}

type DeptResp struct {
	ID       uint       `json:"id"`
	ParentID uint       `json:"parent_id"`
	Name     string     `json:"name"`
	Sort     int        `json:"sort"`
	Leader   string     `json:"leader"`
	Phone    string     `json:"phone"`
	Email    string     `json:"email"`
	Status   int8       `json:"status"`
	Children []DeptResp `json:"children,omitempty"`
}
