package resolver

import (
	"github.com/hulla-hoop/testSobes/internal/DB"
	modelgql "github.com/hulla-hoop/testSobes/internal/modelgql"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	user []*modelgql.User
	DB   DB.DB
}
