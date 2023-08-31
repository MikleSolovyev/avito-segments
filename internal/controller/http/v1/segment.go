package v1

import (
	"avito-segments/internal/model"
	"context"
	"github.com/labstack/echo/v4"
	"net/http"
)

type SegmentService interface {
	CreateSegment(context.Context, *model.Segment) error
	DeleteSegmentBySlug(context.Context, string) error
}

type SegmentRoutes struct {
	service SegmentService
}

func newSegmentRoutes(router *echo.Group, service SegmentService) {
	r := &SegmentRoutes{
		service: service,
	}

	router.POST("", r.Create)
	router.DELETE("", r.Delete)
}

// Create
// @Summary      Create segment
// @Description  Create segment with slug and percent
// @Tags         segment
// @Accept       json
// @Produce      json
// @Param        req    body   v1.SegmentRoutes.Create.request  true  "slug and percent"
// @Success      200  {object}  v1.SegmentRoutes.Create.response
// @Failure      400  {object}  echo.HTTPError
// @Failure      500  {object}  echo.HTTPError
// @Router       /api/v1/segment [post]
func (r *SegmentRoutes) Create(c echo.Context) error {
	type request struct {
		Slug    string `json:"slug" validate:"required"`
		Percent int    `json:"percent,omitempty" validate:"omitempty,gte=0,lte=100"`
	}

	type response struct {
		Id int `json:"id"`
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

	segment := &model.Segment{
		Slug:    req.Slug,
		Percent: req.Percent,
	}

	if err := r.service.CreateSegment(c.Request().Context(), segment); err != nil {
		responseError(c, http.StatusBadRequest, err)
		return err
	}

	return c.JSON(http.StatusOK, response{
		Id: segment.Id,
	})
}

// Delete
// @Summary      Delete segment
// @Description  Delete segment by slug
// @Tags         segment
// @Accept       json
// @Produce      json
// @Param        req    body   v1.SegmentRoutes.Delete.request  true  "slug"
// @Success      200  {object}  v1.SegmentRoutes.Delete.response
// @Failure      400  {object}  echo.HTTPError
// @Failure      500  {object}  echo.HTTPError
// @Router       /api/v1/segment [delete]
func (r *SegmentRoutes) Delete(c echo.Context) error {
	type request struct {
		Slug string `json:"slug" validate:"required"`
	}

	type response struct {
		Message string `json:"message"`
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

	if err := r.service.DeleteSegmentBySlug(c.Request().Context(), req.Slug); err != nil {
		responseError(c, http.StatusBadRequest, err)
		return err
	}

	return c.JSON(http.StatusOK, response{
		Message: "Success",
	})
}
