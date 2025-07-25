package models

import (
	"time"

	"github.com/google/uuid"
)

type PsiUserColData struct {
	ID uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`

	// Undergraduate Data
	UniversityUndergraduate     string    `gorm:"size:255" json:"university_undergraduate"`
	ShowUniversityUndergraduate bool      `gorm:"default:false" json:"show_university_undergraduate"`
	GraduateDate                time.Time `gorm:"type:date" json:"graduate_date"` // Usar time.Time para fechas
	ShowGraduateDate            bool      `gorm:"default:false" json:"show_graduate_date"`
	MentionUndergraduate        string    `gorm:"size:255" json:"mention_undergraduate"`
	ShowMentionUndergraduate    bool      `gorm:"default:false" json:"show_mention_undergraduate"`

	// Undergraduate Data Title Register
	RegisterTitleState string    `gorm:"size:255" json:"register_title_state"` // Puedes usar un enum si lo defines
	RegisterTitleDate  time.Time `gorm:"type:date" json:"register_title_date"`
	RegisterNumber     int       `json:"register_number"`
	RegisterFolio      string    `gorm:"size:255" json:"register_folio"` // Puedes usar un enum si lo defines
	RegisterTome       string    `gorm:"size:255" json:"register_tome"`

	// Professional Data
	GuildDirector       bool `gorm:"default:false" json:"guild_director"`
	SixtyFiveOrPlus     bool `gorm:"default:false" json:"sixty_five_or_plus"`
	GuildCollaborator   bool `gorm:"default:false" json:"guild_collaborator"`
	PublicEmployee      bool `gorm:"default:false" json:"public_employee"`
	UniversityProfessor bool `gorm:"default:false" json:"university_professor"`

	// Otros campos
	DateOfLastSolvency time.Time `gorm:"type:date" json:"date_of_last_solvency"`
	DoubleGuild        bool      `gorm:"default:false" json:"double_guild"`
	CPSM               bool      `gorm:"default:false" json:"cpsm"`

	// Relación con PsiUserModel
	PsiUserModelID uuid.UUID `gorm:"type:uuid" json:"psi_user_model_id"` // Clave foránea
	// PsiUserModel   PsiUserModel `gorm:"foreignKey:PsiUserModelID"`

	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
	CreateBy   string     `gorm:"size:255" json:"create_by"`
	UpdateBy   string     `gorm:"size:255" json:"update_by"`
	CreateById *uuid.UUID `gorm:"type:uuid" json:"create_by_id"`
	UpdateById *uuid.UUID `gorm:"type:uuid" json:"update_by_id"`
}
