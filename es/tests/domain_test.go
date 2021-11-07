package tests

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"testing"

	"github.com/contextgg/pkg/db/pg"
	"github.com/contextgg/pkg/es"
	"github.com/contextgg/pkg/es/tests/commands"
	"github.com/contextgg/pkg/logger"
	"github.com/contextgg/pkg/ns"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/uptrace/bun/extra/bundebug"
	"go.uber.org/zap"
)

func SetupBus() (es.Bus, error) {
	z, _ := zap.NewDevelopment()
	l := logger.NewLogger(z)

	hostname := os.Getenv("DB_HOSTNAME")
	if len(hostname) == 0 {
		hostname = "localhost"
	}

	dbConn := fmt.Sprintf("postgresql://%s:5432/testdb?sslmode=disable", hostname)
	dbName := "testdb"
	dbUser := "contextgg"
	dbPass := "mysecret"

	err := pg.Recreate(func() (*bun.DB, error) {
		conn := pgdriver.NewConnector(
			pgdriver.WithDSN(dbConn),
			pgdriver.WithDatabase("postgres"),
			pgdriver.WithUser(dbUser),
			pgdriver.WithPassword(dbPass),
		)
		sqldb := sql.OpenDB(conn)
		return bun.NewDB(sqldb, pgdialect.New()), nil
	}, dbName)
	if err != nil {
		return nil, err
	}

	conn := pgdriver.NewConnector(
		pgdriver.WithDSN(dbConn),
		pgdriver.WithDatabase(dbName),
		pgdriver.WithUser(dbUser),
		pgdriver.WithPassword(dbPass),
	)

	sqldb := sql.OpenDB(conn)
	db := bun.NewDB(sqldb, pgdialect.New())
	db.AddQueryHook(bundebug.NewQueryHook(bundebug.WithVerbose(true)))

	return NewBus(db, l)

	// fixture := dbfixture.New(db, dbfixture.WithTruncateTables())
	// ferr := fixture.Load(context.Background(), os.DirFS("../data"), "profile.yml")
	// Expect(ferr).ShouldNot(HaveOccurred())
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
			Name: "Hello",
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
