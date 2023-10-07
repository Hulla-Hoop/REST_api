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

/* func (db *sqlPostgres) Update(user *modeldb.User) error {


} */
