package session

import (
	"context"
	"fmt"
)

type key int

const sessionKey key = 0

func SetSession(ctx context.Context, sess *Session) context.Context {
	return context.WithValue(ctx, sessionKey, sess)
}

func FromContext(ctx context.Context) (*Session, error) {
	sess, ok := ctx.Value(sessionKey).(*Session)
	if ok {
		return sess, nil
	}
	return nil, fmt.Errorf("Not found")
}
