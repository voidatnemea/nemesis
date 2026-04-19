package models

import "time"

type KnowledgebaseCategory struct {
	ID          uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Name        string    `gorm:"type:varchar(255)" json:"name"`
	Slug        string    `gorm:"uniqueIndex;type:varchar(255)" json:"slug"`
	Icon        string    `gorm:"type:varchar(255)" json:"icon"`
	Description string    `gorm:"type:text" json:"description"`
	Position    int       `gorm:"default:0" json:"position"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (KnowledgebaseCategory) TableName() string { return "featherpanel_knowledgebase_categories" }

type KnowledgebaseArticle struct {
	ID         uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	CategoryID int       `gorm:"index" json:"category_id"`
	Title      string    `gorm:"type:varchar(255)" json:"title"`
	Slug       string    `gorm:"uniqueIndex;type:varchar(255)" json:"slug"`
	Content    string    `gorm:"type:longtext" json:"content"`
	Status     string    `gorm:"type:enum('draft','published','archived');default:'draft'" json:"status"`
	Pinned     string    `gorm:"type:enum('true','false');default:'false'" json:"pinned"`
	Views      int64     `gorm:"default:0" json:"views"`
	AuthorUUID string    `gorm:"type:varchar(36)" json:"author_uuid"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

func (KnowledgebaseArticle) TableName() string { return "featherpanel_knowledgebase_articles" }

type KnowledgebaseArticleTag struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	ArticleID int       `gorm:"index" json:"article_id"`
	Tag       string    `gorm:"type:varchar(100)" json:"tag"`
	CreatedAt time.Time `json:"created_at"`
}

func (KnowledgebaseArticleTag) TableName() string { return "featherpanel_knowledgebase_articles_tags" }

type KnowledgebaseArticleAttachment struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	ArticleID int       `gorm:"index" json:"article_id"`
	FileName  string    `gorm:"type:varchar(255)" json:"file_name"`
	FilePath  string    `gorm:"type:varchar(500)" json:"file_path"`
	FileSize  int64     `json:"file_size"`
	FileType  string    `gorm:"type:varchar(100)" json:"file_type"`
	CreatedAt time.Time `json:"created_at"`
}

func (KnowledgebaseArticleAttachment) TableName() string {
	return "featherpanel_knowledgebase_articles_attachments"
}
