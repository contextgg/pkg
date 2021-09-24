package pubsub

import (
	"context"
	"fmt"
	"time"

	"cloud.google.com/go/pubsub"

	"github.com/contextgg/pkg/es"
	"github.com/contextgg/pkg/events"
	"github.com/contextgg/pkg/events/json"
	"github.com/contextgg/pkg/logger"
)

func getTopic(l logger.Logger, ctx context.Context, cli *pubsub.Client, topicName string) (*pubsub.Topic, error) {
	topic := cli.Topic(topicName)
	if ok, err := topic.Exists(ctx); err != nil {
		l.Error("topic.Exists", "err", err)
		return nil, err
	} else if !ok {
		if topic, err = cli.CreateTopic(ctx, topicName); err != nil {
			l.Error("cli.CreateTopic", "topicName", topicName)
			return nil, err
		}
	}
	return topic, nil
}

func getSubscription(l logger.Logger, ctx context.Context, cli *pubsub.Client, appId, topicName string) (*pubsub.Subscription, error) {
	topic, err := getTopic(l, ctx, cli, topicName)
	if err != nil {
		return nil, err
	}

	subscriptionId := appId + "__" + topicName
	sub := cli.Subscription(subscriptionId)
	if ok, err := sub.Exists(ctx); err != nil {
		return nil, err
	} else if !ok {
		if sub, err = cli.CreateSubscription(ctx, subscriptionId,
			pubsub.SubscriptionConfig{
				Topic:       topic,
				AckDeadline: 60 * time.Second,
			},
		); err != nil {
			return nil, err
		}
	}

	return sub, nil
}

// Publisher pubsub
type Publisher struct {
	l             logger.Logger
	client        *pubsub.Client
	subscriptions []*pubsub.Subscription
	topic         *pubsub.Topic
	eventBus      es.EventBus
	errCh         chan es.EventBusError
	codec         events.EventCodec
}

// NewPublisher creates a publisher and subscribes
func NewPublisher(l logger.Logger, eventBus es.EventBus, appId string, projectID string, topicName string, listeners ...string) (es.EventPublisher, error) {
	ctx := context.Background()

	cli, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		l.Error("pubsub.NewClient", "err", err, "projectID", projectID, "topicName", topicName)
		return nil, fmt.Errorf("pubsub.NewClient: %v", err)
	}

	topic, err := getTopic(l, ctx, cli, topicName)
	if err != nil {
		l.Error("getTopic", "err", err, "projectID", projectID, "topicName", topicName)
		return nil, fmt.Errorf("getTopic: %v", err)
	}

	p := &Publisher{
		l:        l,
		client:   cli,
		topic:    topic,
		eventBus: eventBus,
		errCh:    make(chan es.EventBusError, 1),
		codec:    &json.EventCodec{},
	}

	for _, l := range listeners {
		if len(l) == 0 {
			continue
		}
		if l == topicName {
			return nil, fmt.Errorf("Can not subscribe to the publishing topic: %s", topicName)
		}

		if err := p.subscription(ctx, appId, l); err != nil {
			return nil, err
		}
	}

	pid := fmt.Sprintf("gcp:%s:%s", appId, topicName)
	eventBus.AddPublisher(pid, p)
	return p, nil
}

// PublishEvent via pubsub
func (c *Publisher) PublishEvent(ctx context.Context, event events.Event) error {
	data, err := c.codec.MarshalEvent(ctx, &event)
	if err != nil {
		c.l.Error("json.Marshal", "err", err)
		return err
	}

	publishCtx := context.Background()
	res := c.topic.Publish(ctx, &pubsub.Message{
		Data: data,
	})
	if _, err := res.Get(publishCtx); err != nil {
		c.l.Error("Could not publish event", "err", err)
		return err
	}

	c.l.Debug("Event Published via GCP pub/sub", "topic_id", c.topic.ID(), "event_type", event.Type, "event_aggregate_id", event.AggregateID, "event_aggregate_type", event.AggregateType)
	return nil
}

func (c *Publisher) Errors() <-chan es.EventBusError {
	return c.errCh
}

func (c *Publisher) Start() {
	for _, sub := range c.subscriptions {
		go c.handle(sub)
	}
}

// Close underlying connection
func (c *Publisher) Close() {
	if c.client != nil {
		c.l.Debug("Closing the pubsub connection")
		c.client.Close()
	}
}

func (c *Publisher) subscription(ctx context.Context, appId, topicName string) error {
	sub, err := getSubscription(c.l, ctx, c.client, appId, topicName)
	if err != nil {
		return err
	}

	c.subscriptions = append(c.subscriptions, sub)
	return nil
}

func (c *Publisher) handle(sub *pubsub.Subscription) {
	for {
		ctx := context.Background()
		if err := sub.Receive(ctx, c.handler); err != context.Canceled {
			select {
			case c.errCh <- es.EventBusError{Ctx: ctx, Err: fmt.Errorf("Could not receive: %v", err)}:
			default:
			}
		}
		time.Sleep(time.Second)
	}
}

func (c *Publisher) handler(ctx context.Context, msg *pubsub.Message) {
	evt, ctx, err := c.codec.UnmarshalEvent(ctx, msg.Data)
	if err != nil {
		select {
		case c.errCh <- es.EventBusError{Err: fmt.Errorf("Could not unmarshal event: %s", err.Error()), Ctx: ctx}:
		default:
		}
		msg.Nack()
		return
	}

	ctx = es.SetIsPublisher(ctx)

	// Notify all observers about the event.
	if err := c.eventBus.HandleEvent(ctx, *evt); err != nil {
		select {
		case c.errCh <- es.EventBusError{Err: fmt.Errorf("Could not handle event: %s", err.Error()), Ctx: ctx, Event: evt}:
		default:
		}
		msg.Nack()
		return
	}

	msg.Ack()
}
