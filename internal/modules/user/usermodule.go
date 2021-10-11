package user

import (
	"github.com/afif0808/sagara-test/internal/modules/user/handler/rest"
	userrepo "github.com/afif0808/sagara-test/internal/modules/user/repository"
	userusecase "github.com/afif0808/sagara-test/internal/modules/user/usecase"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo"
)

func InjectUserModule(e *echo.Echo, readDB, writeDB *sqlx.DB) {
	// the purpose of doing it this way is to make repo and usecase variable
	// could consist of multiple repository

	var repo struct {
		userrepo.UserSQLRepository
	}
	repo.UserSQLRepository = userrepo.NewUserSQLRepository(readDB, writeDB)

	var usecase struct {
		userusecase.UserUsecase
	}
	usecase.UserUsecase = userusecase.NewUserUsecase(&repo)

	restHandler := rest.NewUserRestHandler(&usecase)
	restHandler.Mount(e.Group(""))
}
