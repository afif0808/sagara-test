package usecase

import (
	"context"
	"io"
	"mime/multipart"
	"os"
	"path"
	"strings"

	"github.com/google/uuid"
)

type FileStorageUsecase struct {
}

func NewFileStorageUsecase() FileStorageUsecase {
	return FileStorageUsecase{}
}

func (fsu *FileStorageUsecase) UploadFile(ctx context.Context, fh multipart.FileHeader) (fileName string, err error) {
	fileName = uuid.NewString() + path.Ext(fh.Filename)
	src, err := fh.Open()
	if err != nil {
		return
	}
	storagePath := os.Getenv("FILE_STORAGE")
	storagePath = strings.TrimSuffix(storagePath, "/")
	dst, err := os.Create(storagePath + "/" + fileName)
	if err != nil {
		return
	}
	_, err = io.Copy(dst, src)
	return
}
