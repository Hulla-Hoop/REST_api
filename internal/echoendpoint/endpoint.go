package echoendpoint

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/hulla-hoop/testSobes/internal/modeldb"
	"github.com/hulla-hoop/testSobes/internal/psql"
	"github.com/labstack/echo/v4"
)

type Endpoint struct {
	Db        psql.DB
	inflogger *log.Logger
	errLogger *log.Logger
}

func New(db psql.DB, inflogger *log.Logger, errLogger *log.Logger) *Endpoint {
	return &Endpoint{Db: db,
		inflogger: inflogger,
		errLogger: errLogger}
}

func (e *Endpoint) Insert(c echo.Context) error {
	u := modeldb.User{}
	err := c.Bind(&u)
	if err != nil {
		e.errLogger.Println(err)
		return c.JSON(http.StatusInternalServerError, nil)
	}
	u.CreatedAt = time.Now()
	u.UpdatedAt = time.Now()
	err = e.Db.Create(&u)
	if err != nil {
		e.errLogger.Println(err)
		return c.JSON(http.StatusInternalServerError, nil)
	}
	return c.JSON(http.StatusCreated, u)
}

func (e *Endpoint) Delete(c echo.Context) error {
	id := c.Param("id")
	idi, err := strconv.Atoi(id)
	if err != nil {
		e.errLogger.Println(err)
	}
	err = e.Db.Delete(idi)
	if err != nil {
		e.errLogger.Println(err)
		return c.JSON(http.StatusInternalServerError, nil)
	}
	return c.NoContent(http.StatusNoContent)
}

func (e *Endpoint) Update(c echo.Context) error {
	u := new(modeldb.User)
	err := c.Bind(u)
	if err != nil {
		e.errLogger.Println(err)
		return c.JSON(http.StatusInternalServerError, nil)
	}
	id, _ := strconv.Atoi(c.Param("id"))
	err = e.Db.Update(u, id)
	if err != nil {
		e.errLogger.Println(err)
		return c.JSON(http.StatusInternalServerError, nil)
	}
	return c.JSON(http.StatusOK, u)
}

func (e *Endpoint) Sort(c echo.Context) error {
	users := []modeldb.User{}

	valueStr, err := c.FormParams()
	if err != nil {
		e.errLogger.Println(err)
		return c.JSON(http.StatusInternalServerError, nil)
	}

	field := valueStr["sort"]
	e.inflogger.Println(field[0])

	users, err = e.Db.Sort(field[0])
	if err != nil {
		e.errLogger.Println(err)
		return c.JSON(http.StatusInternalServerError, nil)
	}

	return c.JSON(http.StatusOK, users)
}

func (e *Endpoint) UserPagination(c echo.Context) error {
	valueStr, err := c.FormParams()

	pageStr := valueStr["page"]
	limitStr := valueStr["limit"]

	e.inflogger.Println(pageStr, limitStr)

	page, err := strconv.Atoi(pageStr[0])
	if err != nil {
		return c.JSON(http.StatusInternalServerError, nil)
	}
	limit, err := strconv.Atoi(limitStr[0])
	if err != nil {
		return c.JSON(http.StatusInternalServerError, nil)
	}

	u, err := e.Db.InsertPage(uint(page), limit)
	if err != nil {
		e.errLogger.Println(err)
		return c.JSON(http.StatusInternalServerError, nil)

	}

	return c.JSON(http.StatusOK, u)

}
