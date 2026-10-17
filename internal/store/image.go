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

func (i *ImageStore) GetImage(ctx context.Context, objectName string) ([]byte, error) {
	obj, err := i.minioClient.GetObject(ctx, i.bucketName, objectName, minio.GetObjectOptions{})
	if err != nil {
		return nil, err
	}
	defer obj.Close()

	return io.ReadAll(obj)
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
