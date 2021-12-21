package gcp

import (
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/contextgg/pkg/ns"
	"github.com/contextgg/pkg/storage"
)

const CONCURRENT_SIZE_REQUESTS = 32

type fileStorage struct {
	bucket       string
	useNamespace bool
	service      GCSAPI
}

func (store *fileStorage) WriteChunk(ctx context.Context, id string, offset int64, src io.Reader) (int64, error) {
	cid := fmt.Sprintf("%s_%d", store.keyWithPrefix(ctx, id), offset)
	objectParams := GCSObjectParams{
		Bucket: store.bucket,
		ID:     cid,
	}

	n, err := store.service.WriteObject(ctx, objectParams, src)
	if err != nil {
		return 0, err
	}

	return n, err
}

func (store *fileStorage) GetReader(ctx context.Context, id string) (io.ReadCloser, error) {
	params := GCSObjectParams{
		Bucket: store.bucket,
		ID:     store.keyWithPrefix(ctx, id),
	}

	r, err := store.service.ReadObject(ctx, params)
	if err != nil {
		return nil, err
	}

	return r, nil
}

func (store *fileStorage) GetMetadata(ctx context.Context, id string) (map[string]string, error) {
	params := GCSObjectParams{
		Bucket: store.bucket,
		ID:     store.keyWithPrefix(ctx, id),
	}

	return store.service.GetObjectMetadata(ctx, params)
}

func (store *fileStorage) FinishUpload(ctx context.Context, id string, metadata map[string]string) error {
	prefix := fmt.Sprintf("%s_", store.keyWithPrefix(ctx, id))
	filterParams := GCSFilterParams{
		Bucket: store.bucket,
		Prefix: prefix,
	}

	names, err := store.service.FilterObjects(ctx, filterParams)
	if err != nil {
		return err
	}

	composeParams := GCSComposeParams{
		Bucket:      store.bucket,
		Destination: store.keyWithPrefix(ctx, id),
		Sources:     names,
	}

	err = store.service.ComposeObjects(ctx, composeParams)
	if err != nil {
		return err
	}

	err = store.service.DeleteObjectsWithFilter(ctx, filterParams)
	if err != nil {
		return err
	}

	objectParams := GCSObjectParams{
		Bucket: store.bucket,
		ID:     store.keyWithPrefix(ctx, id),
	}

	err = store.service.SetObjectMetadata(ctx, objectParams, metadata)
	if err != nil {
		return err
	}

	return nil
}

func (store *fileStorage) keyWithPrefix(ctx context.Context, key string) string {
	var prefix string
	if store.useNamespace {
		prefix = ns.FromContext(ctx)
	}
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
func NewFileStorage(bucket string, projectId string, useNamespace bool, service GCSAPI) (storage.FileStorage, error) {
	// todo create bucket?
	ctx := context.Background()
	if err := service.CreateBucket(ctx, GCSBucketParams{
		Bucket:    bucket,
		ProjectId: projectId,
	}); err != nil {
		return nil, err
	}

	return &fileStorage{
		bucket:       bucket,
		service:      service,
		useNamespace: useNamespace,
	}, nil
}
