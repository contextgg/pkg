package x

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/form3tech-oss/jwt-go"
	"github.com/form3tech-oss/jwt-go/request"
	client "github.com/ory/kratos-client-go"
	"golang.org/x/sync/semaphore"
)

// Claims for stuff
type Claims struct {
	*jwt.StandardClaims
	Session client.Session `json:"session"`
}

// JwksClient interface
type JwksClient interface {
	GetSignatureKey(keyId string) (*JSONWebKey, error)
	GetEncryptionKey(keyId string) (*JSONWebKey, error)
	GetKey(keyId string, use string) (jwk *JSONWebKey, err error)
	ParseFromRequestWithClaims(r *http.Request) (*jwt.Token, error)
}

type cacheEntry struct {
	jwk     *JSONWebKey
	refresh int64
}

type jwksClient struct {
	source  JwksSource
	cache   Cache
	refresh time.Duration
	sem     *semaphore.Weighted
}

func (c *jwksClient) GetSignatureKey(keyId string) (*JSONWebKey, error) {
	return c.GetKey(keyId, "sig")
}

func (c *jwksClient) GetEncryptionKey(keyId string) (*JSONWebKey, error) {
	return c.GetKey(keyId, "enc")
}

func (c *jwksClient) GetKey(keyId string, use string) (jwk *JSONWebKey, err error) {
	val, found := c.cache.Get(keyId)
	if found {
		entry := val.(*cacheEntry)
		if time.Now().After(time.Unix(entry.refresh, 0)) && c.sem.TryAcquire(1) {
			go func() {
				defer c.sem.Release(1)
				if _, err := c.refreshKey(keyId, use); err != nil {
					log.Printf("unable to refresh key: %v", err)
				}
			}()
		}
		return entry.jwk, nil
	} else {
		return c.refreshKey(keyId, use)
	}
}

func (c *jwksClient) ParseFromRequestWithClaims(r *http.Request) (*jwt.Token, error) {
	return request.ParseFromRequestWithClaims(r, request.OAuth2Extractor, &Claims{}, c.keyFunc)
}

func (c *jwksClient) refreshKey(keyId string, use string) (*JSONWebKey, error) {
	jwk, err := c.fetchJSONWebKey(keyId, use)
	if err != nil {
		return nil, err
	}

	c.save(keyId, jwk)
	return jwk, nil
}

func (c *jwksClient) save(keyId string, jwk *JSONWebKey) {
	c.cache.Set(keyId, &cacheEntry{
		jwk:     jwk,
		refresh: time.Now().Add(c.refresh).Unix(),
	})
}

func (c *jwksClient) keyFunc(token *jwt.Token) (interface{}, error) {
	keyId, ok := token.Header["kid"].(string)
	if !ok {
		return nil, errors.New("expecting JWT header to have string kid")
	}

	k, err := c.GetSignatureKey(keyId)
	if err != nil {
		return nil, err
	}

	return k.Key, nil
}

func (c *jwksClient) fetchJSONWebKey(keyId string, use string) (*JSONWebKey, error) {
	jsonWebKeySet, err := c.source.JSONWebKeySet()
	if err != nil {
		return nil, err
	}

	keys := jsonWebKeySet.Key(keyId)
	if len(keys) == 0 {
		return nil, fmt.Errorf("JWK is not found: %s", keyId)
	}

	for _, jwk := range keys {
		return &jwk, nil
	}
	return nil, fmt.Errorf("JWK is not found %s, use: %s", keyId, use)
}

func NewJwksClient(jwksUri string, refresh time.Duration, ttl time.Duration) JwksClient {
	if refresh >= ttl {
		panic(fmt.Sprintf("invalid refresh: %v greater or eaquals to ttl: %v", refresh, ttl))
	}
	if refresh < 0 {
		panic(fmt.Sprintf("invalid refresh: %v", refresh))
	}
	return &jwksClient{
		source:  NewJwksSource(jwksUri),
		cache:   DefaultCache(ttl),
		refresh: refresh,
		sem:     semaphore.NewWeighted(1),
	}
}
