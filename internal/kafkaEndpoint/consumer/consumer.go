package consumer

import (
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/hulla-hoop/testSobes/internal/service"
)

type KafkaConsumer struct {
	c  *kafka.Consumer
	wg *sync.WaitGroup
}

func New(c *kafka.Consumer, wg *sync.WaitGroup) *KafkaConsumer {
	return &KafkaConsumer{
		c:  c,
		wg: wg,
	}
}

func (c *KafkaConsumer) Consumer(f chan service.User) {

	topic := "FIO"
	err := c.c.SubscribeTopics([]string{topic}, nil)
	if err != nil {
		os.Exit(1)
	}
	// Set up a channel for handling Ctrl-C, etc
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	// Process messages

	run := true
	for run {
		select {
		case sig := <-sigchan:
			fmt.Printf("Caught signal %v: terminating\n", sig)
			c.wg.Done()
			run = false
			close(f)
		default:
			ev, err := c.c.ReadMessage(100 * time.Millisecond)
			if err != nil {
				// Errors are informational and automatically handled by the consumer
				continue
			}
			var U service.User
			err = json.Unmarshal(ev.Key, &U)
			if err != nil {
				fmt.Println(err)
			}

			fmt.Println(U)
			f <- U
		}
	}

}
