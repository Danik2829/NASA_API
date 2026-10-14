package core

import (
	"context"
	"io"

	"github.com/minio/minio-go/v7"
)

type MetaDataActions interface {
	SaveApod(apod *Apod) error
	GetApodByDate(date string) (*Apod, error)
	GetAllApods() ([]Apod, error)
}

type ImageActions interface {
	GetImage(ctx context.Context, objectName, filePath, contentType string) (minio.UploadInfo, error)
	PutImage(
		ctx context.Context,
		objectName string,
		reader io.Reader,
		size int64,
		contentType string,
	) error
}

type NasaAPI interface {
	GetApodByDate(date string) (*Apod, error)
}

type Apod struct {
	ID          int    `json:"post_id"`
	Date        string `json:"date"`
	Explanation string `json:"explanation"`
	HdURL       string `json:"hdurl"`
	MediaType   string `json:"media_type"`
	Title       string `json:"title"`
	Url         string `json:"url"`
}
