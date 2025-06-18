package admin

import (
	db_admin "github.com/FranSabt/ColPsiCarabobo/src/admin/db"
	"github.com/gofiber/fiber/v2"
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
