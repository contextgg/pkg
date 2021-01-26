package ginhydra

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/ory/hydra-client-go/client"
	"github.com/ory/hydra-client-go/client/admin"
)

func NewClientServices(uri string) admin.ClientService {
	adminURL, _ := url.Parse(uri)
	admin := client.NewHTTPClientWithConfig(nil, &client.TransportConfig{Schemes: []string{adminURL.Scheme}, Host: adminURL.Host, BasePath: adminURL.Path})
	return admin.Admin
}

// TODO: Copied from fosite
func accessTokenFromRequest(req *http.Request) string {
	auth := req.Header.Get("Authorization")
	split := strings.SplitN(auth, " ", 2)
	if len(split) != 2 || !strings.EqualFold(split[0], "bearer") {
		// Empty string returned if there's no such parameter
		err := req.ParseForm()
		if err != nil {
			return ""
		}
		fmt.Println(req.Form.Get("access_token"))
		return req.Form.Get("access_token")
	}
	return split[1]
}

func ScopesRequired(hc IntrospectService, scopes ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		var s *string
		if len(scopes) > 0 {
			str := strings.Join(scopes, " ")
			s = &str
		}

		params := admin.NewIntrospectOAuth2TokenParamsWithContext(c)
		params.Token = accessTokenFromRequest(c.Request)
		params.Scope = s

		ok, err := hc.IntrospectOAuth2Token(params)
		if err != nil {
			c.Error(err)
			c.Abort()
			return
		}

		// All required scopes are found
		c.Set("hydra", ok)
		c.Next()
	}
}
