package psql

import (
	"database/sql"
	"fmt"

	"github.com/hulla-hoop/testSobes/internal/config"
	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
)

type sqlPostgres struct {
	dB *sql.DB
}

func InitDb() (*sqlPostgres, error) {
	config := config.DbNew()

	dsn := fmt.Sprintf("host=%s user=%s dbname=%s password=%s port=%s sslmode=%s", config.Host, config.User, config.DBName, config.Password, config.Port, config.SSLMode)

	dB, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}
	err = goose.Up(dB, "db")
	if err != nil {
		fmt.Println(err)
	}
	return &sqlPostgres{
		dB: dB,
	}, nil

}
