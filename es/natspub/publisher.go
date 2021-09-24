package natspub

import (
	"context"
	"fmt"
	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/nats-io/nats.go"

	"github.com/contextgg/pkg/es"
	"github.com/contextgg/pkg/events"
	"github.com/contextgg/pkg/events/json"
	"github.com/contextgg/pkg/logger"
)

type publisher struct {
	l        logger.Logger
	eventBus es.EventBus

	nc            *nats.Conn
	topicName     string
	subscriptions []*nats.Subscription

	codec events.EventCodec
	errCh chan es.EventBusError
}

// NewPublisher creates a publisher and subscribes
func NewPublisher(l logger.Logger, eventBus es.EventBus, natsUrl string, appId string, topicName string, listeners ...string) (es.EventPublisher, error) {
	b := backoff.NewExponentialBackOff()
	b.InitialInterval = 30 * time.Second
	b.MaxInterval = 5 * time.Minute
	b.Multiplier = 1.7
	b.MaxElapsedTime = 10 * time.Minute

	var nc *nats.Conn

	tryConnect := func() error {
		// try connect
		var err error
		nc, err = nats.Connect(natsUrl)
		if err != nil {
			l.Info("Could not connect to nats server: %s will retry after %d seconds", err, b.NextBackOff()/time.Second)
			return err
		}
		return nil
	}

	err := backoff.Retry(tryConnect, b)
	if err != nil {
		l.Error("BackOff stopped retrying with Error '%s'", err)
		return nil, err
	}

	p := &publisher{
		l:         l,
		eventBus:  eventBus,
		nc:        nc,
		topicName: topicName,
		codec:     &json.EventCodec{},
		errCh:     make(chan es.EventBusError, 1),
	}

	for _, l := range listeners {
		if len(l) == 0 {
			continue
		}
		if l == topicName {
			return nil, fmt.Errorf("Can not subscribe to the publishing topic: %s", topicName)
		}

		// create the subscription
		sub, err := nc.QueueSubscribeSync(l, appId)
		if err != nil {
			nc.Close()
			return nil, err
		}

		p.subscriptions = append(p.subscriptions, sub)
	}

	pid := fmt.Sprintf("nats:%s:%s", appId, topicName)
	eventBus.AddPublisher(pid, p)
	return p, nil
}

// PublishEvent via pubsub
func (c *publisher) PublishEvent(ctx context.Context, event events.Event) error {
	data, err := c.codec.MarshalEvent(ctx, &event)
	if err != nil {
		c.l.Error("json.Marshal", "err", err)
		return err
	}

	if err := c.nc.Publish(c.topicName, data); err != nil {
		c.l.Error("Could not publish event", "err", err)
		return err
	}

	c.l.Debug("Event Published via GCP pub/sub", "topic_name", c.topicName, "event_type", event.Type, "event_aggregate_id", event.AggregateID, "event_aggregate_type", event.AggregateType)
	return nil
}

func (c *publisher) Errors() <-chan es.EventBusError {
	return c.errCh
}

func (c *publisher) Start() {
	for _, sub := range c.subscriptions {
		go c.handle(sub)
	}
}

func (c *publisher) handle(sub *nats.Subscription) {
	for {
		msg, err := sub.NextMsg(10 * time.Second)
		if err != nil {
			c.errCh <- es.EventBusError{Err: fmt.Errorf("Could not receive: %v", err)}
		}
		c.handler(msg)
	}
}

func (c *publisher) handler(msg *nats.Msg) {
	ctx := context.Background()
	evt, ctx, err := c.codec.UnmarshalEvent(ctx, msg.Data)
	if err != nil {
		c.errCh <- es.EventBusError{Err: fmt.Errorf("Could not unmarshal event: %s", err.Error()), Ctx: ctx}
		return
	}

	ctx = es.SetIsPublisher(ctx)

	// Notify all observers about the event.
	if err := c.eventBus.HandleEvent(ctx, *evt); err != nil {
		c.errCh <- es.EventBusError{Err: fmt.Errorf("Could not handle event: %s", err.Error()), Ctx: ctx, Event: evt}
		return
	}
}

// Close underlying connection
func (c *publisher) Close() {
	for _, s := range c.subscriptions {
		s.Unsubscribe()
	}

	if c.nc != nil {
		c.nc.Close()
	}
}
