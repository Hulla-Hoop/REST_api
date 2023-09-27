package producer

import (
	"encoding/json"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/hulla-hoop/testSobes/internal/service"
)

type KafkaProducer struct {
	p         *kafka.Producer
	wg        *sync.WaitGroup
	inflogger *log.Logger
	errLogger *log.Logger
}

func New(p *kafka.Producer, wg *sync.WaitGroup, inflogger *log.Logger, errLogger *log.Logger) *KafkaProducer {
	return &KafkaProducer{
		p:         p,
		wg:        wg,
		inflogger: inflogger,
		errLogger: errLogger,
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
			p.inflogger.Println("Выход из горутины Consumer прекращено сигналом - ", sig)
			p.wg.Done()
			run = false
		default:
			uFailed, err := json.Marshal(<-user)
			if err != nil {
				p.errLogger.Println(err)
			}
			key := uFailed
			p.p.Produce(&kafka.Message{
				TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
				Key:            []byte(key),
			}, nil)
		}
	}
}
