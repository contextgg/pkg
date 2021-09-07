package postgres

import (
	"context"
	"database/sql"
	"testing"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/uptrace/bun/extra/bundebug"

	"github.com/contextgg/pkg/es"
	"github.com/contextgg/pkg/events"
	"github.com/contextgg/pkg/ns"
)

func createDb() *bun.DB {
	conn := pgdriver.NewConnector(
		pgdriver.WithDSN("postgresql://localhost:5432/dev?sslmode=disable"),
		pgdriver.WithDatabase("testDb"),
		pgdriver.WithUser("contextgg"),
		pgdriver.WithPassword("mysecret"),
	)

	sqldb := sql.OpenDB(conn)
	db := bun.NewDB(sqldb, pgdialect.New())
	db.AddQueryHook(bundebug.NewQueryHook(bundebug.WithVerbose()))
	return db
}

type Fake struct {
	es.BaseAggregateHolder

	Name string `json:"name"`
}

func NewFake(id string) es.Entity {
	return &Fake{
		BaseAggregateHolder: es.NewBaseAggregateHolder(id, "Fake"),
	}
}

type FakeCommand struct {
	es.BaseCommand
}
type FakeEvent struct {
	Name string
}

type FakeSourced struct {
	es.BaseAggregateSourced

	Name string `json:"name"`
}

// HandleCommand create events and validate based on such command
func (a *FakeSourced) HandleCommand(ctx context.Context, cmd es.Command) error {
	a.StoreEventData(ctx, &FakeEvent{
		Name: "hello2",
	})
	return nil
}

// ApplyEvent to auth
func (a *FakeSourced) ApplyEvent(ctx context.Context, event events.Event) error {
	e := event.Data.(*FakeEvent)
	a.Name = e.Name
	return nil
}

func NewFakeSourced(id string) es.Entity {
	return &FakeSourced{
		BaseAggregateSourced: es.NewBaseAggregateSourced(id, "FakeSourced"),
	}
}

func TestEntity(t *testing.T) {
	db := createDb()
	data := NewPostgresData(
		db,
		es.InitializeEntities(
			&Fake{},
		),
	)

	id := "72c096f0-d64a-11eb-b8bc-0242ac130003"

	testEntitySave := func(t *testing.T) {
		ctx := ns.SetNamespace(context.TODO(), "temp")

		entity := NewFake(id).(*Fake)
		entity.Name = "hello"

		namespace := ns.FromContext(ctx)
		if err := data.SaveEntity(ctx, namespace, entity); err != nil {
			t.Error(err)
			return
		}
	}
	testEntityLoad := func(t *testing.T) {
		ctx := ns.SetNamespace(context.TODO(), "temp")
		namespace := ns.FromContext(ctx)

		entity := NewFake(id)
		if err := data.LoadEntity(ctx, namespace, entity); err != nil {
			t.Error(err)
			return
		}

		t.Log(entity)
	}

	t.Run("TestEntitySave", testEntitySave)
	t.Run("TestEntityLoad", testEntityLoad)
}

func TestSourcedSave(t *testing.T) {
	db := createDb()
	data := NewPostgresData(
		db,
		es.InitializeEvents(),
		es.InitializeSnapshots(),
		es.InitializeEntities(
			&FakeSourced{},
		),
	)

	eventBus := es.NewEventBus()
	fakeRepo := es.NewStore(data, eventBus, NewFakeSourced)
	fakeHandler := es.NewAggregateSourcedHandler(fakeRepo)

	ctx := ns.SetNamespace(context.TODO(), "temp")

	if err := fakeHandler.HandleCommand(ctx, &FakeCommand{
		BaseCommand: es.BaseCommand{
			AggregateId: "72c096f0-d64a-11eb-b8bc-0242ac130003",
		},
	}); err != nil {
		t.Error(err)
		return
	}

	if err := fakeHandler.HandleCommand(ctx, &FakeCommand{
		BaseCommand: es.BaseCommand{
			AggregateId: "72c096f0-d64a-11eb-b8bc-0242ac130003",
		},
	}); err != nil {
		t.Error(err)
		return
	}
}
