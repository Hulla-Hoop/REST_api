package psql

import (
	"github.com/hulla-hoop/testSobes/internal/modeldb"
)

func (p *Psql) Create(u modeldb.User) {

	p.Db.Create(&u)

}
