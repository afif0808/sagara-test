package rest

import (
	"context"
	"net/http"

	"github.com/afif0808/sagara-test/internal/domain"
	"github.com/afif0808/sagara-test/wrapper"
	"github.com/asaskevich/govalidator"
	"github.com/labstack/echo"
)

type usecase interface {
	Login(ctx context.Context, crd domain.LoginCredentials) (domain.User, string, error)
}

type AuthRestHandler struct {
	uc usecase
}

func NewAuthRestHandler(uc usecase) AuthRestHandler {
	return AuthRestHandler{
		uc: uc,
	}
}

func (arh AuthRestHandler) Mount(root *echo.Group) {
	auth := root.Group("/auth/")
	auth.POST("login", arh.login)
}

func (arh AuthRestHandler) login(c echo.Context) error {
	ctx := c.Request().Context()
	var payload struct {
		Email    string `json:"email" valid:"required"`
		Password string `json:"password" valid:"required"`
	}

	if err := getPayload(c, &payload); err != nil {
		return err
	}
	var response struct {
		User  domain.User `json:"user"`
		Token string      `json:"token"`
	}

	var err error
	response.User, response.Token, err = arh.uc.Login(ctx, domain.LoginCredentials{
		Identity: payload.Email,
		Password: payload.Password,
	})

	if err != nil {
		return wrapper.NewHTTPResponse(http.StatusUnauthorized, "login failed", err).JSON(c.Response())
	}

	return wrapper.NewHTTPResponse(http.StatusOK, "Success", response).JSON(c.Response())
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
