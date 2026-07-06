package thumbnail

import (
	"bytes"
	"context"
	"crypto/md5"
	"fmt"
	"image"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/arthurshafikov/scout-takehome/backend/internal/services/metrics"
	"github.com/disintegration/imaging"
	"github.com/minio/minio-go/v7"
	"github.com/sirupsen/logrus"
)

type ThumbnailGenerator struct {
	client      *minio.Client
	bucket      string
	datasetPath string
	logger      *logrus.Logger
	metrics     *metrics.Metrics
}

type ThumbnailOptions struct {
	Width   int
	Height  int
	Quality int // 1-100
}

// DefaultOptions returns default thumbnail options
func DefaultOptions() ThumbnailOptions {
	return ThumbnailOptions{
		Width:   400,
		Height:  300,
		Quality: 85,
	}
}

func NewThumbnailGenerator(client *minio.Client, bucket, datasetPath string, m *metrics.Metrics) *ThumbnailGenerator {
	logger := logrus.New()
	logger.SetLevel(logrus.DebugLevel)
	return &ThumbnailGenerator{
		client:      client,
		bucket:      bucket,
		datasetPath: datasetPath,
		logger:      logger,
		metrics:     m,
	}
}

// Generate creates a thumbnail from an original photo
func (tg *ThumbnailGenerator) Generate(
	ctx context.Context,
	photoID string,
	opts ThumbnailOptions,
) ([]byte, error) {
	// Try to get from cache first
	cacheKey := tg.cacheKey(photoID, opts)
	cached, err := tg.getFromCache(ctx, cacheKey)
	if err == nil {
		tg.logger.Debugf("Cache HIT: %s", cacheKey)
		if tg.metrics != nil {
			tg.metrics.RecordThumbnailCacheHit()
		}
		return cached, nil
	}
	tg.logger.Debugf("Cache MISS: %s (%v)", cacheKey, err)
	if tg.metrics != nil {
		tg.metrics.RecordThumbnailCacheMiss()
	}

	start := time.Now()

	// Download original
	original, err := tg.downloadOriginal(ctx, photoID)
	if err != nil {
		return nil, fmt.Errorf("download original: %w", err)
	}

	// Decode image
	img, _, err := image.Decode(bytes.NewReader(original))
	if err != nil {
		return nil, fmt.Errorf("decode image: %w", err)
	}

	// Resize
	resized := imaging.Resize(img, opts.Width, opts.Height, imaging.Lanczos)

	// Encode JPEG
	thumbnail := &bytes.Buffer{}
	err = imaging.Encode(thumbnail, resized, imaging.JPEG, imaging.JPEGQuality(opts.Quality))
	if err != nil {
		return nil, fmt.Errorf("encode thumbnail: %w", err)
	}

	duration := time.Since(start).Seconds()
	if tg.metrics != nil {
		tg.metrics.RecordThumbnailGenerationTime(duration)
	}

	// Cache it (best effort, log errors but don't fail)
	if err := tg.cacheToMinIO(ctx, cacheKey, thumbnail.Bytes()); err != nil {
		tg.logger.Warnf("Failed to cache thumbnail to MinIO: %v", err)
	} else {
		tg.logger.Debugf("Successfully cached thumbnail: %s", cacheKey)
	}

	return thumbnail.Bytes(), nil
}

func (tg *ThumbnailGenerator) cacheKey(photoID string, opts ThumbnailOptions) string {
	hash := md5.Sum([]byte(fmt.Sprintf("%s_%d_%d_%d", photoID, opts.Width, opts.Height, opts.Quality)))
	return fmt.Sprintf("thumbs/%s_%x", photoID, hash)
}

func (tg *ThumbnailGenerator) downloadOriginal(ctx context.Context, photoID string) ([]byte, error) {
	// Try to read from local filesystem first (dataset/images/{photoID}.jpg)
	if tg.datasetPath != "" {
		localPath := filepath.Join(tg.datasetPath, "images", photoID+".jpg")
		data, err := os.ReadFile(localPath)
		if err == nil {
			return data, nil
		}
		// If local file not found, continue to try MinIO
	}

	// Fallback to MinIO for cloud-based originals
	objectName := fmt.Sprintf("originals/%s", photoID)

	obj, err := tg.client.GetObject(ctx, tg.bucket, objectName, minio.GetObjectOptions{})
	if err != nil {
		return nil, fmt.Errorf("get object: %w", err)
	}
	defer func() {
		_ = obj.Close() // ignore error
	}()

	buf := &bytes.Buffer{}
	_, err = buf.ReadFrom(obj)
	if err != nil {
		return nil, fmt.Errorf("read object: %w", err)
	}

	return buf.Bytes(), nil
}

func (tg *ThumbnailGenerator) getFromCache(ctx context.Context, cacheKey string) ([]byte, error) {
	obj, err := tg.client.GetObject(ctx, tg.bucket, cacheKey, minio.GetObjectOptions{})
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = obj.Close() // ignore close errors
	}()

	data, err := io.ReadAll(obj)
	if err != nil {
		errResp := minio.ToErrorResponse(err)
		if errResp.Code == "NoSuchKey" {
			return nil, fmt.Errorf("not found")
		}
		return nil, err
	}

	return data, nil
}

func (tg *ThumbnailGenerator) cacheToMinIO(ctx context.Context, cacheKey string, data []byte) error {
	_, err := tg.client.PutObject(ctx, tg.bucket, cacheKey, bytes.NewReader(data), int64(len(data)), minio.PutObjectOptions{
		ContentType: "image/jpeg",
	})
	return err
}

// GetPresignedURL returns a presigned URL for thumbnail
func (tg *ThumbnailGenerator) GetPresignedURL(ctx context.Context, cacheKey string) (string, error) {
	expiresAt := time.Now().Add(1 * time.Hour)
	duration := time.Until(expiresAt)

	presignedURL, err := tg.client.PresignedGetObject(ctx, tg.bucket, cacheKey, duration, nil)
	if err != nil {
		return "", fmt.Errorf("generate presigned url: %w", err)
	}

	return presignedURL.String(), nil
}
