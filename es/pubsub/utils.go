package pubsub

import (
	"context"
	"time"

	"cloud.google.com/go/pubsub"
	"github.com/contextgg/pkg/logger"
)

func getTopic(ctx context.Context, l logger.Logger, cli *pubsub.Client, topicName string) (*pubsub.Topic, error) {
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

func getSubscription(ctx context.Context, l logger.Logger, cli *pubsub.Client, appId, topicName string) (*pubsub.Subscription, error) {
	topic, err := getTopic(ctx, l, cli, topicName)
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
