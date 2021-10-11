package rest

import (
	"context"
	"mime/multipart"
	"net/http"

	"github.com/afif0808/sagara-test/wrapper"
	"github.com/labstack/echo"
)

type usecase interface {
	UploadFile(ctx context.Context, f multipart.FileHeader) (fileName string, err error)
}

type middleware interface {
	BearerAuth(next http.Handler) http.Handler
}

type FileStorageRestHandler struct {
	uc usecase
	mw middleware
}

func NewFileStorageRestHandler(uc usecase, mw middleware) FileStorageRestHandler {
	return FileStorageRestHandler{uc: uc, mw: mw}
}

func (fsrh FileStorageRestHandler) Mount(root *echo.Group) {
	g := root.Group("/file/storage/")
	g.Static("files", "internal/modules/filestorage/storage")
	g.POST("upload", fsrh.uploadFileStorage, echo.WrapMiddleware(fsrh.mw.BearerAuth))
}

func (fsrh FileStorageRestHandler) uploadFileStorage(c echo.Context) (err error) {
	ff, err := c.FormFile("file")
	if err != nil {
		return wrapper.NewHTTPResponse(http.StatusBadRequest, "failed to upload file").JSON(c.Response())
	}
	ctx := c.Request().Context()
	fileName, err := fsrh.uc.UploadFile(ctx, *ff)
	if err != nil {
		return wrapper.NewHTTPResponse(http.StatusInternalServerError, "failed to upload file").JSON(c.Response())
	}
	response := struct {
		FileURL string `json:"file_url"`
	}{
		FileURL: c.Request().Host + "/file/storage/files/" + fileName,
	}

	return wrapper.NewHTTPResponse(http.StatusOK, "Success", response).JSON(c.Response())

}
