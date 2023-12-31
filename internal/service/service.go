package service

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
	"time"

	"github.com/hulla-hoop/testSobes/internal/config"
	"github.com/hulla-hoop/testSobes/internal/modeldb"
	"github.com/pkg/errors"
)

type Service struct {
	errLogger *log.Logger
	cfg       *config.ConfigApi
}

func New(errLogger *log.Logger, cfg *config.ConfigApi) *Service {
	return &Service{
		errLogger: errLogger,
		cfg:       cfg,
	}
}

type Age struct {
	Count int    `json:"count"`
	Name  string `json:"name"`
	Age   int    `json:"age"`
}

func (s *Service) EncrimentAge(uName string) (int, error) {
	userAge := Age{}
	url := (fmt.Sprintf(s.cfg.AGEAPI, uName))
	r, err := http.Get(url)
	if err != nil {
		s.errLogger.Println("Server is not available. Check connection", err)
		time.Sleep(5 * time.Second)
		age, err := s.EncrimentAge(uName)
		if err != nil {
			return 0, err
		}
		return age, nil
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		s.errLogger.Println(err)
		return 0, err
	}

	err = json.Unmarshal(body, &userAge)
	if err != nil {
		s.errLogger.Println(err)
		return 0, err
	}

	return userAge.Age, nil
}

type Gender struct {
	Count  int    `json:"count"`
	Name   string `json:"name"`
	Gender string `json:"gender"`
}

func (s *Service) EncrimentGender(uName string) (string, error) {
	userGender := Gender{}
	url := (fmt.Sprintf(s.cfg.GENDERAPI, uName))
	r, err := http.Get(url)
	if err != nil {
		s.errLogger.Println("Server is not available. Check connection")
		time.Sleep(5 * time.Second)
		name, err := s.EncrimentGender(uName)
		if err != nil {
			return "", err
		}
		return name, nil
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		return "", errors.Wrap(err, "Internal server error")
	}

	err = json.Unmarshal(body, &userGender)
	if err != nil {
		return "", errors.Wrap(err, "Internal server error")
	}

	return userGender.Gender, nil
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

func (s *Service) EncrimentCountry(uName string) (string, error) {
	userNati := Natonality{}
	url := (fmt.Sprintf(s.cfg.NATIONAPI, uName))
	r, err := http.Get(url)

	if err != nil {

		s.errLogger.Println("Server is not available. Check connection")
		time.Sleep(5 * time.Second)
		name, err := s.EncrimentCountry(uName)
		if err != nil {
			return "", err
		}
		return name, nil
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		s.errLogger.Println(err)
		return "", errors.Wrap(err, "Internal server error")
	}

	err = json.Unmarshal(body, &userNati)
	if err != nil {
		s.errLogger.Println(err)
		return "", errors.Wrap(err, "Internal server error")
	}

	return userNati.Country[0].CountryId, nil
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
		return "Неверный формат поля имя", r
	}

	r, err = regexp.MatchString("^[a-zA-Z]+$", U.Surname)
	if err != nil {
		fmt.Println(err)
	}
	if r == false {
		return "Неверный формат поля фамилия", r
	}
	if U.Patronymic == "" {
		return "", true
	} else {
		r, err = regexp.MatchString("^[a-zA-Z]+$", U.Patronymic)
		if err != nil {
			fmt.Println(err)
		}
		if r == false {
			return "Неверный формат поля отчество", r
		}
	}

	return "", true

}

func (s *Service) Encriment(u modeldb.User) (modeldb.User, error) {
	var err error
	u.Age, err = s.EncrimentAge(u.Name)
	if err != nil {
		return modeldb.User{}, err
	}

	u.Gender, err = s.EncrimentGender(u.Name)
	if err != nil {
		return modeldb.User{}, err
	}

	u.Nationality, err = s.EncrimentCountry(u.Name)
	if err != nil {
		return modeldb.User{}, err
	}

	u.CreatedAt = time.Now()
	u.UpdatedAt = time.Now()
	return u, nil
}
