package storage

import (
	"context"
	"io"
)

type FileStorage interface {
	GetPath(ctx context.Context, key string) string
	WriteChunk(ctx context.Context, key string, offset int64, src io.Reader) (int64, error)
	GetReader(ctx context.Context, key string) (io.ReadCloser, error)
	GetMetadata(ctx context.Context, key string) (map[string]string, error)
	FinishUpload(ctx context.Context, key string, metadata map[string]string) error
}
