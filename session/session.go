package session

import "github.com/google/uuid"

type Session struct {
	Id string
}

func NewSession() *Session {
	return &Session{
		Id: uuid.NewString(),
	}
}
