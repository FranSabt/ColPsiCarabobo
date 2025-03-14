package models

import (
	"time"

	"github.com/google/uuid"
)

type PsiUserModel struct {
	ID       uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	Username string    `gorm:"size:25;not null"`
	Email    string    `gorm:"size:50;not null"`
	Password string    `gorm:"size:512;not null"`

	// Identity
	FirstName      string    `gorm:"size:255;not null"`
	SecondName     string    `gorm:"size:255"`
	LastName       string    `gorm:"size:255;not null"`
	SecondLastName string    `gorm:"size:255"`
	FPV            int       `gorm:"not null"`
	CI             int       `gorm:"not null;unique"`
	Nationality    string    `gorm:"size:1;not null"`
	BornDate       time.Time `gorm:"type:date;not null"`
	Genre          string    `gorm:"size:1;not null"`

	// Contact
	ContactEmail             string `gorm:"size:255;not null"`
	ShowContactEmail         bool   `gorm:"default:false"`
	PublicPhone              string `gorm:"size:20"`
	ShowPublicPhone          bool   `gorm:"default:false"`
	ServiceAddress           string `gorm:"size:255"`
	ShowPublicServiceAddress bool   `gorm:"default:false"`

	// PsiCol
	Solvent     bool `gorm:"default:false"`
	ProofOfLife bool `gorm:"default:false"`

	// Carabobo Direction
	MunicipalityCarabobo string `gorm:"size:255"`
	PhoneCarabobo        string `gorm:"size:20"`
	CelPhoneCarabobo     string `gorm:"size:20"`

	// Outside Carabobo Direction
	StateOutside                string `gorm:"size:255"`
	MunicipalityOutSideCarabobo string `gorm:"size:255"`
	PhoneOutSideCarabobo        string `gorm:"size:20"`
	CelPhoneOutSideCarabobo     string `gorm:"size:20"`

	// Relación con PsiUserColData
	PsiUserColDataID *uuid.UUID `gorm:"type:uuid"` // Clave foránea

}
