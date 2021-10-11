package rest

import (
	"context"
	"net/http"

	"github.com/afif0808/sagara-test/errors"
	"github.com/afif0808/sagara-test/internal/domain"
	"github.com/afif0808/sagara-test/wrapper"
	"github.com/asaskevich/govalidator"
	"github.com/labstack/echo"
)

type usecase interface {
	CreateUser(ctx context.Context, u *domain.User) (err error)
}

type UserRestHandler struct {
	uc usecase
}

func NewUserRestHandler(uc usecase) UserRestHandler {
	return UserRestHandler{
		uc: uc,
	}
}

func (urh UserRestHandler) Mount(root *echo.Group) {
	user := root.Group("/user/")
	user.POST("", urh.createUser)
}

func (urh UserRestHandler) createUser(c echo.Context) error {
	var payload struct {
		Email    string `json:"email" valid:"required"`
		Name     string `json:"name" valid:"required"`
		Password string `json:"password" valid:"required"`
	}
	if err := getPayload(c, &payload); err != nil {
		return err
	}
	ctx := c.Request().Context()
	u := domain.User{
		Email:    payload.Email,
		Name:     payload.Name,
		Password: payload.Password,
	}
	if err := urh.uc.CreateUser(ctx, &u); err != nil {
		var code int
		if err == errors.ErrUserExists {
			code = http.StatusBadRequest
		} else {
			code = http.StatusInternalServerError
		}

		return wrapper.NewHTTPResponse(code, "failed to create user", err).JSON(c.Response())
	}

	return wrapper.NewHTTPResponse(http.StatusCreated, "Sucess", u).JSON(c.Response())
}

func getPayload(c echo.Context, payload interface{}) error {
	if err := c.Bind(&payload); err != nil {
		wrapper.NewHTTPResponse(http.StatusBadRequest, "bad payload", err).JSON(c.Response())
		return err
	}

	if _, err := govalidator.ValidateStruct(payload); err != nil {
		wrapper.NewHTTPResponse(http.StatusBadRequest, "bad payload", err).JSON(c.Response())
		return err
	}
	return nil
}
