package model

import "mime/multipart"

type FileRepository interface {
	UploadAvatar(header *multipart.FileHeader, directory string) (string, string, error)
	UploadFile(header *multipart.FileHeader, directory, filename, mimetype string) (string, string, error)
	DeleteImage(fileId string) error
}
