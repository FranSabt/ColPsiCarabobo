package models

import (
	"time"

	"github.com/google/uuid"
)

// PsiSpecialty representa una especialidad en la base de datos.
// El campo 'Active' se usa para lógica interna (soft delete) pero no se expone en la API.
type PsiSpecialty struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Name        string    `gorm:"size:50;not null;unique" json:"name"`
	Description string    `gorm:"type:text" json:"description"`
	Active      bool      `gorm:"default:true;not null" json:"active"` // Correcto: Oculto del JSON
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	CreateBy    string    `gorm:"size:255" json:"create_by"`
	UpdateBy    string    `gorm:"size:255" json:"update_by"`
	CreateById  uuid.UUID `gorm:"type:uuid" json:"create_by_id"`
	UpdateById  uuid.UUID `gorm:"type:uuid" json:"update_by_id"`
}
