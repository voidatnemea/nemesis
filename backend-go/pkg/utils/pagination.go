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
