package es

import "context"

const publisherKey = 0

func SetIsPublisher(ctx context.Context) context.Context {
	return context.WithValue(ctx, publisherKey, true)
}

func IsPublisherFromContext(ctx context.Context) bool {
	isPublisher, ok := ctx.Value(publisherKey).(bool)
	if ok {
		return isPublisher
	}
	return false
}
