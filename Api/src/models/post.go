package models

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// --- Definición del Enum PostType mejorada ---
type PostType int

const (
	Public PostType = iota
	Psi
)

var postTypeToString = map[PostType]string{
	Public: "public",
	Psi:    "psi",
}
var StringToPostType = map[string]PostType{
	"public": Public,
	"psi":    Psi,
}

func (pt PostType) String() string { return postTypeToString[pt] }

func (pt PostType) MarshalJSON() ([]byte, error) { return json.Marshal(pt.String()) }

func (pt *PostType) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	var ok bool
	*pt, ok = StringToPostType[s]
	if !ok {
		return fmt.Errorf("invalid PostType %q", s)
	}
	return nil
}

func (pt *PostType) Scan(value interface{}) error {
	s, ok := value.(string)
	if !ok {
		return fmt.Errorf("invalid type for PostType: %T", value)
	}
	*pt, ok = StringToPostType[s]
	if !ok {
		return fmt.Errorf("invalid PostType value %q", s)
	}
	return nil
}

func (pt PostType) Value() (driver.Value, error) { return pt.String(), nil }

// --- Modelo Post Mejorado ---
type Post struct {
	ID               uint     `gorm:"primaryKey" json:"id"`
	Name             string   `gorm:"size:50;not null;uniqueIndex" json:"name"` // Añadido uniqueIndex para evitar títulos duplicados
	Type             PostType `gorm:"type:varchar(10);default:public;not null" json:"type"`
	ShortDescription string   `gorm:"size:250" json:"short_description"`
	TextID           uint     `json:"text_id,omitempty"`             // no vamos a recuerar el texto directamente, solo si es necesario
	IsActive         bool     `gorm:"default:true" json:"is_active"` // Por defecto, el admin está activo

	// Timestamps y Soft Delete
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	// DeletedAt gorm.DeletedAt `gorm:"index" json:"-"` // Para Soft Deletes

	// Claves foráneas (ocultas en JSON para evitar redundancia)
	CreatedBy *uuid.UUID `gorm:"type:uuid" json:"create_by"`
	UpdatedBy *uuid.UUID `gorm:"type:uuid" json:"update_by"`

	// Relaciones de GORM (expuestas en JSON)
	// 'omitempty' evita que aparezca el campo si es nulo (ej. si un post no tiene updater)
	Creator *UserAdmin `gorm:"foreignKey:CreatedBy;references:ID" json:"creator,omitempty"`
	Updater *UserAdmin `gorm:"foreignKey:UpdatedBy;references:ID" json:"updater,omitempty"`

	// Imagen asociada

}
