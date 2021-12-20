package auth

import "fmt"

type InternalUser struct {
	Id             string `json:"id"`
	Username       string `json:"username"`
	PasswordHash   string `json:"password_hash"`
	EmailEncrypted []byte `json:"email_encrypted"`
	EmailHash      string `json:"email_hash"`
	EmailMasked    string `json:"email_masked"`
}

type User struct {
	Id          string `json:"id"`
	Username    string `json:"username"`
	EmailMasked string `json:"email_masked"`
	Avatar      string `json:"avatar"`
}

func ToUser(u *InternalUser) (*User, error) {
	return &User{
		Id:          u.Id,
		Username:    u.Username,
		EmailMasked: u.EmailMasked,
		Avatar:      fmt.Sprintf("https://www.gravatar.com/avatar/%s.png", u.EmailHash),
	}, nil
}
