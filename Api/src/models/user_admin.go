package models

import "github.com/google/uuid"

type UserAdmin struct {
	ID       uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	Username string    `gorm:"size:25;not null"`
	Email    string    `gorm:"size:50;not null"`
	Password string    `gorm:"size:512;not null"`
	IsActive bool      `gorm:"default:true"` // Por defecto, el admin está activo

	// Permisos sobre administradores
	CanCreateAdmin bool `gorm:"default:false"`
	CanUpdateAdmin bool `gorm:"default:false"`
	CanDeleteAdmin bool `gorm:"default:false"`

	// Permisos sobre publicaciones
	CanPublish       bool `gorm:"default:false"`
	CanUpdatePublish bool `gorm:"default:false"`
	CanDeletePublish bool `gorm:"default:false"`

	// Permisos de notificaciones
	CanSendNotifications   bool `gorm:"default:false"` // Puede enviar notificaciones
	CanManageNotifications bool `gorm:"default:false"` // Puede gestionar notificaciones (editar/eliminar)
	CanReadNotifications   bool `gorm:"default:false"` // Puede leer notificaciones

	// Permisos para crear etiquetas
	CanCreateTags bool `gorm:"default:false"` // Puede crear etiquetas
	CanEditTags   bool `gorm:"default:false"` // Puede editar etiquetas
	CanDeleteTags bool `gorm:"default:false"` // Puede eliminar etiquetas
}
