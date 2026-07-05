package storage

import (
	"context"
	"fmt"
	"time"

	"github.com/arthurshafikov/scout-takehome/backend/internal/config"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type UploadLink struct {
	URL       string            `json:"url"`
	Method    string            `json:"method"`
	Headers   map[string]string `json:"headers,omitempty"`
	ExpiresAt time.Time         `json:"expiresAt"`
}

type StorageService interface {
	GenerateUploadLink(ctx context.Context, photoID string, contentType string) (*UploadLink, error)
	GetOriginalURL(ctx context.Context, photoID string) (string, error)
	ObjectExists(ctx context.Context, photoID string) (bool, error)
	GetMinIOClient() (*minio.Client, error)
}

type MinIOStorageService struct {
	client *minio.Client
	bucket string
}

func NewMinIOStorageService(config *config.MinIOConfig) (StorageService, error) {
	client, err := minio.New(config.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(config.AccessKey, config.SecretKey, ""),
		Secure: config.UseSSL,
	})
	if err != nil {
		return nil, fmt.Errorf("create minio client: %w", err)
	}

	return &MinIOStorageService{
		client: client,
		bucket: config.Bucket,
	}, nil
}

func (s *MinIOStorageService) GenerateUploadLink(
	ctx context.Context,
	photoID string,
	contentType string,
) (*UploadLink, error) {
	objectName := fmt.Sprintf("originals/%s", photoID)
	expiresAt := time.Now().Add(15 * time.Minute)
	duration := expiresAt.Sub(time.Now())

	presignedURL, err := s.client.PresignedPutObject(ctx, s.bucket, objectName, duration)
	if err != nil {
		return nil, fmt.Errorf("generate presigned put url: %w", err)
	}

	return &UploadLink{
		URL:       presignedURL.String(),
		Method:    "PUT",
		ExpiresAt: expiresAt,
	}, nil
}

func (s *MinIOStorageService) GetOriginalURL(ctx context.Context, photoID string) (string, error) {
	objectName := fmt.Sprintf("originals/%s", photoID)
	expiresAt := time.Now().Add(1 * time.Hour)
	duration := expiresAt.Sub(time.Now())

	presignedURL, err := s.client.PresignedGetObject(ctx, s.bucket, objectName, duration, nil)
	if err != nil {
		return "", fmt.Errorf("generate presigned get url: %w", err)
	}

	return presignedURL.String(), nil
}

func (s *MinIOStorageService) ObjectExists(ctx context.Context, photoID string) (bool, error) {
	objectName := fmt.Sprintf("originals/%s", photoID)

	_, err := s.client.StatObject(ctx, s.bucket, objectName, minio.StatObjectOptions{})
	if err != nil {
		errResponse := minio.ToErrorResponse(err)
		if errResponse.Code == "NoSuchKey" {
			return false, nil
		}

		return false, fmt.Errorf("stat object: %w", err)
	}

	return true, nil
}

func (s *MinIOStorageService) GetMinIOClient() (*minio.Client, error) {
	if s.client == nil {
		return nil, fmt.Errorf("minio client not initialized")
	}
	return s.client, nil
}
