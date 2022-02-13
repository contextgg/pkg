package jwks

import (
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

// Algorithms
const (
	AlgorithmHS256 = "HS256"
)

type Extractor func(r *http.Request) string

// HeaderExtractor extracts a string from a header and supports authSchemes.
func HeaderExtractor(header string, authScheme string) Extractor {
	return func(r *http.Request) string {
		auth := r.Header.Get(header)
		l := len(authScheme)
		if len(auth) > l+1 && auth[:l] == authScheme {
			return auth[l+1:]
		}
		return ""
	}
}

type Option func(c *JWTConfig)

func UseJwksClient(jwksUri string) Option {
	var keyfunc jwt.Keyfunc
	if len(jwksUri) > 0 {
		jwksClient := NewJwksClient(jwksUri, time.Hour, 12*time.Hour)
		keyfunc = jwksClient.KeyFunc
	}

	return func(c *JWTConfig) {
		c.KeyFunc = keyfunc
	}
}
func UseSigningKey(signingKey string) Option {
	var keyfunc jwt.Keyfunc
	if len(signingKey) > 0 {
		keyfunc = func(t *jwt.Token) (interface{}, error) {
			// Check the signing method
			if t.Method.Alg() != AlgorithmHS256 {
				return nil, fmt.Errorf("unexpected jwt signing method=%v", t.Header["alg"])
			}
			return signingKey, nil
		}
	}
	return func(c *JWTConfig) {
		c.KeyFunc = keyfunc
	}
}
func UseClaims(claims jwt.Claims) Option {
	return func(c *JWTConfig) {
		c.Claims = claims
	}
}
func UseExtractor(extractor Extractor) Option {
	return func(c *JWTConfig) {
		c.Extractor = extractor
	}
}
