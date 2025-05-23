package models

import (
	"time"

	"github.com/google/uuid"
)

type ImagesModel struct {
	ID           uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	AssociatedID uuid.UUID `gorm:"type:uuid;not null"` // o UserID, ProductID, etc.
	Name         string    `gorm:"size:50"`
	ImageData    []byte    `gorm:"type:bytea;not null"`
	Format       string    `gorm:"size:10"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
