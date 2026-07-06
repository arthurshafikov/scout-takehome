package types

// ListPhotosParams defines cursor-based pagination and filtering parameters
type ListPhotosParams struct {
	Cursor        string
	Limit         int
	ClassID       string
	MinConfidence float64
}
