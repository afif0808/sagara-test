package product

import (
	authusecase "github.com/afif0808/sagara-test/internal/modules/auth/usecase"
	"github.com/afif0808/sagara-test/internal/modules/product/handler/rest"
	productrepo "github.com/afif0808/sagara-test/internal/modules/product/repository"
	userrepo "github.com/afif0808/sagara-test/internal/modules/user/repository"

	productusecase "github.com/afif0808/sagara-test/internal/modules/product/usecase"

	"github.com/afif0808/sagara-test/middleware"

	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo"
)

func InjectProductModule(e *echo.Echo, readDB, writeDB *sqlx.DB) {
	// the purpose of doing it this way is to make repo and usecase variable
	// could consist of multiple repository

	var repo struct {
		productrepo.ProductSQLRepository
		userrepo.UserSQLRepository
	}

	repo.ProductSQLRepository = productrepo.NewProductSQLRepository(readDB, writeDB)
	repo.UserSQLRepository = userrepo.NewUserSQLRepository(readDB, writeDB)

	var usecase struct {
		productusecase.ProductUsecase
		authusecase.AuthUsecase
	}
	usecase.ProductUsecase = productusecase.NewProductUsecase(&repo)
	usecase.AuthUsecase = authusecase.NewAuthUsecase(&repo.UserSQLRepository)

	mw := middleware.NewHTTPMiddleware(&usecase)

	restHandler := rest.NewProductRestHandler(&usecase, &mw)
	restHandler.Mount(e.Group(""))
}
