package session

import (
	"fmt"
	"net/http"

	"github.com/contextgg/pkg/cookier"
)

type Manager interface {
	Get(w http.ResponseWriter, r *http.Request) (*Session, error)
}

type manager struct {
	cm cookier.Manager
}

func (m *manager) Get(w http.ResponseWriter, r *http.Request) (*Session, error) {
	if r == nil {
		return nil, fmt.Errorf("request is nil")
	}

	// get from context first
	sess, err := FromContext(r.Context())
	if err == nil && sess != nil {
		return sess, nil
	}

	// get from cookie
	sess = NewSession()
	err = m.cm.GetCookieValue(r, sess)
	if err != nil && err != http.ErrNoCookie {
		return nil, err
	}

	// Save the cookie
	if err == http.ErrNoCookie {
		// save it in the cookie
		if err := m.cm.StoreCookie(w, sess); err != nil {
			return nil, err
		}
	}

	return sess, nil
}

func NewManager(cm cookier.Manager) Manager {
	return &manager{
		cm: cm,
	}
}
