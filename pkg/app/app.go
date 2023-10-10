package app

import (
	"log"

	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/99designs/gqlgen/handler"
	"github.com/hulla-hoop/testSobes/internal/echoendpoint"
	"github.com/hulla-hoop/testSobes/internal/psql"
	"github.com/hulla-hoop/testSobes/internal/resolver"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type App struct {
	e         *echoendpoint.Endpoint
	echo      *echo.Echo
	psql      psql.DB
	inflogger *log.Logger
	errLogger *log.Logger
}

func New(db psql.DB, inflogger *log.Logger, errLogger *log.Logger) *App {
	a := App{}

	a.psql = db
	a.errLogger = errLogger
	a.inflogger = inflogger
	a.e = echoendpoint.New(a.psql, a.inflogger, a.errLogger)
	a.echo = echo.New()

	a.echo.Use(middleware.Logger())
	a.echo.Use(middleware.Recover())

	a.echo.POST("/graphql", func(c echo.Context) error {
		config := resolver.Config{
			Resolvers: &resolver.Resolver{
				DB: a.psql,
			},
		}
		h := handler.GraphQL(resolver.NewExecutableSchema(config))
		h.ServeHTTP(c.Response(), c.Request())
		return nil
	})
	a.echo.GET("/play", func(c echo.Context) error {
		h := playground.Handler("GraphQL playground", "/graphql")
		h.ServeHTTP(c.Response(), c.Request())
		return nil
	})

	a.echo.POST("/user", a.e.Insert)
	/* a.echo.GET("/userage", a.e.AgeSort) */
	a.echo.DELETE("/user/:id", a.e.Delete)
	a.echo.PUT("/user/:id", a.e.Update)
	/* a.echo.GET("/user/:nat", a.e.NatFilter)
	a.echo.GET("/user/:page", a.e.UserPagination) */

	return &a

}

func (a *App) Start() {
	a.inflogger.Println("Запуск сервера на localhost:1234")
	a.echo.Start(":1234")
}
