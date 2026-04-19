package utils

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

type Pagination struct {
	Page    int
	PerPage int
	Offset  int
}

func GetPagination(c *gin.Context) Pagination {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "20"))
	if page < 1 {
		page = 1
	}
	if perPage < 1 || perPage > 100 {
		perPage = 20
	}
	return Pagination{
		Page:    page,
		PerPage: perPage,
		Offset:  (page - 1) * perPage,
	}
}

// BuildPagination returns the standard pagination map the frontend expects.
func BuildPagination(p Pagination, total int64) gin.H {
	totalPages := int((total + int64(p.PerPage) - 1) / int64(p.PerPage))
	if totalPages < 1 {
		totalPages = 1
	}
	from := p.Offset + 1
	to := p.Offset + p.PerPage
	if int64(to) > total {
		to = int(total)
	}
	if total == 0 {
		from = 0
		to = 0
	}
	return gin.H{
		"current_page":  p.Page,
		"per_page":      p.PerPage,
		"total_records": total,
		"total_pages":   totalPages,
		"has_next":      p.Page < totalPages,
		"has_prev":      p.Page > 1,
		"from":          from,
		"to":            to,
	}
}
