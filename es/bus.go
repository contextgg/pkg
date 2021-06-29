package es

type Bus interface {
	CommandBus
	EventBus

	AddSaga(saga Saga, events ...interface{})
	AddProjector(store Store, projector Projector, events ...interface{})
}

type bus struct {
	CommandBus
	EventBus
}

func (b *bus) AddSaga(saga Saga, events ...interface{}) {
	handler := NewSagaHandler(b, saga)
	b.AddHandler(handler, MatchAnyEventOf(events...))
}

func (b *bus) AddProjector(store Store, projector Projector, events ...interface{}) {
	handler := NewProjectorHandler(store, projector)
	b.AddHandler(handler, MatchAnyEventOf(events...))
}

func NewBus() Bus {
	commandBus := NewCommandBus()
	eventBus := NewEventBus()

	return &bus{
		commandBus,
		eventBus,
	}
}
