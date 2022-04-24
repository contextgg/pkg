package gcp

import (
	"context"
	"testing"
)

func TestIt(t *testing.T) {
	s := &fileStorage{
		bucket:       "contextgg-na",
		useNamespace: false,
	}

	data := []struct {
		id  string
		out string
	}{
		{
			"demo",
			"de/demo",
		},
		{
			"logo/demo",
			"logo/de/demo",
		},
	}

	for _, d := range data {
		t.Run(d.id, func(t *testing.T) {
			ctx := context.Background()
			out := s.keyWithPrefix(ctx, d.id)
			if out != d.out {
				t.Errorf("Expected %s, got %s", d.out, out)
			}
		})
	}
}
