package user

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mythicalsystems/nemesis-backend/pkg/utils"
)

type VmInstancesController struct{}

func (v *VmInstancesController) Index(c *gin.Context) {
	utils.Success(c, gin.H{
		"instances": []interface{}{},
		"pagination": gin.H{
			"current_page":  1,
			"per_page":      25,
			"total_records": 0,
			"total_pages":   1,
			"has_next":      false,
			"has_prev":      false,
			"from":          0,
			"to":            0,
		},
	}, "VM instances retrieved", http.StatusOK)
}
