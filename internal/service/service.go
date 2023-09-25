package service

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"sync"

	"github.com/hulla-hoop/testSobes/internal/models"
	"github.com/hulla-hoop/testSobes/internal/psql"
)

type Service struct {
	wg *sync.WaitGroup
	db *psql.Psql
}

func New(wg *sync.WaitGroup) *Service {
	return &Service{
		wg: wg,
	}
}

type Age struct {
	Count int    `json:"count"`
	Name  string `json:"name"`
	Age   int    `json:"age"`
}

func (s *Service) EncrimentAge(u models.User) (models.User, error) {
	userAge := Age{}
	url := (fmt.Sprintf("https://api.agify.io/?name=%s", u.Name))
	r, err := http.Get(url)
	if err != nil {
		return models.User{}, err
	}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return models.User{}, err
	}
	err = json.Unmarshal(body, &userAge)
	u.Age = userAge.Age
	return u, nil
}

type Gender struct {
	Count  int    `json:"count"`
	Name   string `json:"name"`
	Gender string `json:"gender"`
}

func (s *Service) EncrimentGender(u models.User) (models.User, error) {
	userGender := Gender{}
	url := (fmt.Sprintf("https://api.genderize.io/?name=%s", u.Name))
	r, err := http.Get(url)
	if err != nil {
		return models.User{}, err
	}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return models.User{}, err
	}
	err = json.Unmarshal(body, &userGender)
	u.Gender = userGender.Gender
	return u, nil
}

type Country struct {
	CountryId   string  `json:"country_id"`
	Probability float32 `json:"probability"`
}

type Natonality struct {
	Count   int    `json:"count"`
	Name    string `json:"name"`
	Country []Country
}

func (s *Service) EncrimentCountry(u models.User) (models.User, error) {
	userNati := Natonality{}
	url := (fmt.Sprintf("https://api.nationalize.io/?name=%s", u.Name))
	r, err := http.Get(url)
	if err != nil {
		return models.User{}, err
	}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return models.User{}, err
	}
	err = json.Unmarshal(body, &userNati)
	u.Nationality = userNati.Country[0].CountryId
	return u, nil
}

func (s *Service) CheckErr(U models.User) (string, bool) {
	if U.Name == "" || U.Surname == "" {
		return "Нет обязательного поля", false
	}

	r, err := regexp.MatchString("^[a-zA-Z]+$", U.Name)
	if err != nil {
		fmt.Println(err)
	}
	if r == false {
		return "Неверный формат", r
	}
	r, err = regexp.MatchString("^[a-zA-Z]+$", U.Surname)
	if err != nil {
		fmt.Println(err)
	}
	if r == false {
		return "Неверный формат", r
	}

	return "", true

}
