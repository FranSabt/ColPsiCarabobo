package db_admin

import (
	"errors"
	"fmt"
	"log"
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

func GetAdminById(id uuid.UUID, db *gorm.DB) (*models.UserAdmin, error) {
	var admin models.UserAdmin

	// 1. Ejecutar la consulta. `db.First` devuelve un objeto *gorm.DB.
	// El error, si existe, se encuentra en el campo `result.Error`.
	result := db.First(&admin, "id = ?", id)

	// 2. Comprobar si hubo un error en la consulta.
	if result.Error != nil {
		// Es una buena práctica verificar si el error específico es que no se encontró el registro.
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			// Devolvemos una estructura vacía y el error específico para que quien llame a la función sepa qué pasó.
			return nil, gorm.ErrRecordNotFound
		}
		// Para cualquier otro tipo de error de base de datos (conexión, sintaxis, etc.).
		return nil, result.Error
	}

	// 3. Si no hubo errores, la variable 'admin' ha sido poblada con los datos.
	// Devolvemos el admin encontrado y un error nulo para indicar que todo fue exitoso.
	return &admin, nil
}

func CreateOrUpdateAdmin(admin_model models.UserAdmin, db *gorm.DB) error {
	// - Si admin_model.ID es el valor cero (uuid.Nil), crea un nuevo registro.
	// - Si admin_model.ID tiene un valor, actualiza el registro con ese ID.
	result := db.Save(&admin_model)

	// Es una buena práctica siempre verificar si ocurrió un error.
	if result.Error != nil {
		// En una aplicación real, manejarías este error de forma más elegante.
		// Por ahora, solo lo imprimimos en la consola.
		log.Printf("Error al guardar el administrador: %v", result.Error)
		return result.Error
	}

	// Si todo fue bien, el registro fue creado o actualizado.
	// El `admin_model` ahora tendrá el ID asignado por la base de datos si fue una creación.
	log.Printf("Administrador guardado exitosamente con ID: %s", admin_model.ID)
	log.Printf("Total de filas afectadas: %d", result.RowsAffected)
	return nil
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

func GetAdminByUsernameOrEmal(db *gorm.DB, username, query string) (*models.UserAdmin, error) {
	var admin models.UserAdmin
	err := db.Where(query, username).First(&admin).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &admin, nil
}

func SaveUpdatedAdminOnly(db *gorm.DB, admin *models.UserAdmin) error {
	// Iniciar transacción
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Guardar admin
	if err := tx.Save(admin).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Commit si todo fue bien
	return tx.Commit().Error
}
