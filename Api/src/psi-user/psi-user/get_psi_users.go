package psiuser_presenter

import (
	"net/http"

	psi_user_db "github.com/FranSabt/ColPsiCarabobo/src/psi-user/db"
	psi_user_mapper "github.com/FranSabt/ColPsiCarabobo/src/psi-user/mapper"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
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

//////////////////////////////////////////////////////
//////////////////////////////////////////////////////
//////////////////////////////////////////////////////

type RequestPsiUserById struct {
	ID string `json:"id"`
}

func GetPsiUserById(c *fiber.Ctx, db *gorm.DB) error {
	var request RequestPsiUserById

	if err := c.BodyParser(&request); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Cuerpo de solicitud inválido",
			"details": err.Error(),
		})
	}

	id, err := uuid.Parse(request.ID)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "invalid id format",
			"details": err.Error(),
		})
	}

	// Bucar el usuario y info del colegio
	psi_user, psi_user_col_data, err := psi_user_db.GetPsiUserByIdDetails(db, id)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Error while trying to retrieve the PsiUser",
			"error":   err.Error(),
		})
	}

	if psi_user == nil || psi_user_col_data == nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Error while trying to retrieve the PsiUser",
			// "error":   err.Error(),
		})
	}

	psi_user_public := psi_user_mapper.PsiUserDataToPublic(psi_user, psi_user_col_data)
	if psi_user_public == nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Error while parsing data",
			// "error":   err.Error(),
		})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Psichologyst found",
		"data":    psi_user_public,
	})
}
