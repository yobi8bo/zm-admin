package datetime

import (
	"testing"
	"time"
)

func TestFormat(t *testing.T) {
	value := time.Date(2026, time.June, 12, 15, 30, 45, 0, time.Local)
	if got, want := Format(value), "2026-06-12 15:30:45"; got != want {
		t.Fatalf("Format() = %q, want %q", got, want)
	}
}

func TestFormatMillis(t *testing.T) {
	value := time.Date(2026, time.June, 12, 15, 30, 45, 0, time.Local)
	if got, want := FormatMillis(value.UnixMilli()), "2026-06-12 15:30:45"; got != want {
		t.Fatalf("FormatMillis() = %q, want %q", got, want)
	}
}
