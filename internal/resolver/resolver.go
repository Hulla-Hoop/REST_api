package resolver

import (
	modelgql "github.com/hulla-hoop/testSobes/internal/modelgql"
	"github.com/hulla-hoop/testSobes/internal/psql"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	user []*modelgql.User
	DB   *psql.Psql
}
