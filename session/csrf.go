package session

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

const MIN_KEY_LENGTH = 32
const TSLen = 19

type CsrfService interface {
	New(w http.ResponseWriter, r *http.Request) (string, error)
	Verify(w http.ResponseWriter, r *http.Request, token string) error
}

type csrfService struct {
	secret []byte

	sessionManager Manager
}

func (s *csrfService) createCsrf(ts time.Time, sessionId string) (string, error) {
	str := fmt.Sprintf("%0*d", TSLen, ts.UTC().Unix())

	mac := hmac.New(sha256.New, s.secret)
	mac.Write([]byte(str))
	mac.Write([]byte(sessionId))
	out := mac.Sum(nil)
	return str + hex.EncodeToString(out), nil
}

func (s *csrfService) verify(token, sessionId string, valid time.Duration) error {
	if len(token) < TSLen+1 {
		return fmt.Errorf("CSRF token too short")
	}
	str := token[:TSLen]
	i, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		return err
	}

	ts := time.Unix(i, 0)
	expectedToken, err := s.createCsrf(ts, sessionId)
	if err != nil {
		return fmt.Errorf("CSRF hash error")
	}
	if !hmac.Equal([]byte(token), []byte(expectedToken)) {
		return fmt.Errorf("CSRF token invalid")
	}
	if ts.Add(valid).Before(time.Now()) {
		return fmt.Errorf("CSRF token expired")
	}
	return nil
}

func (s *csrfService) New(w http.ResponseWriter, r *http.Request) (string, error) {
	session, err := s.sessionManager.Get(w, r)
	if err != nil {
		return "", err
	}

	return s.createCsrf(time.Now(), session.Id)
}

func (s *csrfService) Verify(w http.ResponseWriter, r *http.Request, token string) error {
	session, err := s.sessionManager.Get(w, r)
	if err != nil {
		return err
	}
	return s.verify(token, session.Id, 5*time.Minute)
}

func NewCsrfService(secret []byte, sessionManager Manager) (CsrfService, error) {
	if len(secret) < MIN_KEY_LENGTH {
		return nil, fmt.Errorf("Key too short")
	}

	return &csrfService{
		secret:         secret,
		sessionManager: sessionManager,
	}, nil
}
