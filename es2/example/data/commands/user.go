package commands

import "github.com/contextgg/pkg/es"

type NewUser struct {
	es.BaseCommand

	FirstName string
	LastName  string
	Username  string
}
