package auth

import (
	"fmt"
	"net/http"

	"github.com/contextgg/pkg/cookier"
)

var ErrNotValid = fmt.Errorf("Not valid")

type Manager interface {
	GetAuthUser(r *http.Request) (*AuthUser, error)
	StoreCookie(w http.ResponseWriter, authUser *AuthUser) error
	DeleteCookie(w http.ResponseWriter) error
}

type manager struct {
	cm cookier.Manager
}

func (a *manager) GetAuthUser(r *http.Request) (*AuthUser, error) {
	auth := new(AuthUser)
	if err := a.cm.GetCookieValue(r, auth); err != nil {
		return nil, err
	}
	return auth, nil
}

func (a *manager) StoreCookie(w http.ResponseWriter, authUser *AuthUser) error {
	if authUser == nil {
		return fmt.Errorf("No auth user provided")
	}
	return a.cm.StoreCookie(w, authUser)
}

func (a *manager) DeleteCookie(w http.ResponseWriter) error {
	return a.cm.DeleteCookie(w)
}

func NewManager(cm cookier.Manager) Manager {
	return &manager{
		cm: cm,
	}
}
