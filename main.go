package main

import (
	"context"
	"fmt"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

func main() {
	// connect to nats server
	nc, _ := nats.Connect(nats.DefaultURL)

	// create jetstream context from nats connection
	js, _ := jetstream.New(nc)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// get existing stream handle
	stream, _ := js.Stream(ctx, "foo")

	// retrieve consumer handle from a stream
	cons, _ := stream.Consumer(ctx, "cons")

	// consume messages from the consumer in callback
	cc, _ := cons.Consume(func(msg jetstream.Msg) {
		fmt.Println("Received jetstream message: ", string(msg.Data()))
		msg.Ack()
	})
	defer cc.Stop()
}
