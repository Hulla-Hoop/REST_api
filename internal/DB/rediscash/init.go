package rediscash

import (
	"log"
	"strconv"

	"github.com/hulla-hoop/testSobes/internal/DB"
	"github.com/hulla-hoop/testSobes/internal/config"
	"github.com/redis/go-redis/v9"
)

type Cash struct {
	r       *redis.Client
	db      DB.DB
	infoLog *log.Logger
	errLog  *log.Logger
}

func Init(db DB.DB, i *log.Logger, e *log.Logger) *Cash {
	c := config.NewCfgRedis()
	cDB, err := strconv.Atoi(c.DB)
	if err != nil {
		e.Println("Невозможно конвертировать конфиг", err)
	}
	d := redis.NewClient(&redis.Options{
		Addr:     c.HOST + ":" + c.PORT,
		Password: c.PASSWORD,
		DB:       cDB,
	})

	return &Cash{
		r:       d,
		db:      db,
		infoLog: i,
		errLog:  e,
	}
}
