package ginhydra

import "github.com/ory/hydra-client-go/client/admin"

type IntrospectService interface {
	IntrospectOAuth2Token(params *admin.IntrospectOAuth2TokenParams) (*admin.IntrospectOAuth2TokenOK, error)
}
