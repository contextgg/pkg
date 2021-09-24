package natspub

import (
	"context"
	"sync"
	"testing"

	"github.com/contextgg/pkg/es"
	"github.com/contextgg/pkg/events"
	"github.com/contextgg/pkg/events/json"
	"github.com/contextgg/pkg/logger"
	"github.com/contextgg/pkg/types"
	"github.com/nats-io/nats.go"
	"go.uber.org/zap"
)

var (
	_ es.EventHandler = &Handler{}
	_ es.Entity       = &Entity{}
)

type Handler struct {
	l  logger.Logger
	wg sync.WaitGroup
}

func (h *Handler) HandleEvent(ctx context.Context, evt events.Event) error {
	h.l.Debug("Has event", "event", evt)

	h.wg.Done()
	return nil
}

type Entity struct {
	es.BaseAggregateHolder
}

type EventData struct {
	Message string `json:"message"`
}

func toEvent(ctx context.Context, version int) *events.Event {
	agg := &Entity{
		BaseAggregateHolder: es.NewBaseAggregateHolder("82430ae6-96a4-4e6d-911f-92cf1812b7b6", "Entity"),
	}
	evt := es.NewEvent(ctx, agg, version, &EventData{
		Message: "Hello world",
	})
	return &evt
}

func toData(ctx context.Context, version int) ([]byte, error) {
	evt := toEvent(ctx, version)

	codec := &json.EventCodec{}
	return codec.MarshalEvent(ctx, evt)
}

func TestIt(t *testing.T) {
	const (
		url   = "nats://localhost:4222"
		topic = "demo-natspub"
		sub   = "demo-listener"
	)

	types.SetTypeData(&EventData{}, false)

	zapL, err := zap.NewDevelopment()
	if err != nil {
		t.Error(err)
		return
	}

	l := logger.NewLogger(zapL)
	handler := &Handler{
		l: l,
	}
	handler.wg.Add(2)

	ctx := context.TODO()
	d1, err := toData(ctx, 1)
	if err != nil {
		t.Error(err)
		return
	}
	d2, err := toData(ctx, 2)
	if err != nil {
		t.Error(err)
		return
	}

	bus := es.NewEventBus()
	bus.AddHandler(handler, es.MatchAny())

	pub, err := NewPublisher(l, bus, url, topic, topic, sub)
	if err != nil {
		t.Error(err)
		return
	}
	pub.Start()

	// write an via nats!
	nc, err := nats.Connect(url)
	if err != nil {
		t.Error(err)
		return
	}

	if err := nc.Publish(sub, d1); err != nil {
		t.Error(err)
		return
	}
	if err := nc.Publish(sub, d2); err != nil {
		t.Error(err)
		return
	}

	evt := toEvent(ctx, 1)
	if err := pub.PublishEvent(ctx, *evt); err != nil {
		t.Error(err)
		return
	}

	handler.wg.Wait()
}
