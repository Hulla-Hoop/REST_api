package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"regexp"
	"sync"
	"syscall"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type User struct {
	Name        string `json:"name"`
	Surname     string `json:"surname"`
	Patronymic  string `json:"patronymic"`
	Age         int    `json:"age"`
	Gender      string
	Nationality string
}

var Users []User

var wg sync.WaitGroup

func main() {
	UserChan := make(chan User)
	UserChanFailed := make(chan User)
	wg.Add(3)
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s <config-file-path>\n",
			os.Args[0])
		os.Exit(1)
	}

	configFile := os.Args[1]
	conf := ReadConfig(configFile)
	conf2 := ReadConfig(configFile)
	conf["group.id"] = "kafka-go-getting-started"
	conf["auto.offset.reset"] = "earliest"

	p, err := kafka.NewProducer(&conf2)

	if err != nil {
		fmt.Printf("Failed to create producer: %s", err)
		os.Exit(1)
	}

	c, err := kafka.NewConsumer(&conf)

	if err != nil {
		fmt.Printf("Failed to create consumer: %s", err)
		os.Exit(1)
	}

	go KafkaConsumer(c, UserChan)

	go Prinnnn(UserChan, UserChanFailed)
	go KafkaProducer(p, UserChanFailed)
	wg.Wait()
	c.Close()

}

func KafkaConsumer(c *kafka.Consumer, f chan User) {

	topic := "FIO"
	err := c.SubscribeTopics([]string{topic}, nil)
	if err != nil {
		os.Exit(1)
	}
	// Set up a channel for handling Ctrl-C, etc
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	// Process messages

	run := true
	for run {
		select {
		case sig := <-sigchan:
			fmt.Printf("Caught signal %v: terminating\n", sig)
			wg.Done()
			run = false
			close(f)
		default:
			ev, err := c.ReadMessage(100 * time.Millisecond)
			if err != nil {
				// Errors are informational and automatically handled by the consumer
				continue
			}
			var U User
			err = json.Unmarshal(ev.Key, &U)
			if err != nil {
				fmt.Println(err)
			}

			fmt.Println(U)
			f <- U
		}
	}

}

func Prinnnn(u chan User, uFailed chan User) {

	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	run := true
	for run {
		select {
		case sig := <-sigchan:
			fmt.Println("Ybito", sig)
			wg.Done()
			close(uFailed)
			run = false
		default:
			User1 := <-u
			chek := Check(User1)
			if chek {
				User2, err := EncrimentAge(User1)
				if err != nil {
					fmt.Println(err)
				}

				User2, err = EncrimentGender(User2)
				if err != nil {
					fmt.Println(err)
				}

				User2, err = EncrimentCountry(User2)
				if err != nil {
					fmt.Println(err)
				}
				fmt.Println("age encriment", User2)
			}

			fmt.Println(chek)
			/* uFailed <- <-u */
			/* fmt.Println(<-u) */
		}

	}

}

func KafkaProducer(p *kafka.Producer, user chan User) {

	topic := "FIO_FAILED"

	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	run := true
	for run {
		select {
		case sig := <-sigchan:
			fmt.Println("Ybito", sig)
			wg.Done()
			run = false
		default:
			uFailed, err := json.Marshal(<-user)
			if err != nil {
				fmt.Println(err)
			}
			key := uFailed
			p.Produce(&kafka.Message{
				TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
				Key:            []byte(key),
			}, nil)
		}
	}
}

func Check(U User) bool {
	if U.Name == "" || U.Surname == "" {
		return false
	}

	r, err := regexp.MatchString("^[a-zA-Z]+$", U.Name)
	if err != nil {
		fmt.Println(err)
	}
	if r == false {
		return r
	}
	r, err = regexp.MatchString("^[a-zA-Z]+$", U.Surname)
	if err != nil {
		fmt.Println(err)
	}
	if r == false {
		return r
	}

	return true

}

type Age struct {
	Count int    `json:"count"`
	Name  string `json:"name"`
	Age   int    `json:"age"`
}

func EncrimentAge(u User) (User, error) {
	userAge := Age{}
	url := (fmt.Sprintf("https://api.agify.io/?name=%s", u.Name))
	r, err := http.Get(url)
	if err != nil {
		return User{}, err
	}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return User{}, err
	}
	err = json.Unmarshal(body, &userAge)
	u.Age = userAge.Age
	return u, nil
}

type Gender struct {
	Count  int    `json:"count"`
	Name   string `json:"name"`
	Gender string `json:"gender"`
}

func EncrimentGender(u User) (User, error) {
	userGender := Gender{}
	url := (fmt.Sprintf("https://api.genderize.io/?name=%s", u.Name))
	r, err := http.Get(url)
	if err != nil {
		return User{}, err
	}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return User{}, err
	}
	err = json.Unmarshal(body, &userGender)
	u.Gender = userGender.Gender
	return u, nil
}

type Country struct {
	CountryId   string  `json:"country_id"`
	Probability float32 `json:"probability"`
}

type Natonality struct {
	Count   int    `json:"count"`
	Name    string `json:"name"`
	Country []Country
}

func EncrimentCountry(u User) (User, error) {
	userNati := Natonality{}
	url := (fmt.Sprintf("https://api.nationalize.io/?name=%s", u.Name))
	r, err := http.Get(url)
	if err != nil {
		return User{}, err
	}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return User{}, err
	}
	err = json.Unmarshal(body, &userNati)
	u.Nationality = userNati.Country[0].CountryId
	return u, nil
}
