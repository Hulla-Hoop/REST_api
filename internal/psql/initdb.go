package psql

import (
	"fmt"

	"github.com/hulla-hoop/testSobes/internal/config"
	"github.com/hulla-hoop/testSobes/internal/service"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

var db *gorm.DB

func InitDb() {

	config := config.DbNew()

	dsn := fmt.Sprintf("host=%s user=%s dbname=%s password=%s port=%s sslmode=%s", config.Host, config.User, config.DBName, config.Password, config.Port, config.SSLMode)
	fmt.Println(dsn)
	db, err := gorm.Open("postgres", dsn)
	if err != nil {
		fmt.Println(err)
	}

	db.Debug().AutoMigrate(&service.User{})

}

func New() *gorm.DB {
	return db
}
