package app

import (
	"crawler/api/v1/controllers"
	"crawler/utils"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"net/http"
)

func init() {

	//Echo instance
	e := echo.New()
	//Middleware
	//e.Use(middleware.Logger())
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "time=${time_rfc3339}, method='${method}', uri='${uri}', status=${status}, error='${error}', latency_human=${latency_human}, bytes_in:${bytes_in}, bytes_out:${bytes_out}\n",
	}))
	e.Use(middleware.Recover())
	e.Use(middleware.GzipWithConfig(middleware.GzipConfig{
		Level: 7,
	}))

	//Routes
	//group -> api v1
	v1 := e.Group("/api/v1", jsonMiddleware)
	v1.GET("/get/title", controllers.GetTitle)

	//Start server
	e.Logger.Fatal(e.Start(":8080"))

}

func jsonMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if c.Request().Header.Get("Content-Type") == "application/json" {
			//
			return next(c)
		}
		return c.JSON(http.StatusBadRequest, utils.Error("content type is not json"))
	}
}
