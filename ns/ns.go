package ns

import (
	"context"
)

type key int

const namespaceKey key = 0

const defaultNamespace = "default"

func SetNamespace(ctx context.Context, namespace string) context.Context {
	return context.WithValue(ctx, namespaceKey, namespace)
}

func FromContext(ctx context.Context) string {
	namespace, ok := ctx.Value(namespaceKey).(string)
	if ok {
		return namespace
	}
	return defaultNamespace
}
