package tests

import (
	"github.com/contextgg/pkg/es"
	"github.com/contextgg/pkg/es/postgres"
	"github.com/contextgg/pkg/logger"
	"github.com/uptrace/bun"

	"github.com/contextgg/pkg/es/tests/aggregates"
	"github.com/contextgg/pkg/es/tests/commands"
	"github.com/contextgg/pkg/es/tests/eventdata"
	"github.com/contextgg/pkg/es/tests/sagas"
)

func NewBus(db *bun.DB, log logger.Logger) (es.Bus, error) {
	data := postgres.NewPostgresData(
		db,
		es.InitializeEvents(),
		es.InitializeSnapshots(),
		es.InitializeEntities(
			&aggregates.Demo{},
		),
	)

	bus := es.NewBus()

	// uniter := es.NewUniter()
	demoStore := es.NewStore(data, bus, aggregates.NewDemo)
	demoHandler := es.NewAggregateSourcedHandler(demoStore)

	bus.SetHandler(demoHandler, &commands.NewDemo{})

	bus.AddSaga(
		sagas.NewCounter(demoStore),
		&eventdata.DemoCreated{},
	)

	return bus, nil
}
