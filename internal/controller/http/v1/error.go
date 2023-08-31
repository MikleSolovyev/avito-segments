package v1

import "github.com/labstack/echo/v4"

type errorResponse struct {
	Error string `json:"error"`
}

func responseError(c echo.Context, code int, err error) {
	c.JSON(code, errorResponse{Error: err.Error()})
}
