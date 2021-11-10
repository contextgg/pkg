package commands

import "github.com/contextgg/pkg/es"

type NewDemo struct {
	es.BaseCommand

	Name string
}

type AddDescription struct {
	es.BaseCommand

	Description string
}
