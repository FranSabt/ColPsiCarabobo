package models

import "time"

// PsiSpecialty representa una especialidad en la base de datos.
// El campo 'Active' se usa para lógica interna (soft delete) pero no se expone en la API.
type PsiSpecialty struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Name        string    `gorm:"size:50;not null;unique" json:"name"`
	Description string    `gorm:"type:text" json:"description"`
	Active      bool      `gorm:"default:true;not null" json:"-"` // Correcto: Oculto del JSON
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
