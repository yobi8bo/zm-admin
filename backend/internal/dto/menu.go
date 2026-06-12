package dto

type CreateMenuReq struct {
	ParentID   uint   `json:"parent_id"`
	Name       string `json:"name" binding:"required,max=32"`
	Type       int8   `json:"type" binding:"required,oneof=1 2 3"`
	Path       string `json:"path"`
	Component  string `json:"component"`
	Permission string `json:"permission"`
	Icon       string `json:"icon"`
	Sort       int    `json:"sort"`
	Visible    int8   `json:"visible"`
	Status     int8   `json:"status"`
}

type UpdateMenuReq struct {
	ParentID   uint   `json:"parent_id"`
	Name       string `json:"name" binding:"required,max=32"`
	Type       int8   `json:"type" binding:"required,oneof=1 2 3"`
	Path       string `json:"path"`
	Component  string `json:"component"`
	Permission string `json:"permission"`
	Icon       string `json:"icon"`
	Sort       int    `json:"sort"`
	Visible    int8   `json:"visible"`
	Status     int8   `json:"status"`
}

type MenuResp struct {
	ID         uint       `json:"id"`
	ParentID   uint       `json:"parent_id"`
	Name       string     `json:"name"`
	Type       int8       `json:"type"`
	Path       string     `json:"path"`
	Component  string     `json:"component"`
	Permission string     `json:"permission"`
	Icon       string     `json:"icon"`
	Sort       int        `json:"sort"`
	Visible    int8       `json:"visible"`
	Status     int8       `json:"status"`
	Children   []MenuResp `json:"children,omitempty"`
}

// 动态路由菜单（返回给前端用于addRoute）
type RouteResp struct {
	ID        uint        `json:"id"`
	Name      string      `json:"name"`
	Path      string      `json:"path"`
	Component string      `json:"component"`
	Icon      string      `json:"icon"`
	Sort      int         `json:"sort"`
	Children  []RouteResp `json:"children,omitempty"`
}
