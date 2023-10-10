package psql

import (
	"database/sql"
	"fmt"

	"github.com/hulla-hoop/testSobes/internal/config"
	"github.com/hulla-hoop/testSobes/internal/modeldb"
	_ "github.com/lib/pq"
)

type DB interface {
	Create(user *modeldb.User) error
	Delete(id int) error
	InsertAll() ([]modeldb.User, error)
	Update(user *modeldb.User, id int) error
}

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

	return &sqlPostgres{
		dB: dB,
	}, nil

}
