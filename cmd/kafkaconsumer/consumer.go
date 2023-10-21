package main

import (
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/hulla-hoop/testSobes/internal/DB/psql"
	"github.com/hulla-hoop/testSobes/internal/DB/rediscash"
	"github.com/hulla-hoop/testSobes/internal/config"
	"github.com/hulla-hoop/testSobes/internal/kafkaEndpoint/consumer"
	"github.com/hulla-hoop/testSobes/internal/kafkaEndpoint/producer"
	"github.com/hulla-hoop/testSobes/internal/modeldb"
	"github.com/hulla-hoop/testSobes/internal/service"
	"github.com/hulla-hoop/testSobes/pkg/app"
	"github.com/joho/godotenv"
)

var wg sync.WaitGroup

func main() {

	//Инициализируем логеры
	infLogger := log.New(os.Stdout, "\nINFO:  ", log.Ldate|log.Lshortfile)
	errLogger := log.New(os.Stdout, "\nERROR:  ", log.Ldate|log.Lshortfile)

	//Загружаем .env содержащий все конфиги
	err := godotenv.Load()
	if err != nil {
		errLogger.Fatal("Не загружается .env файл")
	}

	//Инициализируем базу данных с gorm библиотекой
	/* sqlGorm := psql.InitDbGorm() */

	//Инициализируем базу данных с стандартной библиотекой
	db, err := psql.InitDb()
	if err != nil {
		errLogger.Println("Проблемы иниициализации БД", err)
	}

	rdb := rediscash.Init(db, infLogger, errLogger)

	//Инициализируем echo роутер и запскаем его
	a := app.New(rdb, infLogger, errLogger)

	go a.Start()

	//Получаем конфиги для сервиса и инициализируем новый сервис
	cfgApi := config.NewCfgApi()

	s := service.New(errLogger, cfgApi)

	//Создаем каналы для передачи данных между горутинами
	UserChan := make(chan modeldb.User)
	UserChanFailed := make(chan service.UserFailed)
	wg.Add(3)

	//Прописываем конфиги для kafka.consumer
	config := config.New()

	conf := make(kafka.ConfigMap)
	conf["bootstrap.servers"] = config.BootstrapService
	conf["group.id"] = config.GroupID
	conf["auto.offset.reset"] = config.AutoOffsetReset

	//Прописываем конфиги для kafka.producer
	conf2 := make(kafka.ConfigMap)
	conf2["bootstrap.servers"] = config.BootstrapService

	//Инициализируем kafka.producer
	p, err := kafka.NewProducer(&conf2)
	if err != nil {
		fmt.Printf("Failed to create producer: %s", err)
		os.Exit(1)
	}
	producer := producer.New(p, &wg, infLogger, errLogger, config)

	//Инициализируем kafka.consumer
	c, err := kafka.NewConsumer(&conf)
	if err != nil {
		fmt.Printf("Failed to create consumer: %s", err)
		os.Exit(1)
	}
	consumer := consumer.New(c, &wg, infLogger, errLogger, config)

	//Запускаем функции в горутинах
	go consumer.Consumer(UserChan)
	go service.Distribution(s, UserChan, UserChanFailed, infLogger, errLogger, &wg, db)
	go producer.Producer(UserChanFailed)

	wg.Wait()
	c.Close()

}
