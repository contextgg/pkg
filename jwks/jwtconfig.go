package jwks

import (
	"errors"
	"fmt"
	"net/http"
	"reflect"

	"github.com/contextgg/pkg/x"
	"github.com/golang-jwt/jwt/v4"
)

// Errors
var (
	ErrJWTMissing = x.NewHTTPError(http.StatusBadRequest, "missing or malformed jwt")
	ErrJWTInvalid = x.NewHTTPError(http.StatusUnauthorized, "invalid or expired jwt")
)

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

func NewJWTConfig(opts ...Option) (*JWTConfig, error) {
	cfg := &JWTConfig{
		Claims:    &UserClaims{},
		Extractor: HeaderExtractor("Authorization", "Bearer"),
	}
	for _, o := range opts {
		o(cfg)
	}

	if cfg.KeyFunc == nil {
		return nil, fmt.Errorf("KeyFunc is required")
	}
	return cfg, nil
}
