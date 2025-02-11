package storage

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3Storage struct {
	client *s3.Client
	bucket string
	prefix string
}

func NewS3Storage(bucket, prefix string) *S3Storage {
	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		panic(fmt.Sprintf("failed to load AWS config: %v", err))
	}

	return &S3Storage{
		client: s3.NewFromConfig(cfg),
		bucket: bucket,
		prefix: strings.TrimPrefix(prefix, "/"), // Remove leading slash
	}
}

func (s *S3Storage) Save(ctx context.Context, data []byte, path string) error {
	key := path
	if s.prefix != "" {
		key = s.prefix + "/" + path
	}

	_, err := s.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
		Body:   bytes.NewReader(data),
	})
	return err
}

func (s *S3Storage) Load(ctx context.Context, path string) ([]byte, error) {
	key := path
	if s.prefix != "" {
		key = s.prefix + "/" + path
	}

	output, err := s.client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to download from S3: %w", err)
	}
	defer output.Body.Close()

	return io.ReadAll(output.Body)
}
