package cookier

import (
	"net/http"

	"github.com/gorilla/securecookie"
)

type Manager interface {
	GetCookieValue(r *http.Request, out interface{}) error
	StoreCookie(w http.ResponseWriter, value interface{}) error
	DeleteCookie(w http.ResponseWriter) error
}

type manager struct {
	name    string
	sc      *securecookie.SecureCookie
	options *Options
}

func (m *manager) GetCookieValue(r *http.Request, out interface{}) error {
	// look up the cookie!
	c, err := r.Cookie(m.name)
	if err != nil {
		return err
	}

	// decode the value
	if err := m.sc.Decode(c.Name, c.Value, out); err != nil {
		return err
	}
	return nil
}
func (m *manager) StoreCookie(w http.ResponseWriter, value interface{}) error {
	e, err := m.sc.Encode(m.name, value)
	if err != nil {
		return err
	}

	// build the cookie
	http.SetCookie(w, NewCookie(m.name, e, m.options))
	return nil
}
func (m *manager) DeleteCookie(w http.ResponseWriter) error {
	opts := &Options{
		Path:     m.options.Path,
		Domain:   m.options.Domain,
		MaxAge:   -1,
		Secure:   m.options.Secure,
		HttpOnly: m.options.HttpOnly,
	}

	// to delete the cookie send one back with maxage negative
	http.SetCookie(w, NewCookie(m.name, "", opts))
	return nil
}

func NewManager(name string, hashKey, blockKey []byte, options *Options) Manager {
	sc := securecookie.New(hashKey, blockKey)
	return &manager{
		name:    name,
		sc:      sc,
		options: options,
	}
}
