package identity

import "fmt"

type SessionRequest struct {
	Method string `json:"method"`
	Path   string `json:"path"`
}

type Identity struct {
	UserId     string                 `json:"user_id" validate:"required"`
	Username   string                 `json:"username" validate:"required"`
	Connection string                 `json:"provider" validate:"required"`
	Metadata   map[string]interface{} `json:"metadata"`
}

// The user object
type User struct {
	// The token for the user
	Token *string `json:"token" validate:"required"`
	// Either cookie, server or api
	AuthType string `json:"auth_type" validate:"required"`
	// The id for the user
	Id string `json:"id" validate:"required"`
	// The connection for the user
	Connection string `json:"connection" validate:"required"`
	// The password for the user
	Username string `json:"username" validate:"required"`
	// The roles the user is in
	Roles []string `json:"roles" validate:"required"`
	// If the user is registered
	Registered bool `json:"registered"`
	// The avatar url for the user
	AvatarUrl *string `json:"avatar_url,omitempty"`
	// The display name for the user
	DisplayName *string `json:"display_name,omitempty"`
	// The password for the user
	EmailMasked *string `json:"email_masked,omitempty"`
	// If the user is verified
	Verified *bool `json:"verified,omitempty"`
	// Extra metadata
	Metadata map[string]interface{} `json:"metadata"`
	// The identities for the user
	Identities []Identity `json:"identities" validate:"required"`
	// Who the user is authenticated for
	Audience string `json:"audience" validate:"required"`
}

type ErrorMessage struct {
	// Error code
	Code int `json:"code"`

	// Error message
	Message string `json:"message"`
}

// Error makes it compatible with `error` interface.
func (he *ErrorMessage) Error() string {
	return fmt.Sprintf("code=%d, message=%v", he.Code, he.Message)
}
