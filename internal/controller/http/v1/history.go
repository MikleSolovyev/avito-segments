package v1

import (
	"context"
	"github.com/labstack/echo/v4"
	"io"
	"net/http"
)

type HistoryService interface {
	GetHistoryByDate(context.Context, int, int, io.Writer) error
}

type HistoryRoutes struct {
	service HistoryService
}

func newHistoryRoutes(router *echo.Group, service HistoryService) {
	r := &HistoryRoutes{
		service: service,
	}

	router.GET("", r.GetByDate)
}

// GetByDate
// @Summary      Get CSV report
// @Description  Get CSV report in a given period
// @Tags         history
// @Accept       json
// @Produce      json
// @Param        req    body   v1.HistoryRoutes.GetByDate.request  true  "from to period"
// @Success      200  {file}  file
// @Failure      400  {object}  echo.HTTPError
// @Failure      500  {object}  echo.HTTPError
// @Router       /api/v1/history [get]
func (r *HistoryRoutes) GetByDate(c echo.Context) error {
	type request struct {
		From int `json:"from" validate:"required"`
		To   int `json:"to" validate:"required"`
	}

	req := &request{}

	if err := c.Bind(req); err != nil {
		responseError(c, http.StatusBadRequest, err)
		return err
	}

	if err := c.Validate(req); err != nil {
		responseError(c, http.StatusBadRequest, err)
		return err
	}

	c.Response().Header().Set("Content-Type", "text/csv")
	c.Response().Header().Set("Content-Disposition", `attachment; filename="report.csv"`)
	err := r.service.GetHistoryByDate(c.Request().Context(), req.From, req.To, c.Response().Writer)
	if err != nil {
		responseError(c, http.StatusBadRequest, err)
	}
	return err
}
