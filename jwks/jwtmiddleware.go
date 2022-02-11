package jwks

import (
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"time"

	"github.com/contextgg/pkg/x"
	"github.com/golang-jwt/jwt/v4"
)

type key int

const tokenKey key = 1
const bearerKey key = 1

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

func signatureKeyFunc(signingKey string) jwt.Keyfunc {
	return func(t *jwt.Token) (interface{}, error) {
		// Check the signing method
		if t.Method.Alg() != AlgorithmHS256 {
			return nil, fmt.Errorf("unexpected jwt signing method=%v", t.Header["alg"])
		}
		return signingKey, nil
	}
}

// NewJWTMiddleware
func NewJWTMiddleware(config JWTConfig, required bool) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			token, err := config.GetToken(r)
			if err != nil && required {
				x.WriteError(w, ErrJWTMissing)
				return
			}
			if err == nil {
				ctx = SetToken(ctx, token)
			}

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

type JWTConfig struct {
	Claims    jwt.Claims
	KeyFunc   jwt.Keyfunc
	Extractor func(r *http.Request) string
}

func (config *JWTConfig) parseToken(auth string) (*jwt.Token, error) {
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

func (config *JWTConfig) GetToken(r *http.Request) (*jwt.Token, error) {
	auth := config.Extractor(r)
	if len(auth) == 0 {
		return nil, ErrJWTMissing
	}
	token, err := config.parseToken(auth)
	if err != nil {
		return nil, err
	}
	return token, nil
}

func NewJWTConfig(signingKey string, jwksUri string, claims jwt.Claims) (*JWTConfig, error) {
	if len(signingKey) == 0 && len(jwksUri) == 0 {
		return nil, fmt.Errorf("signing key or jwks uri is required")
	}

	c := claims
	if c == nil {
		c = &UserClaims{}
	}

	var keyfunc jwt.Keyfunc
	if len(jwksUri) > 0 {
		jwksClient := NewJwksClient(jwksUri, time.Hour, 12*time.Hour)
		keyfunc = jwksClient.KeyFunc
	} else {
		keyfunc = signatureKeyFunc(signingKey)
	}

	return &JWTConfig{
		Claims:    c,
		KeyFunc:   keyfunc,
		Extractor: jwtFromHeader("Authorization", "Bearer"),
	}, nil
}
