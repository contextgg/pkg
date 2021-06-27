package x

import "context"

type User struct {
	Id string
}

type key int

const userKey key = 1

func SetUser(ctx context.Context, user *User) context.Context {
	return context.WithValue(ctx, userKey, user)
}

func FromContext(ctx context.Context) *User {
	return ctx.Value(userKey).(*User)
}
