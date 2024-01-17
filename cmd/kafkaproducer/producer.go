package main

import (
	"fmt"
	"math/rand"
	"os"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

func main() {

	/* c := config.New() */

	conf := make(kafka.ConfigMap)
	conf["bootstrap.servers"] = "localhost:9092"

	topic := "FIO"
	p, err := kafka.NewProducer(&conf)

	if err != nil {
		fmt.Printf("Failed to create producer: %s", err)
		os.Exit(1)
	}

	go func() {
		for e := range p.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error != nil {
					fmt.Printf("Failed to deliver message: %v\n", ev.TopicPartition)
				} else {
					fmt.Printf("Produced event to topic %s: key = %-10s value = %s\n",
						*ev.TopicPartition.Topic, string(ev.Key), string(ev.Value))
				}
			}
		}
	}()
	dmitriy := []byte(`{
		"name":"Dmitriy",
		"surname":"Ushakov"	
	}`)
	users := [][]byte{dmitriy}

	for n := 0; n < len(users); n++ {
		Value := users[rand.Intn(len(users))]
		p.Produce(&kafka.Message{
			TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
			Value:          []byte(Value),
		}, nil)
	}

	p.Flush(15 * 1000)
	p.Close()
}
