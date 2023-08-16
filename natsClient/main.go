package main

import (
	"context"
	"fmt"
	"github.com/nats-io/nats.go"
	"log"
	"os"
	"time"
)

func main() {
	payload, err := os.ReadFile("payload.json")

	nc, err := nats.Connect("nats://localhost:4222")
	if err != nil {
		log.Fatal("can't connect to nats server")
	}

	/*for i := 0; i < 10; i++ {
		simplePublish(err, nc, payload)
	}*/

	ctx := context.Background()
	ctxTimeout, cancel := context.WithTimeout(ctx, 60*time.Second)

	dataChannel := make(chan string, 16)

	for i := 0; i < 100; i++ {
		go func() {
			dataChannel <- string(payload)
			fmt.Println("send to channel")
		}()
	}

	printIntegers(ctxTimeout, nc, dataChannel)

	cancel()
	fmt.Println("IsDone")
}

func simplePublish(err error, nc *nats.Conn, payload []byte) {
	const subject = "json_payload"
	err = nc.Publish(subject, payload)

	if err != nil {
		log.Fatal("can't publish to nats server")
	}
}

func printIntegers(ctx context.Context, conn *nats.Conn, stringChannel <-chan string) {
	for {
		select {
		case str := <-stringChannel:
			{
				fmt.Printf("read from channel %s\n", str)
				const subject = "json_payload"
				_ = conn.Publish(subject, []byte(str))
			}
		case <-ctx.Done():
			return
		}
	}
}
