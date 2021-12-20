package session

import (
	"fmt"
	"net/http"

	"github.com/contextgg/pkg/cookier"
	"github.com/google/uuid"
)

var NotFound = fmt.Errorf("Not found")

type Manager interface {
	New() *Session
	Get(r *http.Request) (*Session, error)
	Save(w http.ResponseWriter, s *Session) error
}

type manager struct {
	cm cookier.Manager
}

func (m *manager) Get(r *http.Request) (*Session, error) {
	if r == nil {
		return nil, fmt.Errorf("request is nil")
	}

	// session?
	sess := new(Session)
	err := m.cm.GetCookieValue(r, sess)
	if err == http.ErrNoCookie {
		return nil, NotFound
	}
	if err != nil {
		return nil, err
	}
	return sess, nil
}
func (m *manager) New() *Session {
	return &Session{
		Id: uuid.NewString(),
	}
}

func (m *manager) Save(w http.ResponseWriter, s *Session) error {
	return m.cm.StoreCookie(w, s)
}

func NewManager(cm cookier.Manager) Manager {
	return &manager{
		cm: cm,
	}
}
