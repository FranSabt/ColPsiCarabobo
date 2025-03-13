package psi_user_admin_presenter

import (
	psi_use_controller "github.com/FranSabt/ColPsiCarabobo/src/psi-user/controller"
	"github.com/gofiber/fiber/v2"
)

func UploadCsv(c *fiber.Ctx) error {
	// Obtener el archivo CSV de la solicitud
	file, err := c.FormFile("csv")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "No se pudo obtener el archivo CSV",
			"details": err.Error(),
		})
	}

	// Abrir el archivo CSV
	src, err := file.Open()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "No se pudo abrir el archivo CSV",
			"details": err.Error(),
		})
	}
	defer src.Close()

	// Procesar el archivo CSV
	result, err := psi_use_controller.ProcessCsv(src)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "No se pudo procesar el CSV",
			"details": err.Error(),
		})
	}

	// Devolver una respuesta exitosa
	return c.JSON(fiber.Map{
		"message": "CSV procesado correctamente",
		"data":    result,
	})
}
