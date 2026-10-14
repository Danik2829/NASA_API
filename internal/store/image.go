package store

import (
	"context"
	"fmt"
	"io"

	"github.com/NASA_API/internal/config"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type ImageStore struct {
	minioClient *minio.Client
	bucketName  string
}

func NewMinioStore(cfg config.MinioConfig) (*ImageStore, error) {
	minioClient, err := minio.New(cfg.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.AccessKey, cfg.SecretKey, ""),
		Secure: false,
	})
	if err != nil {
		return nil, err
	}

	ctx := context.Background()

	err = minioClient.MakeBucket(ctx, cfg.Bucket, minio.MakeBucketOptions{Region: "us-east-1"})
	if err != nil {
		exists, errBucketExists := minioClient.BucketExists(ctx, cfg.Bucket)
		if errBucketExists != nil && !exists {
			return nil, err
		}
	}

	return &ImageStore{minioClient: minioClient, bucketName: cfg.Bucket}, nil
}

func (i *ImageStore) GetImage(ctx context.Context, objectName, filePath, contentType string) (minio.UploadInfo, error) {
	info, err := i.minioClient.FPutObject(
		ctx,
		i.bucketName,
		objectName,
		filePath,
		minio.PutObjectOptions{ContentType: contentType},
	)
	return info, err
}

func (m *ImageStore) PutImage(
	ctx context.Context,
	objectName string,
	reader io.Reader,
	size int64,
	contentType string,
) error {
	_, err := m.minioClient.PutObject(
		ctx,
		m.bucketName,
		objectName,
		reader,
		size,
		minio.PutObjectOptions{
			ContentType: contentType,
		},
	)

	if err != nil {
		return fmt.Errorf("upload image to MinIO: %w", err)
	}

	return nil
}
