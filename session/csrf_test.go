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

	// Create a new HTTP Recorder (implements http.ResponseWriter)
	recorder := httptest.NewRecorder()
	if err := sessionManager.Save(recorder, sessionManager.New()); err != nil {
		t.Error(err)
		return
	}

	// Copy the Cookie over to a new Request
	r := &http.Request{
		Header: http.Header{"Cookie": recorder.HeaderMap["Set-Cookie"]},
	}

	service, err := NewCsrfService(secret, sessionManager)
	if err != nil {
		t.Error(err)
		return
	}

	token, err := service.New(r)
	if err != nil {
		t.Error(err)
		return
	}

	t.Log(token)

	if err := service.Verify(r, token); err != nil {
		t.Error(err)
		return
	}
}
