package repositories

import (
	"gosearch/models"
)

type Finder interface {
	FindByURL(url string) (*models.BaseModel, error)
}

type Saver interface {
	Save(doc *models.Document) error
}

type Creator interface {
	FindOrCreate(body []byte, attrs map[string]any, level int) error
}
