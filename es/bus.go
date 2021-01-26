package es

type Bus interface {
	CommandBus
	EventBus
}

type bus struct {
	CommandBus
	EventBus
}

func NewBus() Bus {
	commandBus := NewCommandBus()
	eventBus := NewEventBus()

	return &bus{
		commandBus,
		eventBus,
	}
}
