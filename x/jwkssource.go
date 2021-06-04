package x

import (
	"crypto/x509"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

// JSONWebKey represents a public or private key in JWK format.
type JSONWebKey struct {
	// Cryptographic key, can be a symmetric or asymmetric key.
	Key interface{}
	// Key identifier, parsed from `kid` header.
	KeyID string
	// Key algorithm, parsed from `alg` header.
	Algorithm string
	// Key use, parsed from `use` header.
	Use string

	// X.509 certificate chain, parsed from `x5c` header.
	Certificates []*x509.Certificate
	// X.509 certificate URL, parsed from `x5u` header.
	CertificatesURL *url.URL
	// X.509 certificate thumbprint (SHA-1), parsed from `x5t` header.
	CertificateThumbprintSHA1 []byte
	// X.509 certificate thumbprint (SHA-256), parsed from `x5t#S256` header.
	CertificateThumbprintSHA256 []byte
}

// JSONWebKeySet represents a JWK Set object.
type JSONWebKeySet struct {
	Keys []JSONWebKey `json:"keys"`
}

// Key convenience method returns keys by key ID. Specification states
// that a JWK Set "SHOULD" use distinct key IDs, but allows for some
// cases where they are not distinct. Hence method returns a slice
// of JSONWebKeys.
func (s *JSONWebKeySet) Key(kid string) []JSONWebKey {
	var keys []JSONWebKey
	for _, key := range s.Keys {
		if key.KeyID == kid {
			keys = append(keys, key)
		}
	}

	return keys
}

// JwksSource for fetching keys
type JwksSource interface {
	JSONWebKeySet() (*JSONWebKeySet, error)
}

type jwksSource struct {
	client  *http.Client
	jwksUri string
}

func (s *jwksSource) JSONWebKeySet() (*JSONWebKeySet, error) {
	resp, err := s.client.Get(s.jwksUri)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("failed request, status: %d", resp.StatusCode)
	}

	jsonWebKeySet := new(JSONWebKeySet)
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
