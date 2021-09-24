package events

import (
	"context"

	"github.com/contextgg/pkg/ns"
)

func MarshalContext(ctx context.Context) map[string]string {
	vals := map[string]string{
		"namespace": ns.FromContext(ctx),
	}
	return vals
}

func UnmarshalContext(ctx context.Context, vals map[string]string) context.Context {
	if vals == nil {
		return ctx
	}

	namespace, ok := vals["namespace"]
	if ok && len(namespace) > 0 {
		ctx = ns.SetNamespace(ctx, namespace)
	}

	return ctx
}
