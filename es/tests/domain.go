package tests

import (
	"github.com/contextgg/pkg/es"

	"github.com/contextgg/pkg/es/tests/aggregates"
	"github.com/contextgg/pkg/es/tests/commands"
	"github.com/contextgg/pkg/es/tests/eventdata"
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
		Saga(sagas.NewCounter()).
		Events(
			&eventdata.DemoCreated{},
		)

	return cfg
}
