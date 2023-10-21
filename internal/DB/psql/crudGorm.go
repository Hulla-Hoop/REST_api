package psql

import (
	"errors"
	"fmt"
	"math"
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

func (p *Psqlgorm) Delete(id int) error {

	p.Db.Delete(id)
	return nil
}
func (p *Psqlgorm) InsertAll() ([]modeldb.User, error) {

	return nil, nil
}
func (p *Psqlgorm) Update(user *modeldb.User, id int) error {
	user.Id = uint(id)
	p.Db.Save(user)
	return nil
}

func (p *Psqlgorm) InsertPage(page uint, limit int) ([]modeldb.User, error) {

	var UserCount int

	err := p.Db.Table("users").Count(&UserCount).Error
	if err != nil {
		return nil, err
	}

	UserPerPage := 3

	pageCount := int(math.Ceil(float64(UserCount) / float64(UserPerPage)))

	if pageCount == 0 {
		pageCount = 1
	}
	if int(page) > pageCount {

		return nil, err

	}

	offset := (int(page) - 1) * UserPerPage

	users := []modeldb.User{}

	err = p.Db.Limit(UserPerPage).Offset(offset).Find(&users).Error
	if err != nil {
		return nil, err
	}

	return users, nil
}
func (p *Psqlgorm) Sort(field string) ([]modeldb.User, error) {
	var users []modeldb.User
	p.Db.Raw("SELECT * FROM users WHERE deleted_at IS NULL ORDER BY %s", field).Scan(&users)
	return users, nil
}
