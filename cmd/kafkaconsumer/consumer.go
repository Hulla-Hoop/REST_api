package main

import (
	"sync"

	"github.com/hulla-hoop/testSobes/internal/echoendpoint"
	"github.com/hulla-hoop/testSobes/internal/psql"
	"github.com/labstack/echo/v4"
)

var wg sync.WaitGroup

func main() {

	p := psql.InitDb()
	e := echo.New()
	end := echoendpoint.New(p)

	e.POST("/user", end.Insert)
	e.GET("/userage", end.AgeSort)
	e.DELETE("/user/:id", end.Delete)
	e.PUT("/user/:id", end.Update)
	e.GET("/user/:nat", end.NatFilter)
	e.GET("/user/:page", end.UserPagination)

	e.Start(":1234")
	/* s := service.New(&wg)

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
	c.Close() */

}
