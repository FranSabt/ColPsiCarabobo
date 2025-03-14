package psiuser_presenter

import (
	psi_user_db "github.com/FranSabt/ColPsiCarabobo/src/psi-user/db"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func GetPsiUsers(c *fiber.Ctx, db *gorm.DB) error {
	// Obtener los parámetros de paginación de la consulta
	page := c.QueryInt("page", 1)          // Página actual (por defecto 1)
	pageSize := c.QueryInt("pageSize", 10) // Tamaño de la página (por defecto 10)

	// Obtener los parámetros de búsqueda (opcionales)
	var ci *int
	if c.Query("ci") != "" {
		ciValue := c.QueryInt("ci")
		ci = &ciValue
	}

	var fpv *int
	if c.Query("fpv") != "" {
		fpvValue := c.QueryInt("fpv")
		fpv = &fpvValue
	}

	// Obtener los registros paginados con los filtros aplicados
	psiUsers, totalRecords, err := psi_user_db.GetPaginatedPsiUsers(db, page, pageSize, ci, fpv)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Error al obtener los registros",
			"details": err.Error(),
		})
	}

	// Devolver la respuesta con los registros y la información de paginación
	return c.JSON(fiber.Map{
		"data":         psiUsers,
		"totalRecords": totalRecords,
		"page":         page,
		"pageSize":     pageSize,
	})
}
