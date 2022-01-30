package session

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/contextgg/pkg/cookier"
	"github.com/gorilla/securecookie"
)

func Test_It(t *testing.T) {
	var secret = securecookie.GenerateRandomKey(32)
	var hashKey = securecookie.GenerateRandomKey(32)
	var blockKey = securecookie.GenerateRandomKey(32)

	cookieManager := cookier.NewManager("session", hashKey, blockKey, &cookier.Options{})
	sessionManager := NewManager(cookieManager)

	csrfService, err := NewCsrfService(secret, sessionManager)
	if err != nil {
		t.Error(err)
		return
	}

	// Create a new HTTP Recorder (implements http.ResponseWriter)
	w1 := httptest.NewRecorder()
	r1 := &http.Request{}

	sess, err := sessionManager.Get(w1, r1)
	if err != nil {
		t.Error(err)
		return
	}
	if sess == nil {
		t.Error("Sess is nil")
		return
	}

	w2 := httptest.NewRecorder()
	r2 := &http.Request{
		Header: http.Header{"Cookie": w1.HeaderMap["Set-Cookie"]},
	}

	token, err := csrfService.New(w2, r2)
	if err != nil {
		t.Error(err)
		return
	}

	if err := csrfService.Verify(w2, r2, token); err != nil {
		t.Error(err)
		return
	}
}
