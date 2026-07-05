package types

type Pagination struct {
	Page    int
	PerPage int
}

func (p *Pagination) GetLimit() int {
	return p.PerPage
}

func (p *Pagination) GetOffset() int {
	if p.Page < 1 {
		return 0
	}

	return (p.Page - 1) * p.PerPage
}

type ListPhotosParams struct {
	Cursor        string
	Limit         int
	ClassID       string
	MinConfidence float64
}
