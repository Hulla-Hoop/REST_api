package main

import (
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/hulla-hoop/testSobes/internal/config"
	"github.com/hulla-hoop/testSobes/internal/kafkaEndpoint/consumer"
	"github.com/hulla-hoop/testSobes/internal/kafkaEndpoint/producer"
	"github.com/hulla-hoop/testSobes/internal/modeldb"
	"github.com/hulla-hoop/testSobes/internal/psql"
	"github.com/hulla-hoop/testSobes/internal/service"
	"github.com/hulla-hoop/testSobes/pkg/app"
	"github.com/joho/godotenv"
)

var wg sync.WaitGroup

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	infLogger := log.New(os.Stdout, "INFO:  ", log.Ldate|log.Lshortfile)
	errLogger := log.New(os.Stdout, "ERROR:  ", log.Ldate|log.Lshortfile)
	psql := psql.InitDb()

	a := app.New(psql, infLogger, errLogger)

	go a.Start()

	cfgApi := config.NewCfgApi()

	s := service.New(errLogger, cfgApi)

	config := config.New()
	UserChan := make(chan modeldb.User)
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
	producer := producer.New(p, &wg, infLogger, errLogger, config)

	c, err := kafka.NewConsumer(&conf)

	if err != nil {
		fmt.Printf("Failed to create consumer: %s", err)
		os.Exit(1)
	}

	consumer := consumer.New(c, &wg, infLogger, errLogger, config)

	go consumer.Consumer(UserChan)
	go service.Distribution(s, UserChan, UserChanFailed, infLogger, errLogger, &wg, psql)
	go producer.Producer(UserChanFailed)

	wg.Wait()
	c.Close()

}
