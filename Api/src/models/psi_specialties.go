package models

import "time"

type PsiSpecialty struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Name        string    `gorm:"size:50;not null;unique" json:"name"`
	Description string    `gorm:"type:text" json:"description"` // Tag JSON añadido
	Active      bool      `gorm:"type:bool" json:"_"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
