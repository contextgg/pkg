package auth

import (
	"context"
	"fmt"
	"net/http"

	"github.com/contextgg/pkg/cookier"
	"github.com/matthewhartstonge/argon2"
)

var ErrNotValid = fmt.Errorf("Not valid")

type Manager interface {
	GetValidUser(ctx context.Context, username string, password string) (*User, error)
	GetAuthedUser(r *http.Request) (*User, error)
	StoreUserCookie(w http.ResponseWriter, user *User) error
	DeleteUserCookie(w http.ResponseWriter) error
}

type manager struct {
	cm    cookier.Manager
	query Query
}

func (a *manager) GetValidUser(ctx context.Context, username string, password string) (*User, error) {
	user, err := a.query.GetByUsername(ctx, username)
	if err != nil {
		return nil, err
	}

	// todo validate password!
	valid, err := argon2.VerifyEncoded([]byte(password), []byte(user.PasswordHash))
	if err != nil {
		return nil, err
	}
	if !valid {
		return nil, ErrNotValid
	}

	return ToUser(user)
}

func (a *manager) GetAuthedUser(r *http.Request) (*User, error) {
	var id string
	if err := a.cm.GetCookieValue(r, &id); err != nil {
		return nil, err
	}

	ctx := r.Context()
	user, err := a.query.GetById(ctx, id)
	if err != nil {
		return nil, err
	}

	return ToUser(user)
}

func (a *manager) StoreUserCookie(w http.ResponseWriter, user *User) error {
	if user == nil {
		return fmt.Errorf("No user provided")
	}
	return a.cm.StoreCookie(w, user.Id)
}

func (a *manager) DeleteUserCookie(w http.ResponseWriter) error {
	return a.cm.DeleteCookie(w)
}

func NewManager(cm cookier.Manager, query Query) Manager {
	return &manager{
		cm:    cm,
		query: query,
	}
}
