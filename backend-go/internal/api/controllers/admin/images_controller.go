package admin

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mythicalsystems/nemesis-backend/internal/database"
	"github.com/mythicalsystems/nemesis-backend/internal/models"
	"github.com/mythicalsystems/nemesis-backend/pkg/utils"
)

type ImagesController struct{}

func (i *ImagesController) Index(c *gin.Context) {
	p := utils.GetPagination(c)
	var images []models.Image
	var total int64
	database.DB.Model(&models.Image{}).Count(&total)
	database.DB.Offset(p.Offset).Limit(p.PerPage).Find(&images)
	utils.Success(c, gin.H{"images": images, "pagination": utils.BuildPagination(p, total)}, "Images retrieved", http.StatusOK)
}

func (i *ImagesController) Show(c *gin.Context) {
	id := c.Param("id")
	var image models.Image
	if err := database.DB.First(&image, id).Error; err != nil {
		utils.Error(c, "Image not found", "NOT_FOUND", http.StatusNotFound, nil)
		return
	}
	utils.Success(c, image, "Image retrieved", http.StatusOK)
}

func (i *ImagesController) Create(c *gin.Context) {
	var req struct {
		Name string `json:"name" binding:"required"`
		URL  string `json:"url" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, err.Error(), "INVALID_REQUEST", http.StatusBadRequest, nil)
		return
	}
	image := models.Image{Name: req.Name, URL: req.URL}
	database.DB.Create(&image)
	utils.Success(c, image, "Image created", http.StatusCreated)
}

func (i *ImagesController) Update(c *gin.Context) {
	id := c.Param("id")
	var image models.Image
	if err := database.DB.First(&image, id).Error; err != nil {
		utils.Error(c, "Image not found", "NOT_FOUND", http.StatusNotFound, nil)
		return
	}
	var req map[string]interface{}
	c.ShouldBindJSON(&req)
	delete(req, "id")
	database.DB.Model(&image).Updates(req)
	utils.Success(c, image, "Image updated", http.StatusOK)
}

func (i *ImagesController) Delete(c *gin.Context) {
	id := c.Param("id")
	database.DB.Delete(&models.Image{}, id)
	utils.Success(c, nil, "Image deleted", http.StatusOK)
}
