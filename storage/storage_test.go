package storage

import (
	"fmt"
	"strings"
	"testing"
)

func TestNewStorage(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		wantType string
	}{
		{
			name:     "local path",
			path:     "./backups",
			wantType: "LocalStorage",
		},
		{
			name:     "s3 path",
			path:     "s3://my-bucket/backups",
			wantType: "S3Storage",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewStorage(tt.path)
			if err != nil {
				t.Errorf("NewStorage() error = %v", err)
				return
			}
			if got == nil {
				t.Error("NewStorage() returned nil")
				return
			}
			gotType := fmt.Sprintf("%T", got)
			if !strings.Contains(gotType, tt.wantType) {
				t.Errorf("NewStorage() = %T, want type containing %v", got, tt.wantType)
			}
		})
	}
}
