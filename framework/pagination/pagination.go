package pagination

import (
	"net/http"
	"strconv"

	"gorm.io/gorm"
)

const (
	DEFAULT_PAGE_SIZE = 50
	MAX_PAGE_SIZE     = 100

	GET_OPTION_KEY_PAGE      = "page"
	GET_OPTION_KEY_PAGE_SIZE = "page_size"
)

// Pagination Exports
type Pagination struct {
	Page     int   `json:"page"`
	NumPages int   `json:"numPages"`
	PageSize int   `json:"pageSize"`
	Count    int64 `json:"count"`
}

// Paginate Exports
func (p *Pagination) Paginate() func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		offset := (p.Page - 1) * p.PageSize
		return db.Offset(offset).Limit(p.PageSize)
	}
}

// NewPagination Exports
func NewPagination(page, pageSize int, model *gorm.DB) *Pagination {
	switch {
	case pageSize > MAX_PAGE_SIZE:
		pageSize = MAX_PAGE_SIZE
	case pageSize <= 0:
		pageSize = DEFAULT_PAGE_SIZE
	}

	var countRows int64
	model.Count(&countRows)

	numPages := float32(countRows) / float32(pageSize)
	if float32(int(numPages)) < numPages {
		numPages++
	}

	if page <= 0 || page > int(numPages) {
		page = 1
	}

	return &Pagination{
		Page:     page,
		NumPages: int(numPages),
		PageSize: pageSize,
		Count:    countRows,
	}
}

func PaginatorFromRequest(r *http.Request, model *gorm.DB) *Pagination {
	page, _ := strconv.Atoi(r.URL.Query().Get(GET_OPTION_KEY_PAGE))
	pageSize, _ := strconv.Atoi(r.URL.Query().Get(GET_OPTION_KEY_PAGE_SIZE))

	return NewPagination(page, pageSize, model)
}
