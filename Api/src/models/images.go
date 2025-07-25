package models

import (
	"time"

	"github.com/google/uuid"
)

type ProfilePicModel struct {
	ID         uuid.UUID  `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	UserID     uuid.UUID  `gorm:"type:uuid;not null" json:"user_id"` // o UserID, ProductID, etc.
	Name       string     `gorm:"size:50" json:"name"`
	ImageData  []byte     `gorm:"type:bytea;not null" json:"image_data"`
	Format     string     `gorm:"size:10" json:"format"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
	CreateBy   string     `gorm:"size:255" json:"create_by"`
	UpdateBy   string     `gorm:"size:255" json:"update_by"`
	CreateById *uuid.UUID `gorm:"type:uuid" json:"create_by_id"`
	UpdateById *uuid.UUID `gorm:"type:uuid" json:"update_by_id"`
}

func (ProfilePicModel) TableName() string {
	return "profile_pics"
}

type PostPicModel struct {
	ID         uuid.UUID  `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	PostId     uint       `json:"post_id,omitempty"`
	Name       string     `gorm:"size:50" json:"name"`
	ImageData  []byte     `gorm:"type:bytea;not null" json:"image_data"`
	Format     string     `gorm:"size:10" json:"format"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
	CreateBy   string     `gorm:"size:255" json:"create_by"`
	UpdateBy   string     `gorm:"size:255" json:"update_by"`
	CreateById *uuid.UUID `gorm:"type:uuid" json:"create_by_id"`
	UpdateById *uuid.UUID `gorm:"type:uuid" json:"update_by_id"`
}

func (PostPicModel) TableName() string {
	return "post_pics"
}

type PostGradePic struct {
	ID         uuid.UUID  `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	Name       string     `gorm:"size:50" json:"name"`
	ImageData  []byte     `gorm:"type:bytea;not null" json:"image_data"`
	Format     string     `gorm:"size:10" json:"format"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
	CreateBy   string     `gorm:"size:255" json:"create_by"`
	UpdateBy   string     `gorm:"size:255" json:"update_by"`
	CreateById *uuid.UUID `gorm:"type:uuid" json:"create_by_id"`
	UpdateById *uuid.UUID `gorm:"type:uuid" json:"update_by_id"`
}

func (PostGradePic) TableName() string {
	return "post_grade_pics"
}
