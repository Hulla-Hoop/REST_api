package psql

import (
	"errors"
	"fmt"
	"regexp"

	"github.com/hulla-hoop/testSobes/internal/modeldb"
)

func (p *Psqlgorm) Create(U *modeldb.User) error {
	if U.Name == "" || U.Surname == "" {

		return errors.New("Нет обязательного поля")
	}

	r, err := regexp.MatchString("^[a-zA-Z]+$", U.Name)
	if err != nil {
		return err
	}

	if r == false {
		return errors.New("Неверный формат")
	}

	r, err = regexp.MatchString("^[a-zA-Z]+$", U.Surname)
	if err != nil {
		fmt.Println(err)
	}

	if r == false {
		return errors.New("Неверный формат")
	}

	p.Db.Create(&U)

	return nil

}

func (p *Psqlgorm) Deleate(id int) error {
	return nil
}
func (p *Psqlgorm) InsertAll(id int) ([]modeldb.User, error) {
	return nil, nil
}
func (p *Psqlgorm) Update(user *modeldb.User, id int) error {
	return nil
}
