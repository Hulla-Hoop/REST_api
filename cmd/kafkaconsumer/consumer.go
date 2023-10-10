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
	//Загружаем .env содержащий все конфиги
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	//Инициализируем логеры
	infLogger := log.New(os.Stdout, "INFO:  ", log.Ldate|log.Lshortfile)
	errLogger := log.New(os.Stdout, "ERROR:  ", log.Ldate|log.Lshortfile)

	//Инициализируем базу данных с gorm библиотекой
	/* sqlGorm := psql.InitDbGorm() */

	//Инициализируем базу данных с стандартной библиотекой
	db, err := psql.InitDb()
	if err != nil {
		fmt.Println("пиздец нахуя блять")
	}

	//Инициализируем echo роутер и запскаем его
	a := app.New(db, infLogger, errLogger)

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

	//Запускаем функции в горутинах для беспрерывной работы программы
	go consumer.Consumer(UserChan)
	go service.Distribution(s, UserChan, UserChanFailed, infLogger, errLogger, &wg, db)
	go producer.Producer(UserChanFailed)

	wg.Wait()
	c.Close()

}
