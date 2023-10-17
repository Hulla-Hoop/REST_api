package rediscash

import (
	"github.com/hulla-hoop/testSobes/internal/psql"
	"github.com/redis/go-redis/v9"
)

type Cash struct {
	r  *redis.Client
	db psql.DB
}

func Init(db psql.DB) *Cash {
	d := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	return &Cash{
		r:  d,
		db: db,
	}
}
