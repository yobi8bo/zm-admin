package service

// BizError 业务错误，携带错误码
type BizError struct {
	Code int
	Msg  string
}

func (e *BizError) Error() string {
	if e.Msg != "" {
		return e.Msg
	}
	return "业务错误"
}
