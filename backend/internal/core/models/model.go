package models

import "time"

// Photo represents a greenhouse photo with metadata and predictions.
type Photo struct {
	ID          string       `json:"id"`
	X           float64      `json:"x"`
	Y           float64      `json:"y"`
	H           float64      `json:"h"`
	Width       int          `json:"width"`
	Height      int          `json:"height"`
	CapturedAt  time.Time    `json:"capturedAt"`
	OriginalURL string       `json:"originalUrl,omitempty"`
	Predictions []Prediction `json:"predictions"`
}

// Prediction represents a pest/disease detection in a photo.
type Prediction struct {
	ID         string      `json:"id"`
	PhotoID    string      `json:"photoId"`
	ClassID    string      `json:"classId"`
	Confidence float64     `json:"confidence"`
	BBox       BoundingBox `json:"bbox"`
}

// BoundingBox represents a normalized bounding box [0,1].
type BoundingBox struct {
	XMin float64 `json:"xMin"`
	YMin float64 `json:"yMin"`
	XMax float64 `json:"xMax"`
	YMax float64 `json:"yMax"`
}
