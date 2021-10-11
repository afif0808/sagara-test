package user

import (
	resthandler "github.com/afif0808/sagara-test/internal/modules/auth/handler/rest"
	authusecase "github.com/afif0808/sagara-test/internal/modules/auth/usecase"
	userrepo "github.com/afif0808/sagara-test/internal/modules/user/repository"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo"
)

func InjectAuthModule(e *echo.Echo, readDB, writeDB *sqlx.DB) {
	// the purpose of doing it this way is to make repo and usecase variable
	// could consist of multiple repository

	var repo struct {
		userrepo.UserSQLRepository
	}
	repo.UserSQLRepository = userrepo.NewUserSQLRepository(readDB, writeDB)

	var usecase struct {
		authusecase.AuthUsecase
	}
	usecase.AuthUsecase = authusecase.NewAuthUsecase(&repo)

	restHandler := resthandler.NewAuthRestHandler(&usecase)
	restHandler.Mount(e.Group(""))
}
