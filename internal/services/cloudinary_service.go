package services

import (
	"context"
	"fmt"
	"mime/multipart"
	"os"
	"time"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

func UploadToCloudinary(file multipart.File, filename string) (string, error) {
	cld, err := cloudinary.NewFromParams(
		os.Getenv("CLOUDINARY_CLOUD_NAME"),
		os.Getenv("CLOUDINARY_API_KEY"),
		os.Getenv("CLOUDINARY_API_SECRET"),
	)
	if err != nil {
		return "", err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	uploadResult, err := cld.Upload.Upload(ctx, file, uploader.UploadParams{
		PublicID: filename,
		Folder:   "appointment-booking",
	})
	if err != nil {
		fmt.Println("Cloudinary Upload Error:", err)
		return "", err
	}

	fmt.Printf("Cloudinary Upload Result: %+v\n", uploadResult)
	return uploadResult.SecureURL, nil
}
