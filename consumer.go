package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type User struct {
	Name       string `json:"name"`
	Surname    string `json:"surname"`
	Patronymic string `json:"patronymic"`
}

var Users []User

var wg sync.WaitGroup

func main() {
	UserChan := make(chan []User)
	wg.Add(1)
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s <config-file-path>\n",
			os.Args[0])
		os.Exit(1)
	}

	configFile := os.Args[1]
	conf := ReadConfig(configFile)
	conf["group.id"] = "kafka-go-getting-started"
	conf["auto.offset.reset"] = "earliest"

	c, err := kafka.NewConsumer(&conf)

	if err != nil {
		fmt.Printf("Failed to create consumer: %s", err)
		os.Exit(1)
	}

	go KafkaConsumer(c, UserChan)

	US := <-UserChan
	Prinnnn(US)
	wg.Wait()
	c.Close()

}

func KafkaConsumer(c *kafka.Consumer, f chan []User) {

	topic := "FIO"
	err := c.SubscribeTopics([]string{topic}, nil)
	if err != nil {
		os.Exit(1)
	}
	// Set up a channel for handling Ctrl-C, etc
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	var Users []User

	// Process messages

	run := true
	for run {
		select {
		case sig := <-sigchan:
			fmt.Printf("Caught signal %v: terminating\n", sig)
			wg.Done()
			run = false
		default:
			ev, err := c.ReadMessage(100 * time.Millisecond)
			if err != nil {
				// Errors are informational and automatically handled by the consumer
				continue
			}
			var U User
			err = json.Unmarshal(ev.Key, &U)
			Users = append(Users, U)
			fmt.Println(U)
		}
	}
	f <- Users
}

func Prinnnn(u []User) {
	fmt.Println(u)
}
