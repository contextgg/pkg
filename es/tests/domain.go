package tests

import (
	"github.com/contextgg/pkg/es"

	"github.com/contextgg/pkg/es/tests/aggregates"
	"github.com/contextgg/pkg/es/tests/commands"
	"github.com/contextgg/pkg/es/tests/eventdata"
	"github.com/contextgg/pkg/es/tests/handlers"
	"github.com/contextgg/pkg/es/tests/sagas"
)

func SetupDomain() es.Config {
	cfg := es.NewConfig()

	cfg.
		Aggregate(aggregates.NewDemo).
		Commands(
			&commands.NewDemo{},
			&commands.AddDescription{},
		)

	cfg.
		Aggregate(aggregates.NewEntry).
		Handler(handlers.NewAddEntryHandler()).
		Commands(
			&commands.LedgerAddEntryCommand{},
		)
	cfg.
		Aggregate(aggregates.NewLineItem)

	cfg.
		Saga(sagas.NewDescription()).
		Events(
			&eventdata.DemoCreated{},
		)
	cfg.
		Saga(sagas.NewDescription2()).
		Events(
			&eventdata.DemoCreated{},
		)

	return cfg
}
