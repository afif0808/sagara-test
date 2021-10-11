package meta

import "math"

type Meta struct {
	Page         int `json:"page"`
	Limit        int `json:"limit"`
	TotalRecords int `json:"total_records"`
	TotalPages   int `json:"total_pages"`
}

func NewMeta(page, limit, totalRecords int) (m Meta) {
	m.Page, m.Limit, m.TotalRecords = page, limit, totalRecords
	m.CalculatePages()
	return m
}

func (m *Meta) CalculatePages() {
	m.TotalPages = int(math.Ceil(float64(m.TotalRecords) / float64(m.Limit)))
}
