package utils

import (
	"context"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/saurav11sarkar/go/internal/config"
)

type Cloudinary struct {
	cloudi *cloudinary.Cloudinary
}

func NewCloudinary(cfg config.Config) (*Cloudinary, error) {
	cloudi, err := cloudinary.NewFromParams(
		cfg.Cloudinary.CloudName,
		cfg.Cloudinary.ApiKey,
		cfg.Cloudinary.ApiSecret,
	)
	if err != nil {
		return nil, err
	}
	return &Cloudinary{cloudi: cloudi}, nil
}

func (c *Cloudinary) UploadFile(ctx context.Context, fileName any, folderName string) (string, error) {
	resp, err := c.cloudi.Upload.Upload(ctx, fileName, uploader.UploadParams{
		Folder:       folderName,
		ResourceType: "auto",
	})
	if err != nil {
		return "", err
	}
	return resp.SecureURL, nil
}
