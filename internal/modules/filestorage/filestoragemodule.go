package filestorage

import (
	authusecase "github.com/afif0808/sagara-test/internal/modules/auth/usecase"
	"github.com/afif0808/sagara-test/internal/modules/filestorage/handler/rest"
	filestorageusecase "github.com/afif0808/sagara-test/internal/modules/filestorage/usecase"
	productrepo "github.com/afif0808/sagara-test/internal/modules/product/repository"
	"github.com/afif0808/sagara-test/middleware"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo"
)

func InjectFileStorageModule(e *echo.Echo, readDB, writeDB *sqlx.DB) {
	var repo struct {
		productrepo.ProductSQLRepository
	}

	repo.ProductSQLRepository = productrepo.NewProductSQLRepository(readDB, writeDB)

	var usecase struct {
		filestorageusecase.FileStorageUsecase
		authusecase.AuthUsecase
	}
	usecase.FileStorageUsecase = filestorageusecase.NewFileStorageUsecase()
	middleware := middleware.NewHTTPMiddleware(&usecase)
	restHandler := rest.NewFileStorageRestHandler(&usecase, &middleware)
	restHandler.Mount(e.Group(""))
}
