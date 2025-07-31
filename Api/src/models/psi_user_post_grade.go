package models

import (
	"time" // Importa el paquete time para los timestamps

	"github.com/google/uuid"
	"gorm.io/gorm" // Importa gorm para gorm.Model
)

type PisUserPostGrade struct {
	// gorm.Model incrustado para campos estándar como ID, CreatedAt, UpdatedAt, DeletedAt
	// Si prefieres usar tu propio ID de UUID, puedes omitir gorm.Model y usar el ID que ya tenías.
	// Vamos a usar tu ID original y añadir los timestamps manualmente, que es una práctica común.

	ID uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`

	// Clave foranea
	// Relación con PsiUser
	PsiUserID uuid.UUID `gorm:"type:uuid" json:"psi_user_id"` // Clave foránea

	PostGradeTitle          string `gorm:"column:post_grade_title;not null;size:255" json:"post_grade_title"`
	PostGradeUniversity     string `gorm:"column:post_grade_university;not null;size:255" json:"post_grade_university"`
	PostGradeGraduationYear string `gorm:"column:post_grade_graduation_year;not null;size:50" json:"post_grade_graduation_year"`
	PostGradeDescription    string `gorm:"column:post_grade_description;type:text" json:"post_grade_description"`
	Active                  bool   `gorm:"default:true;not null" json:"-"` // Correcto: Oculto del JSON

	// Estos campos son probablemente claves foráneas a otra tabla (ej: 'pictures' o 'files')
	// El puntero (*) indica que son "nullable" (pueden ser nulos en la BBDD)
	PicOne   *uuid.UUID `gorm:"column:pic_one;type:uuid" json:"pic_one"`
	PicTwo   *uuid.UUID `gorm:"column:pic_two;type:uuid" json:"pic_two"`
	PicThree *uuid.UUID `gorm:"column:pic_three;type:uuid" json:"pic_three"`

	// Timestamps para auditoría
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"` // Para soft delete, se omite en JSON
}

// TableName especifica explícitamente el nombre de la tabla en la base de datos.
// Es una buena práctica para evitar que GORM adivine el nombre.
func (PisUserPostGrade) TableName() string {
	return "pis_user_post_grades"
}
