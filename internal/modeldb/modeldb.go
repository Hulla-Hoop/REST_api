package modeldb

import "time"

type User struct {
	Id          uint `gorm:"autoIncrement" gorm:"primary_key"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Name        string `json:"name"`
	Surname     string `json:"surname"`
	Patronymic  string `json:"patronymic"`
	Age         int    `json:"age"`
	Gender      string
	Nationality string
}
