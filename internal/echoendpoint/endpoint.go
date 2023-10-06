package echoendpoint

import (
	"log"
	"math"
	"net/http"
	"strconv"

	"github.com/hulla-hoop/testSobes/internal/modeldb"
	"github.com/hulla-hoop/testSobes/internal/psql"
	"github.com/labstack/echo/v4"
)

type Endpoint struct {
	Db        *psql.Psqlgorm
	inflogger *log.Logger
	errLogger *log.Logger
}

func New(db *psql.Psqlgorm, inflogger *log.Logger, errLogger *log.Logger) *Endpoint {
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
	e.Db.Create(u)
	return c.JSON(http.StatusCreated, u)
}

func (e *Endpoint) Delete(c echo.Context) error {
	id := c.Param("id")
	e.Db.Db.Delete(&modeldb.User{}, id)
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
	u.Id = uint(id)
	e.Db.Db.Save(u)
	return c.JSON(http.StatusOK, u)
}

func (e *Endpoint) AgeSort(c echo.Context) error {
	users := []modeldb.User{}
	e.Db.Db.Raw("SELECT * FROM users WHERE deleted_at IS NULL ORDER BY age").Scan(&users)
	return c.JSON(http.StatusOK, users)
}

func (e *Endpoint) NatFilter(c echo.Context) error {
	national, _ := strconv.Atoi(c.Param("nat"))
	users := []modeldb.User{}
	e.Db.Db.Raw("SELECT * FROM users WHERE age = ?", national).Scan(&users)
	return c.JSON(http.StatusOK, users)
}

func (e *Endpoint) UserPagination(c echo.Context) error {
	pageStr := c.Param("page")

	page, err := strconv.Atoi(pageStr)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, nil)
	}

	var UserCount int

	err = e.Db.Db.Table("users").Count(&UserCount).Error
	if err != nil {
		return c.JSON(http.StatusInternalServerError, nil)
	}

	UserPerPage := 3

	pageCount := int(math.Ceil(float64(UserCount) / float64(UserPerPage)))

	if pageCount == 0 {
		pageCount = 1
	}
	if page > pageCount {

		return c.JSON(http.StatusInternalServerError, nil)

	}

	offset := (page - 1) * UserPerPage

	users := []modeldb.User{}

	err = e.Db.Db.Limit(UserPerPage).Offset(offset).Find(&users).Error
	if err != nil {
		return c.JSON(http.StatusInternalServerError, nil)
	}

	return c.JSON(http.StatusInternalServerError, users)
}
