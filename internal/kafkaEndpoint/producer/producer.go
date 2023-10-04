package producer

import (
	"encoding/json"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/hulla-hoop/testSobes/internal/config"
	"github.com/hulla-hoop/testSobes/internal/service"
)

type KafkaProducer struct {
	p         *kafka.Producer
	wg        *sync.WaitGroup
	inflogger *log.Logger
	errLogger *log.Logger
	cfgKafk   *config.Configkafka
}

func New(p *kafka.Producer, wg *sync.WaitGroup, inflogger *log.Logger, errLogger *log.Logger, cfgKafk *config.Configkafka) *KafkaProducer {
	return &KafkaProducer{
		p:         p,
		wg:        wg,
		inflogger: inflogger,
		errLogger: errLogger,
		cfgKafk:   cfgKafk,
	}
}

func (p *KafkaProducer) Producer(user chan service.UserFailed) {

	topic := p.cfgKafk.TopicErr

	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	run := true
	for run {
		select {
		case sig := <-sigchan:
			p.inflogger.Println("Выход из горутины Producer прекращено сигналом - ", sig)
			p.wg.Done()
			run = false
		default:
			uFailed, err := json.Marshal(<-user)
			if err != nil {
				p.errLogger.Println(err)
			}
			value := uFailed
			p.p.Produce(&kafka.Message{
				TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
				Value:          []byte(value),
			}, nil)
		}
	}
}
