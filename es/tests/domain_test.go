package tests

import (
	"context"
	"os"
	"testing"

	"github.com/contextgg/pkg/es"
	"github.com/contextgg/pkg/es/db"
	"github.com/contextgg/pkg/es/tests/aggregates"
	"github.com/contextgg/pkg/es/tests/commands"
	"github.com/contextgg/pkg/logger"
	"github.com/contextgg/pkg/ns"
	"go.uber.org/zap"
)

func SetupBus() (es.CommandHandler, error) {
	conn, err := es.NewConn(
		db.WithDbHost(os.Getenv("DB_HOSTNAME")),
		db.WithDbName("testdb"),
		db.WithDbUser("contextgg"),
		db.WithDbPassword("mysecret"),
		db.WithDebug(true),
		db.Recreate(true),
	)
	if err != nil {
		return nil, err
	}

	cfg := SetupDomain()

	cli, err := es.NewClient(conn, cfg)

	z, _ := zap.NewDevelopment()
	l := logger.NewLogger(z)

	// migrate the DB!
	if err := es.MigrateDatabase(
		db,
		es.InitializeEvents(),
		es.InitializeSnapshots(),
		es.InitializeEntities(
			&aggregates.Demo{},
		),
	); err != nil {
		return nil, err
	}

	return NewBus(db, l)
}

func TestIt(t *testing.T) {
	bus, err := SetupBus()
	if err != nil {
		t.Error(err)
		return
	}

	cmds := []es.Command{
		&commands.NewDemo{
			BaseCommand: es.BaseCommand{
				AggregateId: "d63b875a-a664-410c-9102-21bfd7381f6e",
			},
			Name: "Hello2",
		},
	}

	for _, cmd := range cmds {
		ctx := context.Background()
		ctx = ns.SetNamespace(ctx, "test")

		if err := bus.HandleCommand(ctx, cmd); err != nil {
			t.Error(err)
			return
		}
	}
}
