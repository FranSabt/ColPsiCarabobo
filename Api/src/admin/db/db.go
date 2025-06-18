package db_admin

import (
	"fmt"
	"strings"

	"github.com/FranSabt/ColPsiCarabobo/src/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// AdminUserResponse es un DTO para devolver datos de administrador de forma segura (sin contraseña).
type AdminUserResponse struct {
	ID       uuid.UUID `json:"id"`
	Username string    `json:"username"`
	Email    string    `json:"email"`
	IsActive bool      `json:"is_active"`

	// Permisos
	CanCreateAdmin         bool `json:"can_create_admin"`
	CanUpdateAdmin         bool `json:"can_update_admin"`
	CanDeleteAdmin         bool `json:"can_delete_admin"`
	CanPublish             bool `json:"can_publish"`
	CanUpdatePublish       bool `json:"can_update_publish"`
	CanDeletePublish       bool `json:"can_delete_publish"`
	CanSendNotifications   bool `json:"can_send_notifications"`
	CanManageNotifications bool `json:"can_manage_notifications"`
	CanReadNotifications   bool `json:"can_read_notifications"`
	CanCreateTags          bool `json:"can_create_tags"`
	CanEditTags            bool `json:"can_edit_tags"`
	CanDeleteTags          bool `json:"can_delete_tags"`
}

func GetPaginatedAdmins(db *gorm.DB, page, pageSize int, username, email string, isActive *bool) ([]AdminUserResponse, int64, error) {
	var admins []AdminUserResponse
	var totalRecords int64

	offset := (page - 1) * pageSize

	// Crear la consulta base apuntando al modelo UserAdmin
	query := db.Model(&models.UserAdmin{})

	// --- APLICAR FILTROS (MÉTODO GORM IDIOMÁTICO) ---

	// Filtro por Username.
	// Para búsquedas insensibles a mayúsculas de forma portable entre bases de datos (MySQL, Postgres, etc.)
	// se usa la función LOWER(). GORM no tiene un método .ILike() nativo.
	// La forma correcta de aplicar esta lógica es todavía con una cláusula Where.
	if username != "" {
		// Construimos el patrón de búsqueda.
		searchUsername := fmt.Sprintf("%%%s%%", strings.ToLower(username))
		// Aplicamos la función LOWER() a la columna para hacer la comparación.
		query = query.Where("LOWER(username) LIKE ?", searchUsername)
	}

	// Filtro por Email, usando la misma técnica.
	if email != "" {
		searchEmail := fmt.Sprintf("%%%s%%", strings.ToLower(email))
		query = query.Where("LOWER(email) LIKE ?", searchEmail)
	}

	// Filtro por estado (activo/inactivo).
	// ¡Aquí es donde el método de GORM sin cadenas de SQL brilla!
	if isActive != nil {
		// Para igualdades simples, podemos pasar un struct o un map.
		// GORM lo traducirá a: WHERE `is_active` = 'true' (o false).
		// Esto es más seguro a nivel de tipos y más limpio.
		query = query.Where(&models.UserAdmin{IsActive: *isActive})
	}

	// Contar el total de registros que coinciden con los filtros.
	// La lógica de Count() no cambia.
	if err := query.Count(&totalRecords).Error; err != nil {
		return nil, 0, err
	}

	// Aplicar paginación y ejecutar la consulta.
	// La lógica de paginación y Find() no cambia.
	if err := query.Offset(offset).Limit(pageSize).Find(&admins).Error; err != nil {
		return nil, 0, err
	}

	return admins, totalRecords, nil
}
