package v1

import (
	_ "avito-segments/docs"
	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
	"net/http"
)

type Service interface {
	SegmentService
	UserService
	HistoryService
}

func New(service Service, logger echo.Logger, validator echo.Validator) *echo.Echo {
	router := echo.New()
	router.Logger = logger
	router.Validator = validator

	router.GET("/health", func(c echo.Context) error { return c.NoContent(http.StatusOK) })
	router.GET("/swagger/*", echoSwagger.WrapHandler)

	v1 := router.Group("/api/v1")
	{
		newSegmentRoutes(v1.Group("/segment"), service)
		newUserRoutes(v1.Group("/user"), service)
		newHistoryRoutes(v1.Group("/history"), service)
	}

	return router
}
