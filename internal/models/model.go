package models

import "time"

type Component struct {
	ID            *string `gorm:"type:text;"`
	Name          *string `gorm:"type:text;"`
	Type          *string `gorm:"type:text;"`
	CreatedBy     *string `gorm:"type:integer;"`
	CreatedAt     *time.Time
	LastUpdatedAt *time.Time
}
