package auth

import "context"

type Query interface {
	GetById(ctx context.Context, id string) (*InternalUser, error)
	GetByUsername(ctx context.Context, username string) (*InternalUser, error)
}
