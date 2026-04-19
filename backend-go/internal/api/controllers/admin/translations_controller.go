package admin

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/mythicalsystems/nemesis-backend/internal/services"
	"github.com/mythicalsystems/nemesis-backend/pkg/utils"
)

type TranslationsController struct{}

type translationFile struct {
	Code     string `json:"code"`
	Name     string `json:"name"`
	File     string `json:"file"`
	Size     int    `json:"size"`
	Modified string `json:"modified"`
	Enabled  bool   `json:"enabled"`
}

var langNames = map[string]string{
	"en": "English",
	"es": "Español",
	"de": "Deutsch",
	"fr": "Français",
	"it": "Italiano",
	"pt": "Português",
	"ru": "Русский",
	"tr": "Türkçe",
	"zh": "中文",
	"ja": "日本語",
	"ko": "한국어",
	"pl": "Polski",
	"nl": "Nederlands",
	"sv": "Svenska",
	"uk": "Українська",
}

func (tr *TranslationsController) Index(c *gin.Context) {
	raw := services.GetSetting("translations:files", "")
	var files []translationFile
	if raw != "" {
		_ = json.Unmarshal([]byte(raw), &files)
	}
	if len(files) == 0 {
		files = []translationFile{
			{Code: "en", Name: "English", File: "en.json", Size: 0, Enabled: true},
		}
	}
	utils.Success(c, files, "Translation files retrieved", http.StatusOK)
}

func (tr *TranslationsController) Get(c *gin.Context) {
	lang := c.Param("lang")
	raw := services.GetSetting("translations:content:"+lang, "")
	if raw == "" {
		utils.Success(c, gin.H{}, "Empty translation", http.StatusOK)
		return
	}
	var content map[string]interface{}
	if err := json.Unmarshal([]byte(raw), &content); err != nil {
		utils.Error(c, "Corrupted translation file", "CORRUPT", http.StatusInternalServerError, nil)
		return
	}
	utils.Success(c, content, "Translation retrieved", http.StatusOK)
}

func (tr *TranslationsController) Save(c *gin.Context) {
	lang := c.Param("lang")
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		utils.Error(c, err.Error(), "INVALID_REQUEST", http.StatusBadRequest, nil)
		return
	}
	if !json.Valid(body) {
		utils.Error(c, "Payload is not valid JSON", "INVALID_JSON", http.StatusBadRequest, nil)
		return
	}
	services.SetSetting("translations:content:"+lang, string(body))
	tr.ensureInIndex(lang, len(body))
	utils.Success(c, gin.H{"lang": lang, "size": len(body)}, "Translation saved", http.StatusOK)
}

func (tr *TranslationsController) Create(c *gin.Context) {
	lang := c.Param("lang")
	if lang == "" {
		utils.Error(c, "Language code required", "INVALID_REQUEST", http.StatusBadRequest, nil)
		return
	}
	services.SetSetting("translations:content:"+lang, "{}")
	tr.ensureInIndex(lang, 2)
	utils.Success(c, gin.H{"lang": lang}, "Language created", http.StatusOK)
}

func (tr *TranslationsController) Delete(c *gin.Context) {
	lang := c.Param("lang")
	services.SetSetting("translations:content:"+lang, "")
	tr.removeFromIndex(lang)
	utils.Success(c, nil, "Language deleted", http.StatusOK)
}

func (tr *TranslationsController) Upload(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		utils.Error(c, "File required", "INVALID_REQUEST", http.StatusBadRequest, nil)
		return
	}
	lang := strings.TrimSuffix(file.Filename, ".json")
	f, err := file.Open()
	if err != nil {
		utils.Error(c, err.Error(), "READ_FAIL", http.StatusBadRequest, nil)
		return
	}
	defer f.Close()
	body, _ := io.ReadAll(f)
	if !json.Valid(body) {
		utils.Error(c, "File is not valid JSON", "INVALID_JSON", http.StatusBadRequest, nil)
		return
	}
	services.SetSetting("translations:content:"+lang, string(body))
	tr.ensureInIndex(lang, len(body))
	utils.Success(c, gin.H{"lang": lang}, "Translation uploaded", http.StatusOK)
}

func (tr *TranslationsController) Download(c *gin.Context) {
	lang := c.Param("lang")
	raw := services.GetSetting("translations:content:"+lang, "{}")
	c.Header("Content-Disposition", "attachment; filename="+lang+".json")
	c.Data(http.StatusOK, "application/json", []byte(raw))
}

func (tr *TranslationsController) ensureInIndex(lang string, size int) {
	raw := services.GetSetting("translations:files", "")
	var files []translationFile
	if raw != "" {
		_ = json.Unmarshal([]byte(raw), &files)
	}
	found := false
	for i, f := range files {
		if f.Code == lang {
			files[i].Size = size
			files[i].Enabled = true
			found = true
			break
		}
	}
	if !found {
		name := langNames[lang]
		if name == "" {
			name = strings.ToUpper(lang)
		}
		files = append(files, translationFile{Code: lang, Name: name, File: lang + ".json", Size: size, Enabled: true})
	}
	b, _ := json.Marshal(files)
	services.SetSetting("translations:files", string(b))
}

func (tr *TranslationsController) removeFromIndex(lang string) {
	raw := services.GetSetting("translations:files", "")
	var files []translationFile
	if raw != "" {
		_ = json.Unmarshal([]byte(raw), &files)
	}
	out := make([]translationFile, 0, len(files))
	for _, f := range files {
		if f.Code != lang {
			out = append(out, f)
		}
	}
	b, _ := json.Marshal(out)
	services.SetSetting("translations:files", string(b))
}
