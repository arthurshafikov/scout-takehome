package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func main() {
	endpoint := flag.String("endpoint", "localhost:9000", "MinIO endpoint")
	accessKey := flag.String("access-key", "minioadmin", "MinIO access key")
	secretKey := flag.String("secret-key", "minioadmin", "MinIO secret key")
	bucket := flag.String("bucket", "scout", "MinIO bucket name")
	useSSL := flag.Bool("use-ssl", false, "Use SSL for MinIO connection")
	datasetDir := flag.String("dataset", "../dataset/images", "Dataset directory with images")

	flag.Parse()

	// Validate dataset directory exists
	if _, err := os.Stat(*datasetDir); err != nil {
		fmt.Fprintf(os.Stderr, "Error: dataset directory not found: %v\n", err)
		os.Exit(1)
	}

	// Create MinIO client
	client, err := minio.New(*endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(*accessKey, *secretKey, ""),
		Secure: *useSSL,
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: failed to create MinIO client: %v\n", err)
		os.Exit(1)
	}

	ctx := context.Background()

	// Check if bucket exists
	exists, err := client.BucketExists(ctx, *bucket)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: failed to check bucket existence: %v\n", err)
		os.Exit(1)
	}

	// Create bucket if it doesn't exist
	if !exists {
		fmt.Printf("Creating bucket '%s'...\n", *bucket)
		if err := client.MakeBucket(ctx, *bucket, minio.MakeBucketOptions{}); err != nil {
			fmt.Fprintf(os.Stderr, "Error: failed to create bucket: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Bucket '%s' created successfully\n", *bucket)
	}

	// Upload all images from dataset directory
	count := 0
	err = filepath.Walk(*datasetDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip directories
		if info.IsDir() {
			return nil
		}

		// Check if file is an image (jpg, jpeg, png)
		ext := filepath.Ext(info.Name())
		if ext != ".jpg" && ext != ".jpeg" && ext != ".png" {
			return nil
		}

		// Get relative path for object name
		relPath, err := filepath.Rel(*datasetDir, path)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Warning: failed to get relative path for %s: %v\n", path, err)
			return nil
		}

		objectName := "originals/" + relPath

		fmt.Printf("Uploading %s to %s/%s...\n", relPath, *bucket, objectName)

		// Upload file
		file, err := os.Open(path)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Warning: failed to open file %s: %v\n", path, err)
			return nil
		}
		defer file.Close()

		fileInfo, err := file.Stat()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Warning: failed to stat file %s: %v\n", path, err)
			return nil
		}

		_, err = client.PutObject(ctx, *bucket, objectName, file, fileInfo.Size(), minio.PutObjectOptions{
			ContentType: getContentType(ext),
		})
		if err != nil {
			fmt.Fprintf(os.Stderr, "Warning: failed to upload %s: %v\n", path, err)
			return nil
		}

		count++
		fmt.Printf("✓ Uploaded %d files\n", count)

		return nil
	})

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: failed to walk dataset directory: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("\nSuccessfully uploaded %d images to MinIO bucket '%s'\n", count, *bucket)
}

func getContentType(ext string) string {
	switch ext {
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".png":
		return "image/png"
	default:
		return "application/octet-stream"
	}
}
