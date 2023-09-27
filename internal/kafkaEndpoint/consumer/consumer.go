package consumer

import (
	"encoding/json"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/hulla-hoop/testSobes/internal/modeldb"
)

type KafkaConsumer struct {
	c         *kafka.Consumer
	wg        *sync.WaitGroup
	inflogger *log.Logger
	errLogger *log.Logger
}

func New(c *kafka.Consumer, wg *sync.WaitGroup, inflogger *log.Logger, errLogger *log.Logger) *KafkaConsumer {
	return &KafkaConsumer{
		c:         c,
		wg:        wg,
		inflogger: inflogger,
		errLogger: errLogger,
	}
}

func (c *KafkaConsumer) Consumer(f chan modeldb.User) {

	topic := "FIO"
	err := c.c.SubscribeTopics([]string{topic}, nil)
	if err != nil {
		c.errLogger.Println(err)
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
			c.inflogger.Println("Выход из горутины Consumer прекращено сигналом - ", sig)
			c.wg.Done()
			run = false
			close(f)
		default:
			ev, err := c.c.ReadMessage(100 * time.Millisecond)
			if err != nil {
				continue
			}
			var U modeldb.User
			err = json.Unmarshal(ev.Key, &U)
			c.inflogger.Println("Получено сообщение из очереди FIO  ---- ", U)
			if err != nil {
				c.errLogger.Println(err)
			}
			f <- U
		}
	}

}
