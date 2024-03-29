package tests

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/contextgg/pkg/es"
	"github.com/contextgg/pkg/es/db"
	"github.com/contextgg/pkg/es/tests/aggregates"
	"github.com/contextgg/pkg/es/tests/commands"
	"github.com/contextgg/pkg/ns"
)

func SetupClient() (es.Client, error) {
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

	// nc, err := nats.Connect("nats://localhost:4222")
	// if err != nil {
	// 	return nil, err
	// }

	cfg := SetupDomain()

	cli, err := es.NewClient(cfg, conn)
	if err != nil {
		return nil, err
	}

	return cli, nil
}

func TestIt(t *testing.T) {
	cli, err := SetupClient()
	if err != nil {
		t.Error(err)
		return
	}

	cmds := []es.Command{
		&commands.NewDemo{
			BaseCommand: es.BaseCommand{
				AggregateId: "d63b875a-a664-410c-9102-21bfd7381f6e",
			},
			Name: "Demo",
		},
		&commands.LedgerAddEntryCommand{
			BaseCommand: es.BaseCommand{
				AggregateId: "b3fc2c26-4ef1-4edb-9363-82c52b30ba87",
			},
			LedgerId:    "685d3238-2cb9-4171-a739-cb887b2253c4",
			Book:        "USD",
			Description: "Can have money",
			At:          time.Now(),
			Items: []*commands.LedgerLineItem{
				{
					AccountId:    "d63b875a-a664-410c-9102-21bfd7381f6e",
					SubAccountId: "deposit",
					Credit:       100,
					Debit:        0,
				},
			},
		},
	}

	ctx := context.Background()
	ctx = ns.SetNamespace(ctx, "test")

	// create a unit.
	unit, err := cli.Unit(ctx)
	if err != nil {
		t.Error(err)
		return
	}
	ctx = es.SetUnit(ctx, unit)

	tx, err := unit.Begin(ctx)
	if err != nil {
		t.Error(err)
		return
	}
	defer func() {
		tx.Rollback(ctx)
	}()

	if err := unit.Dispatch(ctx, cmds...); err != nil {
		t.Error(err)
		return
	}

	if err := tx.Commit(ctx); err != nil {
		t.Error(err)
		return
	}

	q := es.NewQuery[*aggregates.Demo]()
	demo, err := q.Load(ctx, "d63b875a-a664-410c-9102-21bfd7381f6e")
	if err != nil {
		t.Error(err)
		return
	}

	t.Log(demo)
}
