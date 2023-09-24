package main

import (
	"fmt"
	"os"
	"sync"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/hulla-hoop/testSobes/internal/config"
	"github.com/hulla-hoop/testSobes/internal/kafkaEndpoint/consumer"
	"github.com/hulla-hoop/testSobes/internal/kafkaEndpoint/producer"
	"github.com/hulla-hoop/testSobes/internal/psql"
	"github.com/hulla-hoop/testSobes/internal/service"
)

var wg sync.WaitGroup

func main() {

	psql.InitDb()
	s := service.New(&wg)

	config := config.New()
	UserChan := make(chan service.User)
	UserChanFailed := make(chan service.UserFailed)
	wg.Add(3)

	conf := make(kafka.ConfigMap)

	conf["bootstrap.servers"] = config.BootstrapService
	conf["group.id"] = config.GroupID
	conf["auto.offset.reset"] = config.AutoOffsetReset

	conf2 := make(kafka.ConfigMap)

	conf2["bootstrap.servers"] = config.BootstrapService

	p, err := kafka.NewProducer(&conf2)

	if err != nil {
		fmt.Printf("Failed to create producer: %s", err)
		os.Exit(1)
	}
	producer := producer.New(p, &wg)

	c, err := kafka.NewConsumer(&conf)

	if err != nil {
		fmt.Printf("Failed to create consumer: %s", err)
		os.Exit(1)
	}

	consumer := consumer.New(c, &wg)

	go consumer.Consumer(UserChan)
	go s.Distribution(UserChan, UserChanFailed)
	go producer.Producer(UserChanFailed)

	wg.Wait()
	c.Close()

}
