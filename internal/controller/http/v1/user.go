package v1

import (
	"avito-segments/internal/model"
	"context"
	"github.com/labstack/echo/v4"
	"net/http"
)

type UserService interface {
	GetUserSegments(context.Context, *model.User) ([]string, error)
	UpdateUserSegments(context.Context, *model.User, map[string]int, []string) error
}

type UserRoutes struct {
	service UserService
}

func newUserRoutes(router *echo.Group, service UserService) {
	r := &UserRoutes{
		service: service,
	}

	router.GET("/:id", r.GetSegments)
	router.PUT("/:id", r.UpdateSegments)
}

// GetSegments
// @Summary      Get user segments
// @Description  Get all active user segments
// @Tags         user
// @Accept       json
// @Produce      json
// @Param        id    path    int  true  "user id"
// @Success      200  {object}  v1.UserRoutes.GetSegments.response
// @Failure      400  {object}  echo.HTTPError
// @Failure      500  {object}  echo.HTTPError
// @Router       /api/v1/user/{id} [get]
func (r *UserRoutes) GetSegments(c echo.Context) error {
	type request struct {
		Id int `param:"id" validate:"required,gte=1"`
	}

	type response struct {
		Slugs []string `json:"slugs"`
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

	slugs, err := r.service.GetUserSegments(c.Request().Context(), &model.User{Id: req.Id})
	if err != nil {
		responseError(c, http.StatusBadRequest, err)
		return err
	}

	return c.JSON(http.StatusOK, response{Slugs: slugs})
}

// UpdateSegments
// @Summary      Update user segments
// @Description  Add user to new segments and delete from current
// @Tags         user
// @Accept       json
// @Produce      json
// @Param        id    path    int  true  "user id"
// @Param        req    body   v1.UserRoutes.UpdateSegments.request  true  "segments to update"
// @Success      200  {object}  v1.UserRoutes.UpdateSegments.response
// @Failure      400  {object}  echo.HTTPError
// @Failure      500  {object}  echo.HTTPError
// @Router       /api/v1/user/{id} [put]
func (r *UserRoutes) UpdateSegments(c echo.Context) error {
	type request struct {
		Id          int            `param:"id" validate:"required,gte=1"`
		AddSlugs    map[string]int `json:"add_slugs,omitempty" validate:"omitempty,dive,keys,required,endkeys,gte=0"`
		DeleteSlugs []string       `json:"delete_slugs,omitempty" validate:"omitempty,dive,required"`
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

	err := r.service.UpdateUserSegments(
		c.Request().Context(),
		&model.User{Id: req.Id},
		req.AddSlugs,
		req.DeleteSlugs,
	)
	if err != nil {
		responseError(c, http.StatusBadRequest, err)
		return err
	}

	return c.JSON(http.StatusOK, response{Message: "Success"})
}
