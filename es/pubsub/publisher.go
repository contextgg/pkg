package pubsub

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"cloud.google.com/go/pubsub"
	"github.com/rs/zerolog/log"

	"github.com/contextgg/pkg/es"
	"github.com/contextgg/pkg/events"
	"github.com/contextgg/pkg/ns"
)

func getTopic(ctx context.Context, cli *pubsub.Client, topicName string) (*pubsub.Topic, error) {
	topic := cli.Topic(topicName)
	if ok, err := topic.Exists(ctx); err != nil {
		log.
			Error().
			Err(err).
			Msg("topic.Exists")
		return nil, err
	} else if !ok {
		if topic, err = cli.CreateTopic(ctx, topicName); err != nil {
			log.
				Error().
				Err(err).
				Str("topicName", topicName).
				Msg("cli.CreateTopic")
			return nil, err
		}
	}
	return topic, nil
}

func getSubscription(ctx context.Context, cli *pubsub.Client, appId, topicName string) (*pubsub.Subscription, error) {
	topic, err := getTopic(ctx, cli, topicName)
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
	client        *pubsub.Client
	subscriptions []*pubsub.Subscription
	topic         *pubsub.Topic
	eventBus      es.EventBus
	errCh         chan es.EventBusError
}

// NewPublisher creates a publisher and subscribes
func NewPublisher(eventBus es.EventBus, appId string, projectID string, topicName string, listeners ...string) (es.EventPublisher, error) {
	ctx := context.Background()

	cli, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		log.
			Error().
			Err(err).
			Str("projectID", projectID).
			Str("topicName", topicName).
			Msg("pubsub.NewClient")
		return nil, fmt.Errorf("pubsub.NewClient: %v", err)
	}

	topic, err := getTopic(ctx, cli, topicName)
	if err != nil {
		log.
			Error().
			Err(err).
			Str("projectID", projectID).
			Str("topicName", topicName).
			Msg("getTopic")
		return nil, fmt.Errorf("getTopic: %v", err)
	}

	p := &Publisher{
		client:   cli,
		topic:    topic,
		eventBus: eventBus,
		errCh:    make(chan es.EventBusError, 1),
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
	msg, err := json.Marshal(event)
	if err != nil {
		log.
			Error().
			Err(err).
			Msg("json.Marshal")
		return err
	}

	namespace := ns.FromContext(ctx)
	attrs := map[string]string{
		"namespace": namespace,
	}

	publishCtx := context.Background()
	res := c.topic.Publish(publishCtx, &pubsub.Message{
		Data:       msg,
		Attributes: attrs,
	})
	if _, err := res.Get(publishCtx); err != nil {
		log.
			Error().
			Err(err).
			Msg("Could not publish event")
		return err
	}

	log.
		Debug().
		Str("topic_id", c.topic.ID()).
		Str("event_type", event.Type).
		Str("event_aggregate_id", event.AggregateID).
		Str("event_aggregate_type", event.AggregateType).
		Msg("Event Published via GCP pub/sub")
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
		log.
			Debug().
			Msg("Closing the pubsub connection")
		c.client.Close()
	}
}

func (c *Publisher) subscription(ctx context.Context, appId, topicName string) error {
	sub, err := getSubscription(ctx, c.client, appId, topicName)
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
	r := bytes.NewReader(msg.Data)
	evt, err := events.EventDecoder(r)
	if err != nil {
		select {
		case c.errCh <- es.EventBusError{Err: fmt.Errorf("Could not unmarshal event: %s", err.Error()), Ctx: ctx}:
		default:
		}
		msg.Nack()
		return
	}

	handlerCtx := ctx
	if msg.Attributes != nil {
		namespace, ok := msg.Attributes["namespace"]
		if ok && len(namespace) > 0 {
			handlerCtx = ns.SetNamespace(ctx, namespace)
		}
	}

	evt.Metadata["publisher"] = true

	// Notify all observers about the event.
	if err := c.eventBus.HandleEvent(handlerCtx, *evt); err != nil {
		select {
		case c.errCh <- es.EventBusError{Err: fmt.Errorf("Could not handle event: %s", err.Error()), Ctx: ctx, Event: evt}:
		default:
		}
		msg.Nack()
		return
	}

	msg.Ack()
}
