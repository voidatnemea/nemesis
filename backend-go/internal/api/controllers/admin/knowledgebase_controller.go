package admin

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mythicalsystems/nemesis-backend/internal/database"
	"github.com/mythicalsystems/nemesis-backend/internal/models"
	"github.com/mythicalsystems/nemesis-backend/pkg/utils"
)

type KnowledgebaseController struct{}

func (k *KnowledgebaseController) CategoriesIndex(c *gin.Context) {
	var cats []models.KnowledgebaseCategory
	database.DB.Order("position ASC").Find(&cats)
	utils.Success(c, cats, "Categories retrieved", http.StatusOK)
}

func (k *KnowledgebaseController) CategoriesShow(c *gin.Context) {
	id := c.Param("id")
	var cat models.KnowledgebaseCategory
	if err := database.DB.First(&cat, id).Error; err != nil {
		utils.Error(c, "Category not found", "NOT_FOUND", http.StatusNotFound, nil)
		return
	}
	utils.Success(c, cat, "Category retrieved", http.StatusOK)
}

func (k *KnowledgebaseController) CategoriesCreate(c *gin.Context) {
	var req struct {
		Name        string `json:"name" binding:"required"`
		Slug        string `json:"slug" binding:"required"`
		Icon        string `json:"icon" binding:"required"`
		Description string `json:"description"`
		Position    int    `json:"position"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, err.Error(), "INVALID_REQUEST", http.StatusBadRequest, nil)
		return
	}
	cat := models.KnowledgebaseCategory{
		Name: req.Name, Slug: req.Slug, Icon: req.Icon,
		Description: req.Description, Position: req.Position,
	}
	database.DB.Create(&cat)
	utils.Success(c, cat, "Category created", http.StatusCreated)
}

func (k *KnowledgebaseController) CategoriesUpdate(c *gin.Context) {
	id := c.Param("id")
	var cat models.KnowledgebaseCategory
	if err := database.DB.First(&cat, id).Error; err != nil {
		utils.Error(c, "Category not found", "NOT_FOUND", http.StatusNotFound, nil)
		return
	}
	var req map[string]interface{}
	c.ShouldBindJSON(&req)
	delete(req, "id")
	database.DB.Model(&cat).Updates(req)
	utils.Success(c, cat, "Category updated", http.StatusOK)
}

func (k *KnowledgebaseController) CategoriesDelete(c *gin.Context) {
	id := c.Param("id")
	database.DB.Delete(&models.KnowledgebaseCategory{}, id)
	utils.Success(c, nil, "Category deleted", http.StatusOK)
}

func (k *KnowledgebaseController) ArticlesIndex(c *gin.Context) {
	p := utils.GetPagination(c)
	categoryID := c.Query("category_id")
	status := c.Query("status")
	var articles []models.KnowledgebaseArticle
	var total int64
	q := database.DB.Model(&models.KnowledgebaseArticle{})
	if categoryID != "" {
		q = q.Where("category_id = ?", categoryID)
	}
	if status != "" {
		q = q.Where("status = ?", status)
	}
	q.Count(&total)
	q.Offset(p.Offset).Limit(p.PerPage).Find(&articles)
	utils.Success(c, gin.H{"data": articles, "total": total}, "Articles retrieved", http.StatusOK)
}

func (k *KnowledgebaseController) ArticlesShow(c *gin.Context) {
	id := c.Param("id")
	var article models.KnowledgebaseArticle
	if err := database.DB.First(&article, id).Error; err != nil {
		utils.Error(c, "Article not found", "NOT_FOUND", http.StatusNotFound, nil)
		return
	}
	var tags []models.KnowledgebaseArticleTag
	database.DB.Where("article_id = ?", article.ID).Find(&tags)
	utils.Success(c, gin.H{"article": article, "tags": tags}, "Article retrieved", http.StatusOK)
}

func (k *KnowledgebaseController) ArticlesCreate(c *gin.Context) {
	var req struct {
		CategoryID int    `json:"category_id" binding:"required"`
		Title      string `json:"title" binding:"required"`
		Slug       string `json:"slug" binding:"required"`
		Content    string `json:"content"`
		Status     string `json:"status"`
		Pinned     string `json:"pinned"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, err.Error(), "INVALID_REQUEST", http.StatusBadRequest, nil)
		return
	}
	ctxUser, _ := c.Get("user")
	user := ctxUser.(models.User)
	if req.Status == "" {
		req.Status = "draft"
	}
	article := models.KnowledgebaseArticle{
		CategoryID: req.CategoryID, Title: req.Title, Slug: req.Slug,
		Content: req.Content, Status: req.Status, Pinned: req.Pinned,
		AuthorUUID: user.UUID,
	}
	database.DB.Create(&article)
	utils.Success(c, article, "Article created", http.StatusCreated)
}

func (k *KnowledgebaseController) ArticlesUpdate(c *gin.Context) {
	id := c.Param("id")
	var article models.KnowledgebaseArticle
	if err := database.DB.First(&article, id).Error; err != nil {
		utils.Error(c, "Article not found", "NOT_FOUND", http.StatusNotFound, nil)
		return
	}
	var req map[string]interface{}
	c.ShouldBindJSON(&req)
	delete(req, "id")
	database.DB.Model(&article).Updates(req)
	utils.Success(c, article, "Article updated", http.StatusOK)
}

func (k *KnowledgebaseController) ArticlesDelete(c *gin.Context) {
	id := c.Param("id")
	database.DB.Delete(&models.KnowledgebaseArticle{}, id)
	utils.Success(c, nil, "Article deleted", http.StatusOK)
}

func (k *KnowledgebaseController) GetTags(c *gin.Context) {
	id := c.Param("id")
	var tags []models.KnowledgebaseArticleTag
	database.DB.Where("article_id = ?", id).Find(&tags)
	utils.Success(c, tags, "Tags retrieved", http.StatusOK)
}

func (k *KnowledgebaseController) CreateTag(c *gin.Context) {
	id := c.Param("id")
	var req struct {
		Tag string `json:"tag" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, err.Error(), "INVALID_REQUEST", http.StatusBadRequest, nil)
		return
	}
	var articleID int
	database.DB.Raw("SELECT id FROM featherpanel_knowledgebase_articles WHERE id = ?", id).Scan(&articleID)
	tag := models.KnowledgebaseArticleTag{ArticleID: articleID, Tag: req.Tag}
	database.DB.Create(&tag)
	utils.Success(c, tag, "Tag created", http.StatusCreated)
}

func (k *KnowledgebaseController) DeleteTag(c *gin.Context) {
	tagID := c.Param("tagId")
	database.DB.Delete(&models.KnowledgebaseArticleTag{}, tagID)
	utils.Success(c, nil, "Tag deleted", http.StatusOK)
}
