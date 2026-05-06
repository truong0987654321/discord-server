package repository

import (
	"bytes"
	"context"
	"discord-server-go/model"
	"discord-server-go/model/apperrors"
	"discord-server-go/service"

	"image"
	"image/jpeg"
	"log"
	"mime/multipart"

	"github.com/disintegration/imaging"
	"github.com/imagekit-developer/imagekit-go/v2"
	"github.com/imagekit-developer/imagekit-go/v2/packages/param"
)

type imageKitRepository struct {
	Client *imagekit.Client
}

func NewFileRepository(client *imagekit.Client) model.FileRepository {
	return &imageKitRepository{
		Client: client,
	}
}
func (r *imageKitRepository) UploadAvatar(header *multipart.FileHeader, directory string) (string, string, error) {
	key := param.NewOpt("/discord-go-clone/" + directory)
	file, err := header.Open()
	if err != nil {
		log.Println("Open file error:", err)
		return "", "", apperrors.NewInternal()
	}
	defer file.Close()

	src, _, err := image.Decode(file)
	if err != nil {

		log.Println("Decode error:", err)
		return "", "", apperrors.NewInternal()
	}
	img := imaging.Resize(src, 150, 0, imaging.Lanczos)
	buf := new(bytes.Buffer)
	err = jpeg.Encode(buf, img, &jpeg.Options{Quality: 80})
	if err != nil {
		log.Println("encode error:", err)
		return "", "", apperrors.NewInternal()
	}
	ctx := context.Background()
	filename := service.GenerateId() + ".jpeg"

	resp, err := r.Client.Files.Upload(ctx, imagekit.FileUploadParams{
		File:     buf,
		FileName: filename,
		Folder:   key,
	})
	if err != nil {
		log.Println("Upload error:", err)
		return "", "", apperrors.NewInternal()
	}

	return resp.URL, resp.FileID, nil
}
func (r *imageKitRepository) UploadFile(header *multipart.FileHeader, directory, filename, mimetype string) (string, string, error) {

	key := param.NewOpt("/discord-go-clone/" + directory)
	file, err := header.Open()
	if err != nil {
		log.Printf("Failed to open header: %v\n", err)
		return "", "", apperrors.NewInternal()
	}
	defer file.Close()

	ctx := context.Background()

	resp, err := r.Client.Files.Upload(ctx, imagekit.FileUploadParams{
		File:     file,
		FileName: filename,
		Folder:   key,
	})
	if err != nil {
		return "", "", apperrors.NewInternal()
	}

	return resp.URL, resp.FileID, nil
}

func (r *imageKitRepository) DeleteImage(fileId string) error {
	ctx := context.Background()
	err := r.Client.Files.Delete(ctx, fileId)
	if err != nil {
		log.Println("Delete error:", err)
		return apperrors.NewInternal()
	}
	return nil
}
