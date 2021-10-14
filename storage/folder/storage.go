package folder

import (
	"context"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"

	"github.com/contextgg/pkg/storage"
)

var defaultFilePerm = os.FileMode(0664)

type fileStorage struct {
	abs string
}

func (store *fileStorage) WriteChunk(ctx context.Context, id string, offset int64, src io.Reader) (int64, error) {
	key := fmt.Sprintf("%s_%d", id, offset)
	path := store.keyWithPrefix(key)
	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, defaultFilePerm)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	n, err := io.Copy(file, src)
	return n, err
}
func (store *fileStorage) GetReader(ctx context.Context, id string) (io.ReadCloser, error) {
	path := store.keyWithPrefix(id)
	return os.Open(path)
}
func (store *fileStorage) FinishUpload(ctx context.Context, id string, metadata map[string]string) error {
	prefix := fmt.Sprintf("%s_", store.keyWithPrefix(id))
	var chunks []string

	// walk it?
	if err := filepath.Walk(store.abs, func(path string, info fs.FileInfo, err error) error {
		if info.IsDir() || !strings.HasPrefix(path, prefix) {
			return nil
		}
		chunks = append(chunks, path)
		return nil
	}); err != nil {
		return err
	}

	if len(chunks) == 0 {
		return nil
	}

	file, err := os.OpenFile(store.keyWithPrefix(id), os.O_CREATE|os.O_WRONLY|os.O_APPEND, defaultFilePerm)
	if err != nil {
		return err
	}
	defer file.Close()

	sort.Strings(chunks)

	for _, p := range chunks {
		chunk, err := os.OpenFile(p, os.O_RDONLY, defaultFilePerm)
		if err != nil {
			return err
		}
		defer chunk.Close()

		if _, err := io.Copy(file, chunk); err != nil {
			return err
		}
	}

	for _, p := range chunks {
		if err := os.Remove(p); err != nil {
			return err
		}
	}
	return nil
}
func (store *fileStorage) keyWithPrefix(key string) string {
	return path.Join(store.abs, key)
}

func NewFileStorage(path string, remake bool) (storage.FileStorage, error) {
	abs, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}

	if remake {
		os.RemoveAll(abs)
	}

	if err := os.MkdirAll(abs, os.ModePerm); err != nil {
		return nil, err
	}

	return &fileStorage{
		abs: abs,
	}, nil
}
