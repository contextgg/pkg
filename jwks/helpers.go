package jwks

import (
	"context"

	"github.com/golang-jwt/jwt/v4"
)

func TokenFromContext(ctx context.Context) interface{} {
	return ctx.Value(tokenKey)
}
func SetToken(ctx context.Context, token interface{}) context.Context {
	return context.WithValue(ctx, tokenKey, token)
}
func ClaimsFromContext(ctx context.Context) interface{} {
	token := ctx.Value(tokenKey).(*jwt.Token)
	if token == nil {
		return nil
	}
	return token.Claims
}
func BearerFromContext(ctx context.Context) string {
	return ctx.Value(bearerKey).(string)
}
func SetBearer(ctx context.Context, token string) context.Context {
	return context.WithValue(ctx, bearerKey, token)
}
