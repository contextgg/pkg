package storage

import (
	"context"
	"io"
)

type FileStorage interface {
	WriteChunk(ctx context.Context, id string, offset int64, src io.Reader) (int64, error)
	GetReader(ctx context.Context, id string) (io.ReadCloser, error)
	GetMetadata(ctx context.Context, id string) (map[string]string, error)
	FinishUpload(ctx context.Context, id string, metadata map[string]string) error
}
