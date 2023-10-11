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

/* func (e *Endpoint) AgeSort(c echo.Context) error {
	users := []modeldb.User{}
	e.Db.Db.Raw("SELECT * FROM users WHERE deleted_at IS NULL ORDER BY age").Scan(&users)
	return c.JSON(http.StatusOK, users)
}

func (e *Endpoint) NatFilter(c echo.Context) error {
	national, _ := strconv.Atoi(c.Param("nat"))
	users := []modeldb.User{}
	e.Db.Db.Raw("SELECT * FROM users WHERE age = ?", national).Scan(&users)
	return c.JSON(http.StatusOK, users)
} */

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

	/*var UserCount int

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

	return c.JSON(http.StatusInternalServerError, users) */
}
