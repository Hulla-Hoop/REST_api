// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package modelgql

import (
	"time"
)

type NewUser struct {
	Name        string  `json:"name"`
	Surname     string  `json:"surname"`
	Patronymic  *string `json:"patronymic,omitempty"`
	Age         *int    `json:"age,omitempty"`
	Gender      *string `json:"gender,omitempty"`
	Nationality *string `json:"nationality,omitempty"`
}

type User struct {
	ID          int       `json:"id"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
	Name        string    `json:"name"`
	Surname     string    `json:"surname"`
	Patronymic  string    `json:"patronymic"`
	Age         int       `json:"age"`
	Gender      string    `json:"gender"`
	Nationality string    `json:"nationality"`
}
