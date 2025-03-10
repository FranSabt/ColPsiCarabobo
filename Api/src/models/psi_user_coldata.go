package models

import (
	"time"

	"github.com/google/uuid"
)

type PsiUserColData struct {
	ID uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`

	// Undergraduate Data
	UniversityUndergraduate string    `gorm:"size:255"`
	GraduateDate            time.Time `gorm:"type:date"` // Usar time.Time para fechas
	MentionUndergraduate    string    `gorm:"size:255"`

	// Undergraduate Data Title Register
	RegisterTitleState string    `gorm:"size:255"` // Puedes usar un enum si lo defines
	RegisterTitleDate  time.Time `gorm:"type:date"`
	RegisterNumber     int
	RegisterFolio      string `gorm:"size:255"` // Puedes usar un enum si lo defines
	RegisterTome       string `gorm:"size:255"`

	// Professional Data
	GuildDirector       bool `gorm:"default:false"`
	SixtyFiveOrPlus     bool `gorm:"default:false"`
	GuildCollaborator   bool `gorm:"default:false"`
	PublicEmployee      bool `gorm:"default:false"`
	UniversityProfessor bool `gorm:"default:false"`

	// Otros campos
	DateOfLastSolvency time.Time `gorm:"type:date"`
	DoubleGuild        bool      `gorm:"default:false"`
	CPSM               bool      `gorm:"default:false"`

	// Relación con PsiUserModel
	PsiUserModelID uuid.UUID `gorm:"type:uuid"` // Clave foránea
	// PsiUserModel   PsiUserModel `gorm:"foreignKey:PsiUserModelID"`
}
