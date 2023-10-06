package psql

import (
	"fmt"

	"github.com/hulla-hoop/testSobes/internal/config"
	_ "github.com/lib/pq"
)

func InitDb() {
	config := config.DbNew()

	dsn := fmt.Sprintf("host=%s user=%s dbname=%s password=%s port=%s sslmode=%s", config.Host, config.User, config.DBName, config.Password, config.Port, config.SSLMode)
	fmt.Print(dsn)
}
