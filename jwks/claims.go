package jwks

import "github.com/golang-jwt/jwt/v4"

// Identity for the
type Identity struct {
	UserId     string                 `json:"user_id" validate:"required"`
	Username   string                 `json:"username" validate:"required"`
	Connection string                 `json:"provider" validate:"required"`
	Metadata   map[string]interface{} `json:"metadata"`
}

// UserClaims are custom claims extending default ones.
type UserClaims struct {
	jwt.RegisteredClaims

	// The connection the user used to authenticate
	Connection string `json:"connection" validate:"required"`
	// The username
	Username string `json:"username" validate:"required"`
	// The roles the user is in
	Roles []string `json:"roles" validate:"required"`
	// If the user is registered
	Registered bool `json:"registered"`

	// The avatar
	AvatarUrl *string `json:"avatar_url" validate:"required"`
	// The display name
	DisplayName *string `json:"display_name" validate:"required"`
	// The email masked
	EmailMasked *string `json:"email_masked" validate:"required"`
	// If the user is verified
	Verified *bool `json:"verified,omitempty"`
	// Extra metadata
	Metadata map[string]interface{} `json:"metadata"`
	// The identities for the user
	Identities []Identity `json:"identities" validate:"required"`
}
