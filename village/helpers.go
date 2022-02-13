package village

import (
	"github.com/contextgg/pkg/identity"
	"github.com/contextgg/pkg/jwks"
	"github.com/golang-jwt/jwt/v4"
)

func ToUser(claims *jwks.UserClaims) *identity.User {
	var audience string
	if len(claims.Audience) > 0 {
		audience = claims.Audience[0]
	}

	var identities []identity.Identity
	for _, id := range claims.Identities {
		identities = append(identities, identity.Identity{
			UserId:     id.UserId,
			Username:   id.Username,
			Connection: id.Connection,
			Metadata:   id.Metadata,
		})
	}

	return &identity.User{
		Id:          claims.Subject,
		Connection:  claims.Connection,
		Username:    claims.Username,
		Roles:       claims.Roles,
		Registered:  claims.Registered,
		AvatarUrl:   claims.AvatarUrl,
		DisplayName: claims.DisplayName,
		EmailMasked: claims.EmailMasked,
		Verified:    claims.Verified,
		Metadata:    claims.Metadata,
		Identities:  identities,
		Audience:    audience,
	}
}

func ToClaims(user *identity.User) *jwks.UserClaims {
	var identities []jwks.Identity
	for _, id := range user.Identities {
		identities = append(identities, jwks.Identity{
			UserId:     id.UserId,
			Username:   id.Username,
			Connection: id.Connection,
			Metadata:   id.Metadata,
		})
	}

	return &jwks.UserClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:  user.Id,
			Audience: jwt.ClaimStrings([]string{user.Audience}),
		},

		Connection:  user.Connection,
		Username:    user.Username,
		Roles:       user.Roles,
		Registered:  user.Registered,
		AvatarUrl:   user.AvatarUrl,
		DisplayName: user.DisplayName,
		EmailMasked: user.EmailMasked,
		Verified:    user.Verified,
		Metadata:    user.Metadata,
		Identities:  identities,
	}
}
