package tests

import (
	"github.com/contextgg/pkg/es"
	"github.com/contextgg/pkg/logger"
	"github.com/uptrace/bun"

	"github.com/contextgg/pkg/es/tests/aggregates"
	"github.com/contextgg/pkg/es/tests/commands"
	"github.com/contextgg/pkg/es/tests/eventdata"
	"github.com/contextgg/pkg/es/tests/sagas"
)

func NewBus(db *bun.DB, log logger.Logger) (es.CommandHandler, error) {
	uniter := es.NewUniter(db)
	eventBus := es.NewEventBus(uniter)

	entityRegistry := es.NewEntityRegistry()
	entityRegistry.SetEntity(
		&aggregates.Demo{},
		es.EntityFactory(aggregates.NewDemo),
	)
	entityStore := es.NewEntityStore(entityRegistry, eventBus)

	commandBus := es.NewCommandBus(uniter)
	commandBus.SetHandler(
		es.NewAggregateSourcedHandler(entityStore, &aggregates.Demo{}),
		&commands.NewDemo{},
		&commands.AddDescription{},
	)

	eventBus.AddHandler(
		es.NewSagaHandler(commandBus, sagas.NewCounter(entityStore)),
		es.MatchAnyEventDataOf(&eventdata.DemoCreated{}),
	)

	// setup subscribers here!.
	external := es.UseEventHandlerMiddleware(eventBus, es.EventUniterMiddleware(uniter))
	log.Info("External setup", "external", external)

	return commandBus, nil
}
