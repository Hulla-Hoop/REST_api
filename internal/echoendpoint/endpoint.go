package echoendpoint

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/hulla-hoop/testSobes/internal/models"
	"github.com/hulla-hoop/testSobes/internal/psql"
	"github.com/labstack/echo/v4"
)

type Endpoint struct {
	Db *psql.Psql
}

func New(db *psql.Psql) *Endpoint {
	return &Endpoint{Db: db}
}

func (e *Endpoint) Insert(c echo.Context) error {
	u := models.User{}
	err := c.Bind(u)
	if err != nil {
		fmt.Println(err)
	}
	e.Db.Create(u)
	return c.JSON(http.StatusCreated, u)
}

func (e *Endpoint) Delete(c echo.Context) error {
	id := c.Param("id")
	e.Db.Db.Delete(&models.User{}, id)
	return c.NoContent(http.StatusNoContent)
}

func (e *Endpoint) Update(c echo.Context) error {
	u := new(models.User)
	err := c.Bind(u)
	if err != nil {
		fmt.Println(err)
	}
	id, _ := strconv.Atoi(c.Param("id"))
	u.ID = uint(id)
	e.Db.Db.Save(u)
	return c.JSON(http.StatusOK, u)
}
