package model

import (
	"encoding/json"
	"testing"
	"time"
)

func TestOperationLogMarshalJSONFormatsCreatedAt(t *testing.T) {
	createdAt := time.Date(2026, time.June, 12, 15, 30, 45, 0, time.Local)
	data, err := json.Marshal(SysOperationLog{CreatedAt: createdAt.UnixMilli()})
	if err != nil {
		t.Fatalf("json.Marshal() error = %v", err)
	}

	var value map[string]any
	if err := json.Unmarshal(data, &value); err != nil {
		t.Fatalf("json.Unmarshal() error = %v", err)
	}
	if got, want := value["created_at"], "2026-06-12 15:30:45"; got != want {
		t.Fatalf("created_at = %v, want %q", got, want)
	}
}

func TestLoginLogMarshalJSONFormatsCreatedAt(t *testing.T) {
	createdAt := time.Date(2026, time.June, 12, 15, 30, 45, 0, time.Local)
	data, err := json.Marshal(SysLoginLog{CreatedAt: createdAt.UnixMilli()})
	if err != nil {
		t.Fatalf("json.Marshal() error = %v", err)
	}

	var value map[string]any
	if err := json.Unmarshal(data, &value); err != nil {
		t.Fatalf("json.Unmarshal() error = %v", err)
	}
	if got, want := value["created_at"], "2026-06-12 15:30:45"; got != want {
		t.Fatalf("created_at = %v, want %q", got, want)
	}
}
