package commands

import "github.com/contextgg/pkg/es"

type NewDemo struct {
	es.BaseCommand

	Name string
}
