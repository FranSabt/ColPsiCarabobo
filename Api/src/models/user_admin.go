package models

type UserAdmin struct {
	AbstractUser `gorm:"embedded"` // Usar embedded para composición

	IsActive bool `gorm:"default:true"` // Por defecto, el admin está activo

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
