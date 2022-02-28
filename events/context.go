package events

import (
	"context"

	"github.com/contextgg/pkg/identity"
	"github.com/contextgg/pkg/ns"
)

func MarshalContext(ctx context.Context) map[string]interface{} {
	vals := map[string]interface{}{
		"namespace": ns.FromContext(ctx),
	}

	user, ok := identity.FromContext(ctx)
	if ok {
		vals["user"] = user
	}

	return vals
}

func UnmarshalContext(ctx context.Context, vals map[string]interface{}) context.Context {
	if vals == nil {
		return ctx
	}

	namespace, ok := vals["namespace"].(string)
	if ok && len(namespace) > 0 {
		ctx = ns.SetNamespace(ctx, namespace)
	}

	user, ok := vals["user"].(*identity.User)
	if ok && user != nil {
		ctx = identity.SetUser(ctx, user)
	}

	return ctx
}
