package admin

import (
	"errors"

	db_admin "github.com/FranSabt/ColPsiCarabobo/src/admin/db"
	"github.com/FranSabt/ColPsiCarabobo/src/models"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func GetAdmins(c *fiber.Ctx, db *gorm.DB) error {
	// Parámetros de paginación
	page := c.QueryInt("page", 1)
	pageSize := c.QueryInt("pageSize", 10)

	// Parámetros de búsqueda de texto
	username := c.Query("username")
	email := c.Query("email")

	// Parámetro de búsqueda booleano
	var isActive *bool
	if c.Query("isActive") != "" {
		isActiveVal := c.QueryBool("isActive") // Devuelve true para "true", "1", etc. y false para el resto
		isActive = &isActiveVal
	}

	// Llamar a la función de la base de datos
	admins, totalRecords, err := db_admin.GetPaginatedAdmins(db, page, pageSize, username, email, isActive)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Error al obtener los administradores",
			"details": err.Error(),
		})
	}

	// Devolver la respuesta
	return c.JSON(fiber.Map{
		"data":         admins,
		"totalRecords": totalRecords,
		"page":         page,
		"pageSize":     pageSize,
	})
}

//////////////////////////////////////////////////////
//////////////////////////////////////////////////////
//////////////////////////////////////////////////////

type adminCreateOrUpdateRequest struct {
	ID string `json:"admin_id"`

	//////
	Username   string `json:"username"`
	Email      string `json:"email"`
	NewAdminId string `json:"new_admin_id"`
	Password   string `json:"password"`

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

func CreateOrUpdateAdminHandler(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// --- 1. AUTENTICACIÓN: Obtener ID del admin que realiza la petición (del token) ---
		var request adminCreateOrUpdateRequest

		if err := c.BodyParser(&request); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"success": false,
				"error":   "Invalid request format",
			})
		}
		uuid_parsed, err := uuid.Parse(request.ID)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"success": false,
				"error":   "not valid id",
			})
		}

		requesterAdmin, err := db_admin.GetAdminById(uuid_parsed, db)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"success": false,
				"error":   "not valid id",
			})
		}

		// --- 3. VALIDACIÓN DEL BODY ---
		var payload adminCreateOrUpdateRequest
		if err := c.BodyParser(&payload); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Formato de petición inválido: " + err.Error()})
		}

		// --- 4. DISTINGUIR ENTRE CREACIÓN Y ACTUALIZACIÓN ---
		isUpdateRequest := payload.NewAdminId != ""
		var targetAdmin *models.UserAdmin // El admin que será creado o actualizado

		// --- 5. AUTORIZACIÓN Y LÓGICA DE ACTUALIZACIÓN ---
		if isUpdateRequest {

			uuid_parsed_to_update, err := uuid.Parse(request.ID)
			if err != nil {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"success": false,
					"error":   "not valid id",
				})
			}
			if !requesterAdmin.CanUpdateAdmin {
				return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "No tienes permisos para actualizar administradores."})
			}
			// Si es actualización, primero obtenemos el admin existente de la BD.
			targetAdmin, err = db_admin.GetAdminById(uuid_parsed_to_update, db)
			if err != nil {
				return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "El administrador a actualizar no fue encontrado."})
			}
		} else { // --- 6. AUTORIZACIÓN Y LÓGICA DE CREACIÓN ---
			if !requesterAdmin.CanCreateAdmin {
				return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "No tienes permisos para crear administradores."})
			}
			if payload.Password == "" {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "La contraseña es obligatoria al crear un nuevo administrador."})
			}
		}

		// --- 7. VERIFICACIÓN DE ESCALADA DE PRIVILEGIOS (EL PASO CLAVE) ---
		if err := checkPermissionEscalation(*requesterAdmin, payload); err != nil {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Intento de escalada de privilegios denegado: " + err.Error()})
		}

		// --- 8. MAPEAR DATOS Y PREPARAR PARA GUARDAR ---
		targetAdmin.Username = payload.Username
		targetAdmin.Email = payload.Email
		// Asignar todos los permisos desde el payload
		targetAdmin.CanCreateAdmin = payload.CanCreateAdmin
		targetAdmin.CanUpdateAdmin = payload.CanUpdateAdmin
		// ... mapea el resto de los permisos ...

		// Hashear la contraseña si se proporcionó una nueva
		if payload.Password != "" {
			hashedPassword, err := bcrypt.GenerateFromPassword([]byte(payload.Password), bcrypt.DefaultCost)
			if err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error al procesar la contraseña."})
			}
			targetAdmin.Password = string(hashedPassword)
		}

		// --- 9. PERSISTENCIA ---
		err = db_admin.CreateOrUpdateAdmin(*targetAdmin, db) // Usamos la función que devuelve el objeto y un error
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "No se pudo guardar el administrador: " + err.Error()})
		}

		if isUpdateRequest {
			return c.Status(fiber.StatusOK).JSON(fiber.Map{
				"success": true,
				"error":   "not valid id",
			})

		}

		return c.Status(fiber.StatusCreated).JSON(fiber.Map{
			"success": true,
			"error":   "not valid id",
		})

	}
}

func checkPermissionEscalation(requester models.UserAdmin, payload adminCreateOrUpdateRequest) error {

	// --- Permisos sobre Administradores ---
	if payload.CanCreateAdmin && !requester.CanCreateAdmin {
		return errors.New("no puedes conceder el permiso 'CanCreateAdmin' porque no lo posees")
	}
	if payload.CanUpdateAdmin && !requester.CanUpdateAdmin {
		return errors.New("no puedes conceder el permiso 'CanUpdateAdmin' porque no lo posees")
	}
	if payload.CanDeleteAdmin && !requester.CanDeleteAdmin {
		return errors.New("no puedes conceder el permiso 'CanDeleteAdmin' porque no lo posees")
	}

	// --- Permisos sobre Publicaciones ---
	if payload.CanPublish && !requester.CanPublish {
		return errors.New("no puedes conceder el permiso 'CanPublish' porque no lo posees")
	}
	if payload.CanUpdatePublish && !requester.CanUpdatePublish {
		return errors.New("no puedes conceder el permiso 'CanUpdatePublish' porque no lo posees")
	}
	if payload.CanDeletePublish && !requester.CanDeletePublish {
		return errors.New("no puedes conceder el permiso 'CanDeletePublish' porque no lo posees")
	}

	// --- Permisos de Notificaciones ---
	if payload.CanSendNotifications && !requester.CanSendNotifications {
		return errors.New("no puedes conceder el permiso 'CanSendNotifications' porque no lo posees")
	}
	if payload.CanManageNotifications && !requester.CanManageNotifications {
		return errors.New("no puedes conceder el permiso 'CanManageNotifications' porque no lo posees")
	}
	if payload.CanReadNotifications && !requester.CanReadNotifications {
		return errors.New("no puedes conceder el permiso 'CanReadNotifications' porque no lo posees")
	}

	// --- Permisos para Etiquetas (Tags) ---
	if payload.CanCreateTags && !requester.CanCreateTags {
		return errors.New("no puedes conceder el permiso 'CanCreateTags' porque no lo posees")
	}
	if payload.CanEditTags && !requester.CanEditTags {
		return errors.New("no puedes conceder el permiso 'CanEditTags' porque no lo posees")
	}
	if payload.CanDeleteTags && !requester.CanDeleteTags {
		return errors.New("no puedes conceder el permiso 'CanDeleteTags' porque no lo posees")
	}

	// Si se llega hasta aquí, no se ha violado ninguna regla.
	return nil
}
