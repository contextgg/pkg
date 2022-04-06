package tests

import (
	"context"
	"testing"

	"github.com/contextgg/pkg/es2"
	"github.com/contextgg/pkg/es2/example/data/aggregates"
	"github.com/contextgg/pkg/es2/example/data/commands"
)

func TestXxx(t *testing.T) {
	commandRegistry := es2.NewCommandRegistry()
	commandRegistry.SetHandler(aggregates.NewUserHandler(), &commands.NewUser{})

	app, err := es2.Build(opts)
	if err != nil {
		t.Error(err)
		return
	}
	defer app.Close()

	ctx := context.TODO()
	if err := app.Dispatch(ctx, &commands.NewUser{
		Username:  "test",
		FirstName: "test",
		LastName:  "test",
	}); err != nil {
		t.Error(err)
		return
	}
}
