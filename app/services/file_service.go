package services

import "mime/multipart"

type FileService interface {
	SaveTempFile(file *multipart.File) (string, error)
}
