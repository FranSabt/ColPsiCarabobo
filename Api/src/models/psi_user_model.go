package models

import (
	"time"

	"github.com/google/uuid"
)

type PsiUserModel struct {
	ID       uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	Username string    `gorm:"size:25;not null" json:"username"`
	Email    string    `gorm:"size:50;not null" json:"email"`
	Password string    `gorm:"size:512;not null" json:"password"`
	Key      string    `gorm:"size:512;" json:"key"`

	// Identity
	FirstName      string    `gorm:"size:255;not null" json:"first_name"`
	SecondName     string    `gorm:"size:255" json:"second_name"`
	LastName       string    `gorm:"size:255;not null" json:"last_name"`
	SecondLastName string    `gorm:"size:255" json:"second_last_name"`
	FPV            int       `gorm:"not null" json:"fpv"`
	CI             int       `gorm:"not null;unique" json:"ci"`
	Nationality    string    `gorm:"size:1;not null" json:"nationality"`
	BornDate       time.Time `gorm:"type:date;not null" json:"born_date"`
	Genre          string    `gorm:"size:1;not null" json:"genre"`

	// Contact
	ContactEmail             string `gorm:"size:255;not null" json:"contact_email"`
	ShowContactEmail         bool   `gorm:"default:false" json:"show_contact_email"`
	PublicPhone              string `gorm:"size:20" json:"public_phone"`
	ShowPublicPhone          bool   `gorm:"default:false" json:"show_public_phone"`
	ServiceAddress           string `gorm:"size:255" json:"service_address"`
	ShowPublicServiceAddress bool   `gorm:"default:false" json:"show_public_service_address"`

	// PsiCol
	Solvent     bool `gorm:"default:false" json:"solvent"`
	ProofOfLife bool `gorm:"default:false" json:"proof_of_life"`

	// Carabobo Direction
	MunicipalityCarabobo string `gorm:"size:255" json:"municipality_carabobo"`
	PhoneCarabobo        string `gorm:"size:20" json:"phone_carabobo"`
	CelPhoneCarabobo     string `gorm:"size:20" json:"cel_phone_carabobo"`

	// Outside Carabobo Direction
	StateOutside                string `gorm:"size:255" json:"state_outside"`
	MunicipalityOutSideCarabobo string `gorm:"size:255" json:"municipality_outside_carabobo"`
	PhoneOutSideCarabobo        string `gorm:"size:20" json:"phone_outside_carabobo"`
	CelPhoneOutSideCarabobo     string `gorm:"size:20" json:"cel_phone_outside_carabobo"`

	// Especialidades (campos opcionales)
	PrimarySpecialty   string `gorm:"size:50" json:"primary_specialty"`   // Puede ser nulo
	SecondarySpecialty string `gorm:"size:50" json:"secondary_specialty"` // Puede ser nulo

	// Relación con PsiUserColData
	PsiUserColDataID *uuid.UUID `gorm:"type:uuid" json:"psi_user_col_data_id"` // Clave foránea
}
