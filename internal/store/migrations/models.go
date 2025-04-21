package migrations

import (
	"gosearch/models"
)

func Models() []any {
	return []any{
		&models.Document{},
	}
}
