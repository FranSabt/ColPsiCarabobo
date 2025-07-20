package models

import "github.com/google/uuid"

type UserAdmin struct {
	ID       uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	Username string    `gorm:"size:25;unique;not null" json:"username"`
	// RECOMENDADO
	Email    string `gorm:"size:255;unique;not null" json:"email"`
	Password string `gorm:"size:512;not null" json:"-"`
	IsActive bool   `gorm:"default:true" json:"is_active"` // Por defecto, el admin está activo
	Key      string `gorm:"size:512;" json:"key"`
	Sudo     bool   `gorm:"default:false" json:"-"`

	// Permisos sobre administradores
	CanCreateAdmin bool `gorm:"default:false" json:"can_create_admin"`
	CanUpdateAdmin bool `gorm:"default:false" json:"can_update_admin"`
	CanDeleteAdmin bool `gorm:"default:false" json:"can_delete_admin"`

	// Permisos sobre publicaciones
	CanPublish       bool `gorm:"default:false" json:"can_publish"`
	CanUpdatePublish bool `gorm:"default:false" json:"can_update_publish"`
	CanDeletePublish bool `gorm:"default:false" json:"can_delete_publish"`

	// Permisos de notificaciones
	CanSendNotifications   bool `gorm:"default:false" json:"can_send_notifications"`   // Puede enviar notificaciones
	CanManageNotifications bool `gorm:"default:false" json:"can_manage_notifications"` // Puede gestionar notificaciones (editar/eliminar)
	CanReadNotifications   bool `gorm:"default:false" json:"can_read_notifications"`   // Puede leer notificaciones

	// Permisos para crear etiquetas
	CanCreateTags bool `gorm:"default:false" json:"can_create_tags"` // Puede crear etiquetas
	CanEditTags   bool `gorm:"default:false" json:"can_edit_tags"`   // Puede editar etiquetas
	CanDeleteTags bool `gorm:"default:false" json:"can_delete_tags"` // Puede eliminar etiquetas
}
