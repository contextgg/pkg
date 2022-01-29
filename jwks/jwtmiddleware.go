package jwks

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"time"

	"github.com/contextgg/pkg/x"
	"github.com/golang-jwt/jwt/v4"
)

type key int

const userKey key = 1

func TokenFromContext(ctx context.Context) interface{} {
	return ctx.Value(userKey)
}
func ClaimsFromContext(ctx context.Context) interface{} {
	token := ctx.Value(userKey).(*jwt.Token)
	if token == nil {
		return nil
	}
	return token.Claims
}

// Algorithms
const (
	AlgorithmHS256 = "HS256"
)

// Errors
var (
	ErrJWTMissing = x.NewHTTPError(http.StatusBadRequest, "missing or malformed jwt")
	ErrJWTInvalid = x.NewHTTPError(http.StatusUnauthorized, "invalid or expired jwt")
)

// jwtFromHeader returns a `jwtExtractor` that extracts token from the request header.
func jwtFromHeader(header string, authScheme string) func(r *http.Request) string {
	return func(r *http.Request) string {
		auth := r.Header.Get(header)
		l := len(authScheme)
		if len(auth) > l+1 && auth[:l] == authScheme {
			return auth[l+1:]
		}
		return ""
	}
}

func SetToken(ctx context.Context, token interface{}) context.Context {
	return context.WithValue(ctx, userKey, token)
}

// NewJWTMiddleware
func NewJWTMiddleware(config JWTConfig, required bool) func(next http.Handler) http.Handler {
	if config.KeyFunc == nil {
		config.KeyFunc = config.keyFunc
	}

	extractor := jwtFromHeader("Authorization", "Bearer")

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			auth := extractor(r)
			if len(auth) == 0 && required {
				x.WriteError(w, ErrJWTMissing)
				return
			}

			if len(auth) > 0 {
				token, err := config.parseToken(auth)
				if err != nil {
					x.WriteError(w, err)
					return
				}
				ctx = SetToken(ctx, token)
			}

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

type JWTConfig struct {
	SigningKey interface{}
	Claims     jwt.Claims
	KeyFunc    jwt.Keyfunc
}

func (config *JWTConfig) parseToken(auth string) (interface{}, error) {
	token := new(jwt.Token)
	var err error
	// Issue #647, #656
	if _, ok := config.Claims.(jwt.MapClaims); ok {
		token, err = jwt.Parse(auth, config.KeyFunc)
	} else {
		t := reflect.ValueOf(config.Claims).Type().Elem()
		claims := reflect.New(t).Interface().(jwt.Claims)
		token, err = jwt.ParseWithClaims(auth, claims, config.KeyFunc)
	}
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, errors.New("invalid token")
	}
	return token, nil
}

// defaultKeyFunc returns a signing key of the given token.
func (config *JWTConfig) keyFunc(t *jwt.Token) (interface{}, error) {
	// Check the signing method
	if t.Method.Alg() != AlgorithmHS256 {
		return nil, fmt.Errorf("unexpected jwt signing method=%v", t.Header["alg"])
	}
	return config.SigningKey, nil
}

func NewJWTConfig(signingKey string, jwksUri string, claims jwt.Claims) JWTConfig {
	var keyfunc jwt.Keyfunc

	if len(jwksUri) > 0 {
		jwksClient := NewJwksClient(jwksUri, time.Hour, 12*time.Hour)
		keyfunc = jwksClient.KeyFunc
	}

	return JWTConfig{
		Claims:     claims,
		SigningKey: []byte(signingKey),
		KeyFunc:    keyfunc,
	}
}
