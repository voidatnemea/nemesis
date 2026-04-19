package user

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mythicalsystems/nemesis-backend/internal/database"
	"github.com/mythicalsystems/nemesis-backend/internal/models"
	"github.com/mythicalsystems/nemesis-backend/pkg/utils"
)

type UserKnowledgebaseController struct{}

func (k *UserKnowledgebaseController) GetCategories(c *gin.Context) {
	var cats []models.KnowledgebaseCategory
	database.DB.Order("position ASC").Find(&cats)
	utils.Success(c, cats, "Categories retrieved", http.StatusOK)
}

func (k *UserKnowledgebaseController) GetArticles(c *gin.Context) {
	categoryID := c.Query("category_id")
	search := c.Query("search")
	var articles []models.KnowledgebaseArticle
	q := database.DB.Where("status = ?", "published")
	if categoryID != "" {
		q = q.Where("category_id = ?", categoryID)
	}
	if search != "" {
		q = q.Where("title LIKE ? OR content LIKE ?", "%"+search+"%", "%"+search+"%")
	}
	q.Find(&articles)
	utils.Success(c, articles, "Articles retrieved", http.StatusOK)
}

func (k *UserKnowledgebaseController) GetArticle(c *gin.Context) {
	id := c.Param("id")
	var article models.KnowledgebaseArticle
	if err := database.DB.Where("(id = ? OR slug = ?) AND status = ?", id, id, "published").First(&article).Error; err != nil {
		utils.Error(c, "Article not found", "NOT_FOUND", http.StatusNotFound, nil)
		return
	}
	// Increment view count
	database.DB.Model(&article).UpdateColumn("views", article.Views+1)
	var tags []models.KnowledgebaseArticleTag
	database.DB.Where("article_id = ?", article.ID).Find(&tags)
	utils.Success(c, gin.H{"article": article, "tags": tags}, "Article retrieved", http.StatusOK)
}
