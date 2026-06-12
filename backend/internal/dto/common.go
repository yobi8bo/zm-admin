package dto

// 通用分页请求
type PageQuery struct {
	Page     int    `form:"page"`
	PageSize int    `form:"page_size"`
	OrderBy  string `form:"order_by"`
	Order    string `form:"order"`
}

// 通用ID路径参数
type IDParam struct {
	ID uint `uri:"id" binding:"required"`
}
