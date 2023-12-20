package redisc

import (
	"fmt"
	"ginchat/config"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestInit(t *testing.T) {
	c := config.Init("../etc")
	rdb := Init(c)

	if err := rdb.Set(ctx, "redisc-test1", "ok", 0).Err(); err != nil {
		panic(err)
	}

	val, err := rdb.Get(ctx, "redisc-test1").Result()
	if err != nil {
		panic(err)
	}
	assert.Equal(t, "ok", val)
}

func TestPSubscribe(t *testing.T) {
	c := config.Init("../etc")
	Init(c)

	pubsub := PSubscribe(ctx, "redis-chan-*")
	defer pubsub.Close()

	channel := pubsub.Channel()
	for msg := range channel {
		fmt.Printf("Channel: %v, Message: %v\n", msg.Channel, msg.Payload)
	}

	//for {
	//	message, err := pubsub.ReceiveMessage(ctx)
	//	if err != nil {
	//		panic(err)
	//	}
	//
	//	fmt.Printf("Channel: %v, Message: %v", message.Channel, message.Payload)
	//}
}

func TestSubscribe(t *testing.T) {
	c := config.Init("../etc")
	Init(c)

	pubsub := Subscribe(ctx, "redis1")
	defer pubsub.Close()
	//fmt.Println("test")

	channel := pubsub.Channel()
	for msg := range channel {
		fmt.Printf("Channel: %s, Message: %s\n", msg.Channel, msg.Payload)
	}

	//for {
	//	message, err := pubsub.ReceiveMessage(ctx)
	//	if err != nil {
	//		panic(err)
	//	}
	//
	//	fmt.Printf("Channel: %v, Message: %v", message.Channel, message.Payload)
	//}
}
