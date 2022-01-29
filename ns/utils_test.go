package ns

import (
	"context"
	"net/http"
	"testing"

	"github.com/contextgg/pkg/jwks"
	"github.com/golang-jwt/jwt/v4"
)

type CustomClaims struct {
	jwt.RegisteredClaims

	Username string `json:"username"`
}

func Test_Slug(t *testing.T) {
	data := []struct {
		in  string
		out string
	}{
		{in: "localhost:3000", out: "localhost"},
		{in: "localhost", out: "localhost"},
		{in: "inflow.pro", out: "inflowpro"},
		{in: "contextgg.inflow.pro", out: "contextgg"},
		{in: "abc.inflow.pro", out: "abc"},
	}

	for _, d := range data {
		t.Run(d.in, func(t *testing.T) {
			r := &http.Request{
				Host: d.in,
			}

			b := Slug(r, ".inflow.pro")
			if b != d.out {
				t.Errorf("Slug mismatch %s != %s", b, d.out)
			}
		})
	}
}

func Test_Jwt(t *testing.T) {
	data := []struct {
		in  string
		out string
	}{
		{in: "localhost", out: "localhost"},
		{in: "ctx", out: "ctx"},
	}

	extract := func(c interface{}) string {
		claims, ok := c.(*CustomClaims)
		if !ok || len(claims.Audience) == 0 {
			return ""
		}
		return claims.Audience[0]
	}

	for _, d := range data {
		t.Run(d.in, func(t *testing.T) {
			token := &jwt.Token{
				Claims: &CustomClaims{
					RegisteredClaims: jwt.RegisteredClaims{
						Audience: jwt.ClaimStrings{d.in},
					},
				},
			}
			ctx := jwks.SetToken(context.Background(), token)

			r := &http.Request{}
			r = r.WithContext(ctx)

			b := JwtValue(r, extract)
			if b != d.out {
				t.Errorf("Aud mismatch %s != %s", b, d.out)
			}
		})
	}
}
