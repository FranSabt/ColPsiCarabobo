package models

import (
	"time"

	"github.com/google/uuid"
)

type ProfilePicModel struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	UserID    uuid.UUID `gorm:"type:uuid;not null" json:"user_id"` // o UserID, ProductID, etc.
	Name      string    `gorm:"size:50" json:"name"`
	ImageData []byte    `gorm:"type:bytea;not null" json:"image_data"`
	Format    string    `gorm:"size:10" json:"format"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type PostPicModel struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	PostId    uint      `json:"post_id,omitempty"`
	Name      string    `gorm:"size:50" json:"name"`
	ImageData []byte    `gorm:"type:bytea;not null" json:"image_data"`
	Format    string    `gorm:"size:10" json:"format"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
