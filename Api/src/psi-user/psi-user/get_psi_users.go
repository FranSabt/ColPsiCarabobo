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
	// --- PARÁMETROS DE PAGINACIÓN ---
	page := c.QueryInt("page", 1)
	pageSize := c.QueryInt("pageSize", 10)

	// --- PARÁMETROS DE BÚSQUEDA ---

	// Parámetros numéricos (se mantienen como estaban)
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

	// Nuevos parámetros de texto
	name := c.Query("name")
	location := c.Query("location")
	specialty := c.Query("specialty")

	// --- LLAMADA AL SERVICIO DE BASE DE DATOS ---
	// Pasamos todos los parámetros a la función de búsqueda
	psiUsers, totalRecords, err := psi_user_db.GetPaginatedPsiUsers(
		db,
		page,
		pageSize,
		ci,
		fpv,
		name,
		location,
		specialty,
	)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Error al obtener los psicólogos",
			"details": err.Error(),
		})
	}

	// --- RESPUESTA ---
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
