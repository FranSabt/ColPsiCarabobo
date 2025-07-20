package specialties

import (
	"errors"
	"fmt"
	"net/http"
	"regexp"

	"github.com/FranSabt/ColPsiCarabobo/src/admin/admin"
	specialties_controller "github.com/FranSabt/ColPsiCarabobo/src/specialties/controller"
	specialties_mapper "github.com/FranSabt/ColPsiCarabobo/src/specialties/mapper"
	specialties_structs "github.com/FranSabt/ColPsiCarabobo/src/specialties/request-structs"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

var unicodeLetterNumberRegex = regexp.MustCompile(`^[\p{L}\p{N}]+$`)

func CreatePsiSpecialty(c *fiber.Ctx, db *gorm.DB) error {
	var request specialties_structs.SpecialtiesRequest

	if err := c.BodyParser(&request); err != nil {
		return c.JSON(fiber.Map{
			"success": false,
			"error":   "Error while parsing the request",
			"details": err.Error(),
		})
	}

	fmt.Println(request)

	admin_id, err := uuid.Parse(request.AdmindId)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "invalid id format",
			"details": err.Error(),
		})
	}

	admin_exist, err := admin.AdminExists(admin_id, db)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "invalid id format",
			"details": err.Error(),
		})
	}

	if !admin_exist {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "invalid id format",
			"details": "admin does not exist",
		})
	}

	// TODO: bindear el id
	// fmt.Println(id)

	if len(request.Name) < 4 {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Specialty name must be at least 4 characters long.",
		})
	}

	// 2. Validación de caracteres especiales (usando la regex compilada)
	if !unicodeLetterNumberRegex.MatchString(request.Name) {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Specialty name must contain only letters and numbers.",
		})
	}

	if len(request.Description) < 150 {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Specialty description must be at least 150 characters long.",
		})
	}

	specialty_model := specialties_mapper.SpecialtyRequestToSpecialtyModel(request)

	err = specialties_controller.SaveNewSpecialty(db, specialty_model)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   err,
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    nil,
	})
}

//////////////////////////////////////////////////////
//////////////////////////////////////////////////////
//////////////////////////////////////////////////////

func GetPsiSpecialtiesCount(c *fiber.Ctx, db *gorm.DB) error {
	count, last_id, err := specialties_controller.GetPsiSpecialtiesCountController(db)
	if err != nil {
		// Si hay un error de base de datos, es un problema del servidor.
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Failed to retrieve specialties count",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"total_count": count,
			"last_id":     last_id,
		},
	})
}

//////////////////////////////////////////////////////
//////////////////////////////////////////////////////
//////////////////////////////////////////////////////

func GetPsiSpecialtiesNames(c *fiber.Ctx, db *gorm.DB) error {
	specialties, err := specialties_controller.GetPsiSpecialtiesNamesController(db)
	if err != nil {
		// Error del servidor al consultar la base de datos.
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   "Failed to retrieve specialties names",
		})
	}

	// Si no hay especialidades, GORM devuelve un slice vacío, no un error.
	// Esto es correcto y el cliente recibirá un array vacío `[]`.
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data":    specialties,
	})
}

//////////////////////////////////////////////////////
//////////////////////////////////////////////////////
//////////////////////////////////////////////////////

func GetPsiSpecialtiesDescription(c *fiber.Ctx, db *gorm.DB) error {
	// 1. Parsear el ID desde el parámetro de la URL (ej: /specialties/123/description)
	id, err := c.ParamsInt("id")
	if err != nil {
		// Si el ID no es un número, es una petición mal formada.
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Invalid ID format. Must be an integer.",
		})
	}

	// 2. Llamar al controlador para obtener los datos.
	description, err := specialties_controller.GetPsiSpecialtiesDescriptionController(db, uint(id))

	// 3. Manejar los errores de forma específica.
	if err != nil {
		// Si el error es específicamente 'registro no encontrado'...
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// ...devolvemos un 404 Not Found, que es el código correcto.
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"success": false,
				"error":   err.Error(),
			})
		}
		// Para cualquier otro error de base de datos, es un 500.
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   "Database error while fetching description",
		})
	}

	// 4. Devolver la respuesta exitosa.
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"description": description,
		},
	})
}

//////////////////////////////////////////////////////
//////////////////////////////////////////////////////
//////////////////////////////////////////////////////

func UpdatePsiSepecialty(c *fiber.Ctx, db *gorm.DB) error {
	var request specialties_structs.SpecialtyUpdate

	if err := c.BodyParser(&request); err != nil {
		return c.JSON(fiber.Map{
			"success": false,
			"error":   "Internal server error",
		})
	}

	if request.ID <= 0 {
		return c.JSON(fiber.Map{
			"success": false,
			"error":   "No Id to update",
		})
	}

	// Verificamos que haya algo que actualizar
	if request.Description == "" && request.Name == "" {
		return c.JSON(fiber.Map{
			"success": false,
			"error":   "Nothing to update",
		})
	}

	admin_id, err := uuid.Parse(request.AdmindId)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "invalid id format",
			"details": err.Error(),
		})
	}

	admin_exist, err := admin.AdminExists(admin_id, db)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "invalid id format",
			"details": err.Error(),
		})
	}

	if !admin_exist {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "invalid id format",
			"details": "admin does not exist",
		})
	}

	// // TODO: bindear el id
	// fmt.Println(admin_id)

	err = specialties_controller.UpdatePsiSpecialtyController(&request, db)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "error updating",
			"details": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data":    nil,
	})
}

//////////////////////////////////////////////////////
//////////////////////////////////////////////////////
//////////////////////////////////////////////////////

type deleteSpecialty struct {
	AdminId     string `json:"admin_id"`
	SpecialtyId int64  `json:"specialty_id"`
}

// TODO Eliminar el query
func DeletePsiSpecialty(c *fiber.Ctx, db *gorm.DB) error {
	var request deleteSpecialty

	if err := c.BodyParser(&request); err != nil {
		return c.JSON(fiber.Map{
			"success": false,
			"error":   "Internal server error",
		})
	}

	admin_id := request.AdminId
	specialty_id := request.SpecialtyId

	admin_uuid, err := uuid.Parse(admin_id)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "invalid id format",
			"details": err.Error(),
		})
	}
	if specialty_id <= 0 {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "invalid id format",
			"details": "cant use a zero or less value id",
		})
	}

	admin_exist, err := admin.AdminExists(admin_uuid, db)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "invalid id format",
			"details": err.Error(),
		})
	}

	if !admin_exist {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "invalid id format",
			"details": "admin does not exist",
		})
	}

	// TODO: verificar el administrador
	// fmt.Println(admin_uuid)

	err = specialties_controller.DeleteSpecialtyController(int64(specialty_id), db)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "error updating",
			"details": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data":    nil,
	})
}

//////////////////////////////////////////////////////
//////////////////////////////////////////////////////
//////////////////////////////////////////////////////

// ------       Auxiliary Functions      ------     //

//////////////////////////////////////////////////////
//////////////////////////////////////////////////////
//////////////////////////////////////////////////////
