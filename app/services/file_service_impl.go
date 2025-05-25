package services

import (
	"io"
	"mime/multipart"
	"os"
)

type FileServiceImpl struct{}

func NewFileService() *FileServiceImpl {
	return &FileServiceImpl{}
}

func (c *FileServiceImpl) SaveTempFile(file *multipart.File) (string, error) {
	f, err := os.CreateTemp("", "sample")
	if err != nil {
		return "", err
	}

	_, err = io.Copy(f, *file)
	if err != nil {
		return "", err
	}

	return f.Name(), nil
}
