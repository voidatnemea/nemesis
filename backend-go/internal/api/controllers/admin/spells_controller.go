package admin

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mythicalsystems/nemesis-backend/internal/database"
	"github.com/mythicalsystems/nemesis-backend/internal/models"
	"github.com/mythicalsystems/nemesis-backend/pkg/utils"
)

type SpellsController struct{}

func (s *SpellsController) Index(c *gin.Context) {
	p := utils.GetPagination(c)
	search := c.Query("search")
	var spells []models.Spell
	var total int64
	q := database.DB.Model(&models.Spell{})
	if search != "" {
		q = q.Where("name LIKE ?", "%"+search+"%")
	}
	q.Count(&total)
	q.Offset(p.Offset).Limit(p.PerPage).Find(&spells)
	utils.Success(c, gin.H{"data": spells, "total": total}, "Spells retrieved", http.StatusOK)
}

func (s *SpellsController) Show(c *gin.Context) {
	id := c.Param("id")
	var spell models.Spell
	if err := database.DB.First(&spell, id).Error; err != nil {
		utils.Error(c, "Spell not found", "NOT_FOUND", http.StatusNotFound, nil)
		return
	}
	var vars []models.SpellVariable
	database.DB.Where("spell_id = ?", spell.ID).Find(&vars)
	utils.Success(c, gin.H{"spell": spell, "variables": vars}, "Spell retrieved", http.StatusOK)
}

func (s *SpellsController) Create(c *gin.Context) {
	var req struct {
		RealmID      int    `json:"realm_id" binding:"required"`
		Author       string `json:"author" binding:"required"`
		Name         string `json:"name" binding:"required"`
		Description  string `json:"description"`
		DockerImages string `json:"docker_images"`
		Startup      string `json:"startup"`
		Config       string `json:"config"`
		Scripts      string `json:"scripts"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, err.Error(), "INVALID_REQUEST", http.StatusBadRequest, nil)
		return
	}
	spell := models.Spell{
		UUID:         utils.GenerateUUID(),
		RealmID:      req.RealmID,
		Author:       req.Author,
		Name:         req.Name,
		Description:  req.Description,
		DockerImages: req.DockerImages,
		Startup:      req.Startup,
		Config:       req.Config,
		Scripts:      req.Scripts,
	}
	database.DB.Create(&spell)
	utils.Success(c, spell, "Spell created", http.StatusCreated)
}

func (s *SpellsController) Update(c *gin.Context) {
	id := c.Param("id")
	var spell models.Spell
	if err := database.DB.First(&spell, id).Error; err != nil {
		utils.Error(c, "Spell not found", "NOT_FOUND", http.StatusNotFound, nil)
		return
	}
	var req map[string]interface{}
	c.ShouldBindJSON(&req)
	delete(req, "id")
	delete(req, "uuid")
	database.DB.Model(&spell).Updates(req)
	utils.Success(c, spell, "Spell updated", http.StatusOK)
}

func (s *SpellsController) Delete(c *gin.Context) {
	id := c.Param("id")
	database.DB.Delete(&models.Spell{}, id)
	utils.Success(c, nil, "Spell deleted", http.StatusOK)
}

func (s *SpellsController) GetByRealm(c *gin.Context) {
	realmID := c.Param("realmId")
	var spells []models.Spell
	database.DB.Where("realm_id = ?", realmID).Find(&spells)
	utils.Success(c, spells, "Spells retrieved", http.StatusOK)
}

func (s *SpellsController) ListVariables(c *gin.Context) {
	spellID := c.Param("id")
	var vars []models.SpellVariable
	database.DB.Where("spell_id = ?", spellID).Find(&vars)
	utils.Success(c, vars, "Variables retrieved", http.StatusOK)
}

func (s *SpellsController) CreateVariable(c *gin.Context) {
	spellID := c.Param("id")
	var req struct {
		Name         string `json:"name" binding:"required"`
		Description  string `json:"description"`
		EnvVariable  string `json:"env_variable" binding:"required"`
		DefaultValue string `json:"default_value"`
		UserViewable string `json:"user_viewable"`
		UserEditable string `json:"user_editable"`
		Rules        string `json:"rules"`
		FieldType    string `json:"field_type"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, err.Error(), "INVALID_REQUEST", http.StatusBadRequest, nil)
		return
	}
	// Parse spellID to int
	var spellIDInt int
	database.DB.Raw("SELECT id FROM featherpanel_spells WHERE id = ?", spellID).Scan(&spellIDInt)
	sv := models.SpellVariable{
		SpellID: spellIDInt, Name: req.Name, Description: req.Description,
		EnvVariable: req.EnvVariable, DefaultValue: req.DefaultValue,
		UserViewable: req.UserViewable, UserEditable: req.UserEditable,
		Rules: req.Rules, FieldType: req.FieldType,
	}
	if sv.FieldType == "" {
		sv.FieldType = "text"
	}
	database.DB.Create(&sv)
	utils.Success(c, sv, "Variable created", http.StatusCreated)
}

func (s *SpellsController) UpdateVariable(c *gin.Context) {
	id := c.Param("id")
	var sv models.SpellVariable
	if err := database.DB.First(&sv, id).Error; err != nil {
		utils.Error(c, "Variable not found", "NOT_FOUND", http.StatusNotFound, nil)
		return
	}
	var req map[string]interface{}
	c.ShouldBindJSON(&req)
	delete(req, "id")
	database.DB.Model(&sv).Updates(req)
	utils.Success(c, sv, "Variable updated", http.StatusOK)
}

func (s *SpellsController) DeleteVariable(c *gin.Context) {
	id := c.Param("id")
	database.DB.Delete(&models.SpellVariable{}, id)
	utils.Success(c, nil, "Variable deleted", http.StatusOK)
}
