package aggregates

import (
	"context"

	"github.com/contextgg/pkg/es2"

	"github.com/contextgg/pkg/es2/example/data/commands"
	"github.com/contextgg/pkg/es2/example/data/eventdata"
)

type User struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type UserHandler struct {
	es2.AggregateSourced
}

func (h *UserHandler) handleNewUser(ctx context.Context, cmd *commands.NewUser) error {
	h.StoreEventData(ctx, &eventdata.UserCreated{
		FirstName: cmd.FirstName,
		LastName:  cmd.LastName,
		Username:  cmd.Username,
	})
	return nil
}

func (h *UserHandler) Handle(ctx context.Context, cmd es2.Command) error {
	switch c := cmd.(type) {
	case *commands.NewUser:
		return h.handleNewUser(ctx, c)
	default:
		return es2.ErrNotHandled
	}
}

func NewUserHandler() es2.CommandHandler {
	return &UserHandler{
		AggregateSourced: es2.NewBaseAggregateSourced("User"),
	}
}
