package document_repository

import (
	"fmt"
	"strings"

	"gosearch/models"

	"gosearch/internal/global"
	"gosearch/internal/types"
)

type DocumentRepository interface {
	FindByURL(url string) (*models.Document, error)
	Save(doc *models.Document) error
	FindOrCreate(body []uint8, lg *types.LinkGrouped, level int) error
}

func FindOrCreate(body []uint8, lg *types.LinkGrouped, level int) error {
	db := global.Container.DB

	url := lg.Url
	title := strings.Join(lg.Titles, ", ")

	doc := models.Document{
		URL:   url,
		Title: title,
		Body:  string(body),
		Level: level,
	}

	err := db.Where(models.Document{URL: url, Title: title}).FirstOrCreate(&doc).Error
	if err != nil {
		return fmt.Errorf("failed to create document > %w", err)
	}

	return nil
}
