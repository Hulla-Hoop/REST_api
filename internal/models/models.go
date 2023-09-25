package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Name        string `json:"name"`
	Surname     string `json:"surname"`
	Patronymic  string `json:"patronymic"`
	Age         int    `json:"age"`
	Gender      string
	Nationality string
}
