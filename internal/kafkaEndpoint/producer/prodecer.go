package producer

import (
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/hulla-hoop/testSobes/internal/service"
)

type KafkaProducer struct {
	p  *kafka.Producer
	wg *sync.WaitGroup
}

func New(p *kafka.Producer, wg *sync.WaitGroup) *KafkaProducer {
	return &KafkaProducer{
		p:  p,
		wg: wg,
	}
}

func (p *KafkaProducer) Producer(user chan service.UserFailed) {

	topic := "FIO_FAILED"

	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	run := true
	for run {
		select {
		case sig := <-sigchan:
			fmt.Println("Ybito", sig)
			p.wg.Done()
			run = false
		default:
			uFailed, err := json.Marshal(<-user)
			if err != nil {
				fmt.Println(err)
			}
			key := uFailed
			p.p.Produce(&kafka.Message{
				TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
				Key:            []byte(key),
			}, nil)
		}
	}
}
