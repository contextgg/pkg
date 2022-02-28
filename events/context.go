package events

import (
	"context"
	"encoding/json"

	"github.com/contextgg/pkg/identity"
	"github.com/contextgg/pkg/ns"
)

func MarshalContext(ctx context.Context) map[string]string {
	vals := map[string]string{
		"namespace": ns.FromContext(ctx),
	}

	if user, ok := identity.FromContext(ctx); ok {
		if out, err := json.Marshal(user); err != nil {
			vals["user"] = string(out)
		}
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

	user, ok := vals["user"]
	if ok && len(user) > 0 {
		u := identity.User{}
		if err := json.Unmarshal([]byte(user), &u); err != nil {
			ctx = identity.SetUser(ctx, &u)
		}
	}

	return ctx
}
