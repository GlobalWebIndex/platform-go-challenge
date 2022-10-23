package pagination

import "strconv"

// Pagination struct
type Pagination struct {
	Page    int
	PerPage int
}

// Offset returns the offset based on page and PerPage.
func (p *Pagination) Offset() int {
	return (p.Page - 1) * p.PerPage
}

// NewPagination constructor
func NewPagination(p map[string][]string) *Pagination {
	pn := Pagination{Page: 1, PerPage: 10}

	value, found := p["page"]
	if found {
		page, err := strconv.Atoi(value[0])
		if err == nil {
			pn.Page = page
		}
	}

	value, found = p["per_page"]
	if found {
		perPage, err := strconv.Atoi(value[0])
		if err == nil {
			pn.PerPage = perPage
		}
	}

	return &pn
}
