package psql

import (
	"fmt"

	"github.com/hulla-hoop/testSobes/internal/models"
)

func (p *Psql) Create(u models.User) {

	result := p.Db.Create(&u)
	fmt.Println(u.ID)
	fmt.Println(result.Error)
	fmt.Println(result.RowsAffected)
}
