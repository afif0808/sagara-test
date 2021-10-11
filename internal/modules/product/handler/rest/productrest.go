package rest

import (
	"context"
	"log"
	"net/http"
	"strconv"

	"github.com/afif0808/sagara-test/dataquery"
	"github.com/afif0808/sagara-test/internal/domain"
	"github.com/afif0808/sagara-test/meta"
	"github.com/afif0808/sagara-test/wrapper"
	"github.com/asaskevich/govalidator"
	"github.com/labstack/echo"
)

type usecase interface {
	GetProductList(ctx context.Context, dq dataquery.DataQuery) ([]domain.Product, meta.Meta, error)
	GetProduct(ctx context.Context, id int64) (domain.Product, error)
	CreateProduct(ctx context.Context, p *domain.Product) error
	UpdateProduct(ctx context.Context, p *domain.Product) error
	DeleteProduct(ctx context.Context, id int64) error
}

type middleware interface {
	BearerAuth(next http.Handler) http.Handler
}

type ProductRestHandler struct {
	uc usecase
	mw middleware
}

func NewProductRestHandler(uc usecase, mw middleware) ProductRestHandler {
	return ProductRestHandler{
		uc: uc,
		mw: mw,
	}
}

func (prh ProductRestHandler) Mount(root *echo.Group) {
	product := root.Group("/product/")
	product.GET("", prh.getProductList, echo.WrapMiddleware(prh.mw.BearerAuth))
	product.POST("", prh.createProduct, echo.WrapMiddleware(prh.mw.BearerAuth))
	product.PUT(":id", prh.updateProduct, echo.WrapMiddleware(prh.mw.BearerAuth))
	product.DELETE(":id", prh.deleteProduct, echo.WrapMiddleware(prh.mw.BearerAuth))
	product.GET(":id", prh.getProduct, echo.WrapMiddleware(prh.mw.BearerAuth))

}

func (prh ProductRestHandler) getProductList(c echo.Context) error {
	dq, err := dataquery.ParseFromURLQuery(c.QueryParams())
	if err != nil {
		return wrapper.NewHTTPResponse(http.StatusBadRequest, "bad payload", err).JSON(c.Response())
	}
	ctx := c.Request().Context()
	data, meta, err := prh.uc.GetProductList(ctx, dq)
	if err != nil {
		return wrapper.NewHTTPResponse(http.StatusInternalServerError, "failed to get product list", err).JSON(c.Response())
	}
	return wrapper.NewHTTPResponse(http.StatusOK, "Success", data, meta).JSON(c.Response())

}

func (prh ProductRestHandler) getProduct(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return wrapper.NewHTTPResponse(http.StatusNoContent, "product with given id is not found", nil).JSON(c.Response())
	}
	ctx := c.Request().Context()
	product, err := prh.uc.GetProduct(ctx, id)
	if err != nil {
		return wrapper.NewHTTPResponse(http.StatusInternalServerError, "failed to get product", err).JSON(c.Response())
	}
	return wrapper.NewHTTPResponse(http.StatusCreated, "Success", product).JSON(c.Response())
}

func (prh ProductRestHandler) createProduct(c echo.Context) error {
	ctx := c.Request().Context()
	var payload struct {
		Name     string `json:"name" valid:"required"`
		ImageURL string `json:"image_url"`
	}
	if err := getPayload(c, &payload); err != nil {
		return err
	}
	log.Println("lanjut")
	product := domain.Product{
		Name:     payload.Name,
		ImageURL: payload.ImageURL,
	}

	if err := prh.uc.CreateProduct(ctx, &product); err != nil {
		return wrapper.NewHTTPResponse(http.StatusInternalServerError, "failed to create product", err).JSON(c.Response())
	}

	return wrapper.NewHTTPResponse(http.StatusCreated, "Success", product).JSON(c.Response())
}

func (prh ProductRestHandler) updateProduct(c echo.Context) error {
	ctx := c.Request().Context()
	var payload struct {
		Name     string `json:"name" valid:"required"`
		ImageURL string `json:"image_url"`
	}
	if err := getPayload(c, &payload); err != nil {
		return err
	}

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return wrapper.NewHTTPResponse(http.StatusNoContent, "product with given id is not found", nil).JSON(c.Response())
	}

	product := domain.Product{
		Name:     payload.Name,
		ImageURL: payload.ImageURL,
		ID:       id,
	}

	if err := prh.uc.UpdateProduct(ctx, &product); err != nil {
		return wrapper.NewHTTPResponse(http.StatusInternalServerError, "failed to update product", err).JSON(c.Response())
	}

	return wrapper.NewHTTPResponse(http.StatusOK, "Success", product).JSON(c.Response())
}

func (prh ProductRestHandler) deleteProduct(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return wrapper.NewHTTPResponse(http.StatusNoContent, "product with given id is not found", nil).JSON(c.Response())
	}
	ctx := c.Request().Context()
	err = prh.uc.DeleteProduct(ctx, id)
	if err != nil {
		return wrapper.NewHTTPResponse(http.StatusInternalServerError, "failed to delete product", err).JSON(c.Response())
	}
	return wrapper.NewHTTPResponse(http.StatusOK, "Success").JSON(c.Response())
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
