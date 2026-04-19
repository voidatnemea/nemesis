package user

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mythicalsystems/nemesis-backend/internal/database"
	"github.com/mythicalsystems/nemesis-backend/internal/models"
	"github.com/mythicalsystems/nemesis-backend/pkg/utils"
)

type ApiClientController struct{}

func (a *ApiClientController) GetApiClients(c *gin.Context) {
	ctxUser, _ := c.Get("user")
	user := ctxUser.(models.User)
	var clients []models.ApiClient
	database.DB.Where("user_uuid = ?", user.UUID).Find(&clients)
	utils.Success(c, clients, "API clients retrieved", http.StatusOK)
}

func (a *ApiClientController) GetApiClient(c *gin.Context) {
	ctxUser, _ := c.Get("user")
	user := ctxUser.(models.User)
	id := c.Param("id")
	var client models.ApiClient
	if err := database.DB.Where("id = ? AND user_uuid = ?", id, user.UUID).First(&client).Error; err != nil {
		utils.Error(c, "API client not found", "NOT_FOUND", http.StatusNotFound, nil)
		return
	}
	utils.Success(c, client, "API client retrieved", http.StatusOK)
}

func (a *ApiClientController) CreateApiClient(c *gin.Context) {
	ctxUser, _ := c.Get("user")
	user := ctxUser.(models.User)
	var req struct {
		Name       string `json:"name" binding:"required"`
		AllowedIps string `json:"allowed_ips"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, err.Error(), "INVALID_REQUEST", http.StatusBadRequest, nil)
		return
	}
	publicKey := "pub_" + utils.GenerateRandomString(16)
	privateKey := "priv_" + utils.GenerateRandomString(32)
	client := models.ApiClient{
		UserUUID:   user.UUID,
		Name:       req.Name,
		PublicKey:  publicKey,
		PrivateKey: privateKey,
		AllowedIps: req.AllowedIps,
	}
	database.DB.Create(&client)
	// Return private key once on creation
	utils.Success(c, gin.H{
		"id":          client.ID,
		"name":        client.Name,
		"public_key":  client.PublicKey,
		"private_key": privateKey,
		"allowed_ips": client.AllowedIps,
		"created_at":  client.CreatedAt,
	}, "API client created", http.StatusCreated)
}

func (a *ApiClientController) UpdateApiClient(c *gin.Context) {
	ctxUser, _ := c.Get("user")
	user := ctxUser.(models.User)
	id := c.Param("id")
	var client models.ApiClient
	if err := database.DB.Where("id = ? AND user_uuid = ?", id, user.UUID).First(&client).Error; err != nil {
		utils.Error(c, "API client not found", "NOT_FOUND", http.StatusNotFound, nil)
		return
	}
	var req struct {
		Name       string `json:"name"`
		AllowedIps string `json:"allowed_ips"`
	}
	c.ShouldBindJSON(&req)
	if req.Name != "" {
		client.Name = req.Name
	}
	client.AllowedIps = req.AllowedIps
	database.DB.Save(&client)
	utils.Success(c, client, "API client updated", http.StatusOK)
}

func (a *ApiClientController) DeleteApiClient(c *gin.Context) {
	ctxUser, _ := c.Get("user")
	user := ctxUser.(models.User)
	id := c.Param("id")
	database.DB.Where("id = ? AND user_uuid = ?", id, user.UUID).Delete(&models.ApiClient{})
	utils.Success(c, nil, "API client deleted", http.StatusOK)
}

func (a *ApiClientController) RegenerateApiKeys(c *gin.Context) {
	ctxUser, _ := c.Get("user")
	user := ctxUser.(models.User)
	id := c.Param("id")
	var client models.ApiClient
	if err := database.DB.Where("id = ? AND user_uuid = ?", id, user.UUID).First(&client).Error; err != nil {
		utils.Error(c, "API client not found", "NOT_FOUND", http.StatusNotFound, nil)
		return
	}
	newPublic := "pub_" + utils.GenerateRandomString(16)
	newPrivate := "priv_" + utils.GenerateRandomString(32)
	database.DB.Model(&client).Updates(map[string]interface{}{
		"public_key":  newPublic,
		"private_key": newPrivate,
	})
	utils.Success(c, gin.H{
		"public_key":  newPublic,
		"private_key": newPrivate,
	}, "API keys regenerated", http.StatusOK)
}
