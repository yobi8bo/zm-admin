package datetime

import "time"

const Layout = "2006-01-02 15:04:05"

func Format(value time.Time) string {
	if value.IsZero() {
		return ""
	}
	return value.Format(Layout)
}

func FormatPointer(value *time.Time) *string {
	if value == nil {
		return nil
	}
	formatted := Format(*value)
	return &formatted
}

func FormatMillis(value int64) string {
	if value == 0 {
		return ""
	}
	return time.UnixMilli(value).Format(Layout)
}
