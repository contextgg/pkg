package es

type Bus interface {
	CommandBus
	EventBus

	AddSaga(saga Saga, events ...interface{})
}

type bus struct {
	CommandBus
	EventBus
}

func (b *bus) AddSaga(saga Saga, events ...interface{}) {
	handler := NewSagaHandler(b, saga)
	b.AddHandler(handler, MatchAnyEventOf(events...))
}

// bus.AddSaga(registerSaga, &eventdata.BankAccountRegistered{})

func NewBus() Bus {
	commandBus := NewCommandBus()
	eventBus := NewEventBus()

	return &bus{
		commandBus,
		eventBus,
	}
}
