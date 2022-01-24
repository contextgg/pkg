package auth

import (
	"context"
)

type key int

const namespaceKey key = 0

const defaultNamespace = "default"

func SetAuth(ctx context.Context, user *AuthUser) context.Context {
	return context.WithValue(ctx, namespaceKey, user)
}

func FromContext(ctx context.Context) (*AuthUser, bool) {
	user, ok := ctx.Value(namespaceKey).(*AuthUser)
	if ok {
		return user, ok
	}
	return nil, false
}
