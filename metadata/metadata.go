package metadata

import (
	"context"

	"google.golang.org/grpc/metadata"
)

func FromContext(ctx context.Context) map[string]interface{} {
	m := make(map[string]interface{})

	if md, ok := metadata.FromIncomingContext(ctx); ok {
		for k, v := range md {
			if len(v) == 1 {
				m[k] = v[0]
				continue
			}
			m[k] = v
		}
	}

	// what about tracing?

	return m
}
