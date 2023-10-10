package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.38

import (
	"context"
	"time"

	"github.com/hulla-hoop/testSobes/internal/modeldb"
	"github.com/hulla-hoop/testSobes/internal/modelgql"
)

// CreateUser is the resolver for the createUser field.
func (r *mutationResolver) CreateUser(ctx context.Context, input modelgql.NewUser) (*modelgql.User, error) {
	AddUser := modelgql.User{
		Name:    input.Name,
		Surname: input.Surname,
		Age:     *input.Age,
	}

	user := modeldb.User{
		Name:    AddUser.Name,
		Surname: AddUser.Surname,
		Age:     AddUser.Age,
	}

	err := r.DB.Create(&user)
	if err != nil {
		return nil, err
	}

	return &AddUser, nil
}

// UpdateUser is the resolver for the updateUser field.
func (r *mutationResolver) UpdateUser(ctx context.Context, userID int, input *modelgql.NewUser) (*modelgql.User, error) {
	UpdateUser := modelgql.User{
		Name:        input.Name,
		Surname:     input.Surname,
		Patronymic:  *input.Patronymic,
		Age:         *input.Age,
		Gender:      *input.Gender,
		Nationality: *input.Nationality,
	}
	user := modeldb.User{
		UpdatedAt:   time.Now(),
		Name:        UpdateUser.Name,
		Surname:     UpdateUser.Surname,
		Patronymic:  UpdateUser.Patronymic,
		Age:         UpdateUser.Age,
		Gender:      UpdateUser.Gender,
		Nationality: UpdateUser.Nationality,
	}
	err := r.DB.Update(&user, userID)
	if err != nil {
		return nil, err
	}
	UpdateUser.ID = userID
	return &UpdateUser, nil
}

// Users is the resolver for the users field.
func (r *queryResolver) Users(ctx context.Context) ([]*modelgql.User, error) {
	f := []*modelgql.User{}
	u := []modeldb.User{}

	u, err := r.DB.InsertAll()
	if err != nil {
		return nil, err
	}

	for _, t := range u {
		f = append(f, &modelgql.User{
			Name:        t.Name,
			Surname:     t.Surname,
			Patronymic:  t.Patronymic,
			Age:         t.Age,
			Gender:      t.Gender,
			Nationality: t.Nationality,
		})
	}

	return f, nil
}

// Pages is the resolver for the pages field.
func (r *queryResolver) Pages(ctx context.Context, page int) ([]*modelgql.User, error) {
	/* var UserCount int */

	/* err := r.DB.Db.Table("users").Count(&UserCount).Error
	if err != nil {
		return nil, err
	}

	UserPerPage := 3

	pageCount := int(math.Ceil(float64(UserCount) / float64(UserPerPage)))

	if pageCount == 0 {
		pageCount = 1
	}
	if page > pageCount {

		return nil, err

	}

	offset := (page - 1) * UserPerPage */

	users := []*modelgql.User{}

	/* err = r.DB.Db.Limit(UserPerPage).Offset(offset).Find(&users).Error
	if err != nil {
		return nil, err
	} */

	return users, nil
}

// Mutation returns MutationResolver implementation.
func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
