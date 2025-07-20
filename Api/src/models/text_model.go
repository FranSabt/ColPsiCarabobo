package models

import (
	"time"
)

type TextModel struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Text      string    `gorm:"type:text" json:"description"`
	Active    bool      `gorm:"default:true;not null" json:"-"` // Correcto: Oculto del JSON
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type TextDTO struct {
	Text   string `json:"text"`
	Active bool   `json:"active"` // Puntero para que sea opcional en la actualización
}
