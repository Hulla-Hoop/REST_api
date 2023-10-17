package rediscash

import "github.com/redis/go-redis/v9"

func Init() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     "localhost",
		Password: "",
		DB:       0,
	})
}
