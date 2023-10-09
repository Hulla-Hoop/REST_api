package psql

import (
	"database/sql"
	"fmt"

	"github.com/hulla-hoop/testSobes/internal/modeldb"
)

func (db *sqlPostgres) Create(user *modeldb.User) error {

	var id int
	err := db.dB.QueryRow("insert into users(created_at,updated_at,name,surname,patronymic,age,gender,nationality)  values ($1, $2,$3,$4,$5,$6,$7,$8) returning id",
		user.CreatedAt,
		user.UpdatedAt,
		user.Name,
		user.Surname,
		user.Patronymic,
		user.Age,
		user.Gender,
		user.Nationality).Scan(&id)
	if err != nil {

		switch err {
		case sql.ErrNoRows:
			return fmt.Errorf("Пользователь добавлен но не удалось записать ID %s", err)
		default:
			return fmt.Errorf("Ошибка при создании пользователя %s", err)
		}
	}

	fmt.Println("id созданого пользователя", id)
	return nil

}

func (db *sqlPostgres) Update(user *modeldb.User, id int) error {
	result, err := db.dB.Exec("update users set created_at = $1,updated_at=$2,name=$3,surname=$4,patronymic=$5,age=$6,gender=$7,nationality=$8 where id=$9",
		user.CreatedAt,
		user.UpdatedAt,
		user.Name,
		user.Surname,
		user.Patronymic,
		user.Age,
		user.Gender,
		user.Nationality, id)
	if err != nil {
		return err
	}
	fmt.Println(result.RowsAffected())
	return nil
}

func (db *sqlPostgres) Deleate(id int) error {
	result, err := db.dB.Exec("delete from users where id = $1", id)
	if err != nil {
		return err
	}
	fmt.Println(result.RowsAffected())
	return nil
}

func (db *sqlPostgres) InsertAll(id int) ([]modeldb.User, error) {
	rows, err := db.dB.Query("select * from users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	user := []modeldb.User{}

	for rows.Next() {
		u := modeldb.User{}
		err := rows.Scan(&u.Id, &u.CreatedAt, &u.UpdatedAt, &u.Name, &u.Surname, &u.Patronymic, &u.Age, &u.Gender, &u.Nationality)
		if err != nil {
			fmt.Println(err)
			continue
		}
		user = append(user, u)
	}

	return user, nil
}
