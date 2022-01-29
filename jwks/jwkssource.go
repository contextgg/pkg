package jwks

import (
	"encoding/json"
	"fmt"
	"net/http"

	"gopkg.in/square/go-jose.v2"
)

// JwksSource for fetching keys
type JwksSource interface {
	JSONWebKeySet() (*jose.JSONWebKeySet, error)
}

type jwksSource struct {
	client  *http.Client
	jwksUri string
}

func (s *jwksSource) JSONWebKeySet() (*jose.JSONWebKeySet, error) {
	resp, err := s.client.Get(s.jwksUri)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("failed request, status: %d", resp.StatusCode)
	}

	jsonWebKeySet := new(jose.JSONWebKeySet)
	if err = json.NewDecoder(resp.Body).Decode(jsonWebKeySet); err != nil {
		return nil, err
	}

	return jsonWebKeySet, err
}

func NewJwksSource(jwksUri string) JwksSource {
	return &jwksSource{
		client:  new(http.Client),
		jwksUri: jwksUri,
	}
}
