package storage

import (
	"context"
	"mime/multipart"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/osamikoyo/dark-fantasy-land/pkg/logger"
	"go.uber.org/zap"
)

type Storage struct {
	logger  *logger.Logger
	client  *minio.Client
	timeout time.Duration
}

func NewStorage(client *minio.Client, logger *logger.Logger, timeout time.Duration) *Storage {
	return &Storage{
		logger:  logger,
		client:  client,
		timeout: timeout,
	}
}

func (s *Storage) context() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), s.timeout)
}

func (s *Storage) UploadFile(file *multipart.FileHeader, bucketName string) error {
	ctx, cancel := s.context()
	defer cancel()

	src, err := file.Open()
	if err != nil {
		s.logger.Error("failed upload file", zap.Error(err))

		return err
	}
	defer src.Close()

	_, err = s.client.PutObject(
		ctx,
		bucketName,
		file.Filename,
		src,
		file.Size,
		minio.PutObjectOptions{ContentType: "application/octet-stream"},
	)
	if err != nil {
		s.logger.Error("failed put object", zap.Error(err))
	}

	return nil
}

func (s *Storage) DownloadFile(filename, bucketName string) (*minio.Object, error) {
	ctx, cancel := s.context()
	defer cancel()

	obj, err := s.client.GetObject(
		ctx,
		bucketName,
		filename,
		minio.GetObjectOptions{},
	)
	if err != nil {
		s.logger.Error("failed download file",
			zap.String("filename", filename),
			zap.String("bucket_name", bucketName),
			zap.Error(err))

		return nil, err
	}

	return obj, nil
}
