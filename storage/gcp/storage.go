package gcp

import (
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/contextgg/pkg/storage"
)

const CONCURRENT_SIZE_REQUESTS = 32

type fileStorage struct {
	Bucket       string
	ObjectPrefix string
	Service      GCSAPI
}

func (store *fileStorage) WriteChunk(ctx context.Context, id string, offset int64, src io.Reader) (int64, error) {
	cid := fmt.Sprintf("%s_%d", store.keyWithPrefix(id), offset)
	objectParams := GCSObjectParams{
		Bucket: store.Bucket,
		ID:     cid,
	}

	n, err := store.Service.WriteObject(ctx, objectParams, src)
	if err != nil {
		return 0, err
	}

	return n, err
}

func (store *fileStorage) GetReader(ctx context.Context, id string) (io.ReadCloser, error) {
	params := GCSObjectParams{
		Bucket: store.Bucket,
		ID:     store.keyWithPrefix(id),
	}

	r, err := store.Service.ReadObject(ctx, params)
	if err != nil {
		return nil, err
	}

	return r, nil
}

func (store *fileStorage) FinishUpload(ctx context.Context, id string, metadata map[string]string) error {
	prefix := fmt.Sprintf("%s_", store.keyWithPrefix(id))
	filterParams := GCSFilterParams{
		Bucket: store.Bucket,
		Prefix: prefix,
	}

	names, err := store.Service.FilterObjects(ctx, filterParams)
	if err != nil {
		return err
	}

	composeParams := GCSComposeParams{
		Bucket:      store.Bucket,
		Destination: store.keyWithPrefix(id),
		Sources:     names,
	}

	err = store.Service.ComposeObjects(ctx, composeParams)
	if err != nil {
		return err
	}

	err = store.Service.DeleteObjectsWithFilter(ctx, filterParams)
	if err != nil {
		return err
	}

	objectParams := GCSObjectParams{
		Bucket: store.Bucket,
		ID:     store.keyWithPrefix(id),
	}

	err = store.Service.SetObjectMetadata(ctx, objectParams, metadata)
	if err != nil {
		return err
	}

	return nil
}

func (store *fileStorage) keyWithPrefix(key string) string {
	prefix := store.ObjectPrefix
	if prefix != "" && !strings.HasSuffix(prefix, "/") {
		prefix += "/"
	}
	if len(key) > 1 {
		prefix += key[0:2] + "/"
	}
	return prefix + key
}

// New constructs a new GCS storage backend using the supplied GCS bucket name
// and service object.
func NewFileStorage(bucket string, projectId string, service GCSAPI) (storage.FileStorage, error) {
	// todo create bucket?
	ctx := context.Background()
	if err := service.CreateBucket(ctx, GCSBucketParams{
		Bucket:    bucket,
		ProjectId: projectId,
	}); err != nil {
		return nil, err
	}

	return &fileStorage{
		Bucket:  bucket,
		Service: service,
	}, nil
}
