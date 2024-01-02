package redisc

import (
	"context"
	"github.com/go-redis/redis/v8"
)

const (
	SubscribeKey = "0"
	PublishKey = "1"
)

func Publish(ctx context.Context, channel string, message string) error {
	return Redisc.Publish(ctx, channel, message).Err()
}

func PSubscribe(ctx context.Context, channels ...string) *redis.PubSub {
	pubsub := Redisc.Subscribe(ctx)
	if err := pubsub.PSubscribe(ctx, channels...); err != nil {
		panic(err)
	}
	return pubsub
}

func Subscribe(ctx context.Context, channels ...string) *redis.PubSub {
	return Redisc.Subscribe(ctx, channels...)
}

