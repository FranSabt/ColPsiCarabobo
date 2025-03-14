package psi_user_admin_presenter

import (
	psi_use_controller "github.com/FranSabt/ColPsiCarabobo/src/psi-user/controller"
	psi_user_db "github.com/FranSabt/ColPsiCarabobo/src/psi-user/db"
	psi_user_mapper "github.com/FranSabt/ColPsiCarabobo/src/psi-user/mapper"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func UploadCsv(c *fiber.Ctx, db *gorm.DB) error {
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

	count := 0
	fail := []struct {
		Record interface{} // Registro original
		Error  string      // Mensaje de error
	}{}

	//TODO: levar a una funcion aparte y generar la entrada con hilos idependientes
	// Procesar cada registro del CSV
	for _, r := range *result {
		// Mapear los datos a los modelos
		psi_model_mapped := psi_user_mapper.PsiUserCsv_To_PsiUserModel(r)

		// Intentar guardar el PsiUserModel
		err = psi_user_db.CreatePsiUseDb(db, psi_model_mapped)
		if err != nil {
			// Si falla, añadir el registro fallido a la lista de fallos
			fail = append(fail, struct {
				Record interface{}
				Error  string
			}{
				Record: r,
				Error:  "Error al crear PsiUserModel: " + err.Error(),
			})
			continue // Saltar al siguiente registro
		}

		// se reliza este proceso despues de que el registro del psiuser es exitoso
		psi_col_data_mapped := psi_user_mapper.PsiUserCsv_To_PsiUserColData(r)
		// Asignar el ID del PsiUserModel al PsiUserColData
		psi_col_data_mapped.PsiUserModelID = psi_model_mapped.ID

		// Intentar guardar el PsiUserColData
		err = psi_user_db.CreatePsiColDataDb(db, psi_col_data_mapped)
		if err != nil {
			// Si falla, añadir el registro fallido a la lista de fallos
			fail = append(fail, struct {
				Record interface{}
				Error  string
			}{
				Record: r,
				Error:  "Error al crear PsiUserColData: " + err.Error(),
			})
			continue // Saltar al siguiente registro
		}

		count++ // Incrementar el contador de registros exitosos
	}

	// Devolver una respuesta exitosa
	return c.JSON(fiber.Map{
		"message":                    "CSV procesado correctamente",
		"success_registres":          count,
		"number_of_failed_registres": len(fail),
		"failed_registres":           fail,
	})
}
