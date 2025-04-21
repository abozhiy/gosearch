package models

type Document struct {
	ID    int        `gorm:"primaryKey"`
	URL   string
	Title string
	Body  string     `gorm:"type:text"`
	Level int
}
