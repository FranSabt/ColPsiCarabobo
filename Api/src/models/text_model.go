package models

import (
	"time"

	"github.com/google/uuid"
)

type TextModel struct {
	ID         uint       `gorm:"primaryKey" json:"id"`
	Text       string     `gorm:"type:text" json:"description"`
	Active     bool       `gorm:"default:true;not null" json:"-"` // Correcto: Oculto del JSON
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
	CreateBy   string     `gorm:"size:255" json:"create_by"`
	UpdateBy   string     `gorm:"size:255" json:"update_by"`
	CreateById *uuid.UUID `gorm:"type:uuid" json:"create_by_id"`
	UpdateById *uuid.UUID `gorm:"type:uuid" json:"update_by_id"`
}

type TextDTO struct {
	Text       string     `json:"text"`
	Active     bool       `json:"active"` // Puntero para que sea opcional en la actualización
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
	CreateBy   string     `gorm:"size:255" json:"create_by"`
	UpdateBy   string     `gorm:"size:255" json:"update_by"`
	CreateById *uuid.UUID `gorm:"type:uuid" json:"create_by_id"`
	UpdateById *uuid.UUID `gorm:"type:uuid" json:"update_by_id"`
}
