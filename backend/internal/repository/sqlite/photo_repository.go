package sqlite

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/arthurshafikov/scout-takehome/backend/internal/core/models"
	"github.com/arthurshafikov/scout-takehome/backend/internal/core/types"
)

type SQLitePhotoRepository struct {
	db *sql.DB
}

func NewSQLitePhotoRepository(db *sql.DB) *SQLitePhotoRepository {
	return &SQLitePhotoRepository{
		db: db,
	}
}

func (r *SQLitePhotoRepository) GetPhotoByID(
	ctx context.Context,
	id string,
) (*models.Photo, error) {
	query := `
		SELECT id, x, y, h, width, height, captured_at
		FROM photos
		WHERE id = ?
	`

	row := r.db.QueryRowContext(ctx, query, id)

	photo := &models.Photo{}
	var capturedAt string

	if err := row.Scan(&photo.ID, &photo.X, &photo.Y, &photo.H, &photo.Width, &photo.Height, &capturedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("photo not found: %w", err)
		}

		return nil, fmt.Errorf("query photo: %w", err)
	}

	parsedTime, err := time.Parse(time.RFC3339, capturedAt)
	if err != nil {
		return nil, fmt.Errorf("parse captured_at: %w", err)
	}

	photo.CapturedAt = parsedTime

	predictions, err := r.getPredictionsForPhoto(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get predictions: %w", err)
	}

	photo.Predictions = predictions

	return photo, nil
}

func (r *SQLitePhotoRepository) ListPhotos(
	ctx context.Context,
	params types.ListPhotosParams,
) ([]models.Photo, string, error) {
	baseQuery := `
		SELECT DISTINCT p.id, p.x, p.y, p.h, p.width, p.height, p.captured_at
		FROM photos p
	`

	whereConditions := []string{}
	args := []interface{}{}

	if params.ClassID != "" || params.MinConfidence > 0 {
		baseQuery += ` JOIN predictions pred ON p.id = pred.photo_id`

		if params.ClassID != "" {
			whereConditions = append(whereConditions, `pred.class_id = ?`)
			args = append(args, params.ClassID)
		}

		if params.MinConfidence > 0 {
			whereConditions = append(whereConditions, `pred.confidence >= ?`)
			args = append(args, params.MinConfidence)
		}
	}

	if params.Cursor != "" {
		var cursorData struct {
			CapturedAt string
			ID         string
		}

		if err := json.Unmarshal([]byte(params.Cursor), &cursorData); err != nil {
			return nil, "", fmt.Errorf("decode cursor: %w", err)
		}

		whereConditions = append(whereConditions, `(p.captured_at, p.id) < (?, ?)`)
		args = append(args, cursorData.CapturedAt, cursorData.ID)
	}

	if len(whereConditions) > 0 {
		baseQuery += " WHERE " + whereConditions[0]
		for _, cond := range whereConditions[1:] {
			baseQuery += " AND " + cond
		}
	}

	baseQuery += ` ORDER BY p.captured_at DESC, p.id DESC`

	limit := params.Limit
	if limit <= 0 || limit > 200 {
		limit = 50
	}

	baseQuery += ` LIMIT ?`
	args = append(args, limit+1)

	rows, err := r.db.QueryContext(ctx, baseQuery, args...)
	if err != nil {
		return nil, "", fmt.Errorf("query photos: %w", err)
	}

	defer rows.Close()

	photos := []models.Photo{}
	var nextCursor string

	for rows.Next() {
		photo := models.Photo{}
		var capturedAt string

		if err := rows.Scan(
			&photo.ID,
			&photo.X,
			&photo.Y,
			&photo.H,
			&photo.Width,
			&photo.Height,
			&capturedAt,
		); err != nil {
			return nil, "", fmt.Errorf("scan photo: %w", err)
		}

		parsedTime, err := time.Parse(time.RFC3339, capturedAt)
		if err != nil {
			return nil, "", fmt.Errorf("parse captured_at: %w", err)
		}

		photo.CapturedAt = parsedTime

		predictions, err := r.getPredictionsForPhoto(ctx, photo.ID)
		if err != nil {
			return nil, "", fmt.Errorf("get predictions for photo %s: %w", photo.ID, err)
		}

		photo.Predictions = predictions
		photos = append(photos, photo)
	}

	if err := rows.Err(); err != nil {
		return nil, "", fmt.Errorf("rows error: %w", err)
	}

	if len(photos) > limit {
		photos = photos[:limit]
		lastPhoto := photos[limit-1]

		cursorData := map[string]string{
			"captured_at": lastPhoto.CapturedAt.Format(time.RFC3339),
			"id":          lastPhoto.ID,
		}

		cursorBytes, err := json.Marshal(cursorData)
		if err != nil {
			return nil, "", fmt.Errorf("encode cursor: %w", err)
		}

		nextCursor = string(cursorBytes)
	}

	return photos, nextCursor, nil
}

func (r *SQLitePhotoRepository) getPredictionsForPhoto(
	ctx context.Context,
	photoID string,
) ([]models.Prediction, error) {
	query := `
		SELECT id, photo_id, class_id, confidence, bbox_xmin, bbox_ymin, bbox_xmax, bbox_ymax
		FROM predictions
		WHERE photo_id = ?
		ORDER BY confidence DESC
	`

	rows, err := r.db.QueryContext(ctx, query, photoID)
	if err != nil {
		return nil, fmt.Errorf("query predictions: %w", err)
	}

	defer rows.Close()

	predictions := []models.Prediction{}

	for rows.Next() {
		pred := models.Prediction{}

		if err := rows.Scan(
			&pred.ID,
			&pred.PhotoID,
			&pred.ClassID,
			&pred.Confidence,
			&pred.BBox.XMin,
			&pred.BBox.YMin,
			&pred.BBox.XMax,
			&pred.BBox.YMax,
		); err != nil {
			return nil, fmt.Errorf("scan prediction: %w", err)
		}

		predictions = append(predictions, pred)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return predictions, nil
}
