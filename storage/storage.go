package storage

import (
	"context"
	"strings"
)

// Storage defines the interface for backup storage operations
type Storage interface {
	Save(ctx context.Context, data []byte, path string) error
	Load(ctx context.Context, path string) ([]byte, error)
}

// NewStorage creates a storage implementation based on the path
// Supports:
// - Local file system: path starts with "/" or "./" or is a relative path
// - S3: path starts with "s3://"
func NewStorage(path string) (Storage, error) {
	if strings.HasPrefix(path, "s3://") {
		parts := strings.SplitN(path[5:], "/", 2) // Skip "s3://" and split on first "/"
		bucket := parts[0]
		prefix := ""
		if len(parts) > 1 {
			prefix = parts[1]
			// Remove the filename from prefix
			if lastSlash := strings.LastIndex(prefix, "/"); lastSlash != -1 {
				prefix = prefix[:lastSlash]
			}
		}
		return NewS3Storage(bucket, prefix), nil
	}
	return NewLocalStorage(), nil
}
