package service

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
	"sync"
	"time"

	"github.com/hulla-hoop/testSobes/internal/modeldb"
	"github.com/hulla-hoop/testSobes/internal/psql"
	"github.com/pkg/errors"
)

type Service struct {
	wg        *sync.WaitGroup
	db        *psql.Psql
	inflogger *log.Logger
	errLogger *log.Logger
}

func New(wg *sync.WaitGroup, db *psql.Psql, inflogger *log.Logger, errLogger *log.Logger) *Service {
	return &Service{
		wg:        wg,
		db:        db,
		inflogger: inflogger,
		errLogger: errLogger,
	}
}

type Age struct {
	Count int    `json:"count"`
	Name  string `json:"name"`
	Age   int    `json:"age"`
}

func (s *Service) EncrimentAge(u modeldb.User) (modeldb.User, error) {
	userAge := Age{}
	url := (fmt.Sprintf("https://api.agify.io/?name=%s", u.Name))
	r, err := http.Get(url)
	if err != nil {
		s.errLogger.Println("Server is not available. Check connection")
		time.Sleep(5 * time.Second)
		u, err = s.EncrimentAge(u)
		return u, nil
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		s.errLogger.Println(err)
		return modeldb.User{}, err
	}

	err = json.Unmarshal(body, &userAge)
	if err != nil {
		s.errLogger.Println(err)
		return modeldb.User{}, err
	}

	u.Age = userAge.Age
	return u, nil
}

type Gender struct {
	Count  int    `json:"count"`
	Name   string `json:"name"`
	Gender string `json:"gender"`
}

func (s *Service) EncrimentGender(u modeldb.User) (modeldb.User, error) {
	userGender := Gender{}
	url := (fmt.Sprintf("https://api.genderize.io/?name=%s", u.Name))
	r, err := http.Get(url)
	if err != nil {
		s.errLogger.Println("Server is not available. Check connection")
		time.Sleep(5 * time.Second)
		u, err = s.EncrimentGender(u)
		return u, nil
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		return modeldb.User{}, errors.Wrap(err, "Internal server error")
	}

	err = json.Unmarshal(body, &userGender)
	if err != nil {
		return modeldb.User{}, errors.Wrap(err, "Internal server error")
	}

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

func (s *Service) EncrimentCountry(u modeldb.User) (modeldb.User, error) {
	userNati := Natonality{}
	url := (fmt.Sprintf("https://api.nationalize.io/?name=%s", u.Name))
	r, err := http.Get(url)

	if err != nil {

		s.errLogger.Println("Server is not available. Check connection")
		time.Sleep(5 * time.Second)
		u, err = s.EncrimentCountry(u)
		return u, nil
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		s.errLogger.Println(err)
		return modeldb.User{}, err
	}

	err = json.Unmarshal(body, &userNati)
	if err != nil {
		s.errLogger.Println(err)
		return modeldb.User{}, err
	}

	u.Nationality = userNati.Country[0].CountryId
	return u, nil
}

func (s *Service) CheckErr(U modeldb.User) (string, bool) {
	if U.Name == "" || U.Surname == "" {

		return ("Нет обязательного поля"), false
	}

	r, err := regexp.MatchString("^[a-zA-Z]+$", U.Name)
	if err != nil {
		s.errLogger.Println(err)
	}
	if r == false {
		return ("Неверный формат"), r
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
