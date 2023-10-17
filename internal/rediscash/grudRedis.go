package rediscash

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hulla-hoop/testSobes/internal/modeldb"
)

/*
	// interface for redis

	type DB interface {
	Create(user *modeldb.User) error
	Delete(id int) error
	InsertAll() ([]modeldb.User, error)
	Update(user *modeldb.User, id int) error
	InsertPage(page uint, limit int) ([]modeldb.User, error)
	Sort(field string) ([]modeldb.User, error)
	Filter(field string, operator string, value string) ([]modeldb.User, error)
} */

func (c *Cash) Create(user *modeldb.User) error {
	err := c.db.Create(user)
	if err != nil {
		return err
	}
	return nil
}

func (c *Cash) Delete(id int) error {
	err := c.db.Delete(id)
	if err != nil {
		return err
	}
	return nil
}

func (c *Cash) Update(user *modeldb.User, id int) error {
	err := c.db.Update(user, id)
	if err != nil {
		return err
	}
	return nil
}

func (c *Cash) InsertAll() ([]modeldb.User, error) {

	users, err := c.db.InsertAll()
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (c *Cash) InsertPage(page uint, limit int) (modeldb.Users, error) {

	key := strconv.Itoa(int(page)) + strconv.Itoa(limit)

	ctx := context.Background()

	v, err := c.r.Get(ctx, key).Result()
	if v == "" {
		fmt.Println("Ваяяяяя закинули")
	} else {

		fmt.Println(v)
	}
	if err != nil {
		fmt.Println(err)
	}

	users, err := c.db.InsertPage(page, limit)
	if err != nil {
		return nil, err
	}

	err = c.r.Set(ctx, key, users, 0).Err()
	if err != nil {
		fmt.Println(err)
	}

	return users, nil
}

func (c *Cash) Sort(field string) ([]modeldb.User, error) {
	users, err := c.db.Sort(field)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (c *Cash) Filter(field string, operator string, value string) ([]modeldb.User, error) {

	users, err := c.db.Filter(field, operator, value)
	if err != nil {
		return nil, err
	}
	return users, nil
}
