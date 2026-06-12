package pagination

type Query struct {
	Page     int    `form:"page"`
	PageSize int    `form:"page_size"`
	OrderBy  string `form:"order_by"`
	Order    string `form:"order"`
}

func (q *Query) Normalize() {
	if q.Page <= 0 {
		q.Page = 1
	}
	if q.PageSize <= 0 {
		q.PageSize = 20
	}
	if q.PageSize > 100 {
		q.PageSize = 100
	}
	if q.Order != "asc" && q.Order != "desc" {
		q.Order = "desc"
	}
	if q.OrderBy == "" {
		q.OrderBy = "created_at"
	}
}

func (q *Query) Offset() int {
	return (q.Page - 1) * q.PageSize
}
