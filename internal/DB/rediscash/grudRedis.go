package rediscash

import (
	"context"
	"encoding/json"
	"strconv"
	"time"

	"github.com/hulla-hoop/testSobes/internal/modeldb"
)

/*
	// Интерфейс который должен реализовать редис

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

	v, err := c.r.Get(ctx, key).Bytes()
	if v == nil {
		c.infoLog.Println("Значения нет в Кэше")
	} else if err != nil {
		c.errLog.Println(err)
	} else {
		c.infoLog.Println("Значения взято из Кэша")
		Users := modeldb.Users{}
		err = json.Unmarshal(v, &Users)
		if err != nil {
			return nil, err
		}
		return Users, nil
	}

	users, err := c.db.InsertPage(page, limit)
	if err != nil {
		return nil, err
	}

	err = c.r.Set(ctx, key, users, time.Minute).Err()
	if err != nil {
		c.errLog.Println(err)
	}

	return users, nil
}

func (c *Cash) Sort(field string) ([]modeldb.User, error) {

	key := field

	ctx := context.Background()

	v, err := c.r.Get(ctx, key).Bytes()
	if v == nil {
		c.infoLog.Println("Значения нет в Кэше")
	} else if err != nil {
		c.errLog.Println(err)
	} else {
		c.infoLog.Println("Значения взято из Кэша")
		Users := modeldb.Users{}
		err = json.Unmarshal(v, &Users)
		if err != nil {
			return nil, err
		}
		return Users, nil
	}

	users, err := c.db.Sort(field)
	if err != nil {
		return nil, err
	}

	err = c.r.Set(ctx, key, users, time.Minute).Err()
	if err != nil {
		c.errLog.Println(err)
	}

	return users, nil
}

func (c *Cash) Filter(field string, operator string, value string) ([]modeldb.User, error) {

	key := field + operator + value

	ctx := context.Background()

	v, err := c.r.Get(ctx, key).Bytes()
	if v == nil {
		c.infoLog.Println("Значения нет в Кэше")
	} else if err != nil {
		c.errLog.Println(err)
	} else {
		c.infoLog.Println("Значения взято из Кэша")
		Users := modeldb.Users{}
		err = json.Unmarshal(v, &Users)
		if err != nil {
			return nil, err
		}
		return Users, nil
	}

	users, err := c.db.Filter(field, operator, value)
	if err != nil {
		return nil, err
	}

	err = c.r.Set(ctx, key, users, time.Minute).Err()
	if err != nil {
		c.errLog.Println(err)
	}

	return users, nil
}
