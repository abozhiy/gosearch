package repositories

import (
	"fmt"
	"strings"

	"gosearch/models"

	"gosearch/internal/global"
	"gosearch/internal/types"

	"gorm.io/gorm"
)

type DocumentRepository interface {
    Creator
}

type documentRepo struct {
	db *gorm.DB
}

func NewDocumentRepository() *documentRepo {
	return &documentRepo{db: global.Container.DB}
}

func (repo *documentRepo) FindOrCreate(body []byte, attrs map[string]any, level int) error {
	raw, ok := attrs["lg"]
	if !ok {
		return fmt.Errorf("lg key not found")
	}

	lg, ok := raw.(*types.LinkGrouped)
	if !ok {
		return fmt.Errorf("lg has wrong type")
	}

	url := lg.Url
	title := strings.Join(lg.Titles, ", ")

	doc := models.Document{
		URL:   url,
		Title: title,
		Body:  string(body),
		Level: level,
	}

	err := repo.db.Where(models.Document{URL: url, Title: title}).FirstOrCreate(&doc).Error
	if err != nil {
		return fmt.Errorf("failed to create document > %w", err)
	}

	return nil
}
