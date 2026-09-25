package store

import (
	"context"

	"github.com/NASA_API/internal/config"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type ImageStore struct {
	minioClient *minio.Client
	bucketName  string
}

func NewMinioStore(cfg config.MinioConfig) (*minio.Client, error) {
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

	return minioClient, nil
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

func (i *ImageStore) PutImage(ctx context.Context, objectName, filePath string) (minio.UploadInfo, error) {
	info, err := i.minioClient.FPutObject(
		ctx,
		i.bucketName,
		objectName,
		filePath,
		minio.PutObjectOptions{},
	)
	return info, err
}
