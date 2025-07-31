package psiuser_presenter

import (
	"errors"
	"fmt"
	"mime/multipart"
	"net/http"
	"strings"
	"time"

	"github.com/FranSabt/ColPsiCarabobo/db"
	image_controller "github.com/FranSabt/ColPsiCarabobo/src/images/controller"
	db_images "github.com/FranSabt/ColPsiCarabobo/src/images/db"
	"github.com/FranSabt/ColPsiCarabobo/src/models"
	psi_user_controller "github.com/FranSabt/ColPsiCarabobo/src/psi-user/controller"
	psi_user_db "github.com/FranSabt/ColPsiCarabobo/src/psi-user/db"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Definimos una struct para mapear los datos del formulario.
// Las etiquetas `form:"..."` le dicen a Fiber qué campo del formulario corresponde a cada campo de la struct.
type create_postgrade_request struct {
	PsiUserID               string `form:"psi_user_id"`
	PostGradeTitle          string `form:"post_grade_title"`
	PostGradeUniversity     string `form:"post_grade_university"`
	PostGradeGraduationYear string `form:"post_grade_graduation_year"`
	PostGradeDescription    string `form:"post_grade_description"`
}

// send_error_response centraliza el formato de las respuestas de error.
func send_error_response(c *fiber.Ctx, status_code int, err_msg string, details interface{}) error {
	return c.Status(status_code).JSON(fiber.Map{
		"success": false,
		"error":   err_msg,
		"details": details,
	})
}

//////////////////////////////////////////////////////
//////////////////////////////////////////////////////
//////////////////////////////////////////////////////

// En el paquete psiuser_presenter

// GetPsiUserPostgradeById busca y devuelve un postgrado específico por su ID.
// Es una ruta pública.
func GetPsiUserPostgradeById(c *fiber.Ctx, db *gorm.DB) error {
	// 1. Obtener y validar el ID de los parámetros de la ruta
	postgrade_id_str := c.Params("id")
	postgrade_id, err := uuid.Parse(postgrade_id_str)
	if err != nil {
		return send_error_response(c, http.StatusBadRequest, "invalid postgrade id format", err.Error())
	}

	// 2. Buscar el registro en la base de datos
	var postgrade models.PisUserPostGrade
	// Usamos .First() que automáticamente añade "WHERE deleted_at IS NULL"
	// para no mostrar registros con soft delete.
	result := db.First(&postgrade, "id = ?", postgrade_id)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return send_error_response(c, http.StatusNotFound, "postgrade record not found", nil)
		}
		return send_error_response(c, http.StatusInternalServerError, "database error", result.Error.Error())
	}

	// 3. Devolver el resultado
	return c.Status(http.StatusOK).JSON(fiber.Map{
		"success": true,
		"data":    postgrade,
	})
}

//////////////////////////////////////////////////////
//////////////////////////////////////////////////////
//////////////////////////////////////////////////////

func GetPostgradeImageById(c *fiber.Ctx, db db.StructDb) error {
	// 1. Obtener y validar el ID de la imagen
	image_id_str := c.Params("id")
	image_id, err := uuid.Parse(image_id_str)
	if err != nil {
		// No usamos send_error_response aquí porque no queremos devolver JSON en caso de error
		return c.Status(http.StatusBadRequest).SendString("Invalid image ID format")
	}

	// 2. Buscar el registro de la imagen en la base de datos de imágenes
	var image_record models.PostGradePic
	result := db.Image.First(&image_record, "id = ?", image_id)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return c.Status(http.StatusNotFound).SendString("Image not found")
		}
		return c.Status(http.StatusInternalServerError).SendString("Database error")
	}

	// 3. Establecer el Content-Type correcto según el formato de la imagen
	content_type := "application/octet-stream" // Tipo por defecto si no se reconoce
	switch image_record.Format {
	case "jpg", "jpeg":
		content_type = "image/jpeg"
	case "png":
		content_type = "image/png"
		// Añade más casos si soportas otros formatos como webp, gif, etc.
	}
	c.Set("Content-Type", content_type)

	// 4. Enviar los datos crudos de la imagen
	return c.Send(image_record.ImageData)
}

// ////////////////////////////////////////////////////
// ////////////////////////////////////////////////////
// ////////////////////////////////////////////////////
func CreatePsiUserPotgradeRefactored(c *fiber.Ctx, db db.StructDb) error {
	// --- 1. Recopilación y Validación de Datos del Formulario ---
	var req create_postgrade_request
	if err := c.BodyParser(&req); err != nil {
		return send_error_response(c, http.StatusBadRequest, "invalid form data", err.Error())
	}

	// Validación de campos obligatorios
	if req.PsiUserID == "" || req.PostGradeTitle == "" || req.PostGradeUniversity == "" || req.PostGradeGraduationYear == "" {
		return send_error_response(c, http.StatusBadRequest, "empty field on the request", nil)
	}

	// ---- 2. Validación de la Fecha y el Usuario ---
	post_grade_graduation_year_clean, err := formatearFecha(req.PostGradeGraduationYear)
	if err != nil {
		return send_error_response(c, http.StatusBadRequest, "invalid date format", err.Error())
	}

	uuid, err := uuid.Parse(req.PsiUserID)
	if err != nil {
		return send_error_response(c, http.StatusBadRequest, "invalid user id format", err.Error())
	}

	user, err := psi_user_db.GetPsiUserById(db.DB, uuid)
	if err != nil {
		// Asumimos que un error aquí es un problema de DB, no que el usuario no existe.
		return send_error_response(c, http.StatusInternalServerError, "could not query user", err.Error())
	}
	if user == nil {
		return send_error_response(c, http.StatusNotFound, "user not found", nil)
	}
	if !user.Solvent {
		return send_error_response(c, http.StatusForbidden, "user is not active", nil)
	}

	// --- 3. Procesamiento de Imágenes (Usando un Bucle) ---
	var processed_images []models.PostGradePic
	var processing_errors []string
	max_images := 3 // Número máximo de imágenes a procesar

	for i := 1; i <= max_images; i++ {
		// Construir dinámicamente los nombres de los campos del formulario
		file_field_name := fmt.Sprintf("image_%d", i)
		file_name_field := fmt.Sprintf("file_name_%d", i)
		mime_type_field := fmt.Sprintf("mime_type_%d", i)

		file_header, err := c.FormFile(file_field_name)
		if err != nil {
			// Si el error es "no such file", es normal, simplemente no se envió esa imagen.
			if err == http.ErrMissingFile {
				continue
			}
			processing_errors = append(processing_errors, fmt.Sprintf("Error getting %s: %v", file_field_name, err))
			continue
		}

		// Obtener metadatos de la imagen
		file_name := c.FormValue(file_name_field)
		mime_type := c.FormValue(mime_type_field)
		if file_name == "" || mime_type == "" {
			processing_errors = append(processing_errors, fmt.Sprintf("Missing file_name or mime_type for %s", file_field_name))
			continue
		}

		// Procesar la imagen
		image_data, err := psi_user_controller.ImageProcesser(file_header)
		if err != nil {
			processing_errors = append(processing_errors, fmt.Sprintf("Error processing %s: %v", file_field_name, err))
			continue
		}

		image_decoded, format, err := psi_user_controller.ImageDecoder(image_data, mime_type)
		if err != nil {
			processing_errors = append(processing_errors, fmt.Sprintf("Error decoding %s: %v", file_field_name, err))
			continue
		}

		// Truncar nombre de archivo si es muy largo
		if len(file_name) > 50 {
			file_name = file_name[:50]
		}

		// Crear el modelo de la imagen
		image_model, err := image_controller.CreatePostGradePicModel(file_name, format, user.Username, user.ID, image_data, image_decoded)
		if err != nil {
			processing_errors = append(processing_errors, fmt.Sprintf("Error creating model for %s: %v", file_field_name, err))
			continue
		}

		// Si todo fue bien, añadirlo a nuestra lista de imágenes procesadas
		processed_images = append(processed_images, image_model)
	}

	// --- 4. Verificación de Requisitos y Errores ---
	// REQUISITO: "un usuario deberia enviar al menos una fotos del modelo a crear"
	if len(processed_images) == 0 {
		return send_error_response(c, http.StatusBadRequest, "at least one image is required", processing_errors)
	}

	// --- 5. Creación y Persistencia en la Base de Datos ---
	postgrade := psi_user_controller.CreatePsiUserPostGradeModel(
		uuid,
		req.PostGradeTitle,
		req.PostGradeUniversity,
		post_grade_graduation_year_clean,
		req.PostGradeDescription,
	)

	// Guardar cada imagen en la DB y asignar su ID al modelo principal
	for i, img_model := range processed_images {
		if err := db_images.SavePostGradePicModel(img_model, db.Image); err != nil {
			// Si falla al guardar una imagen, lo añadimos a los errores pero podríamos continuar
			processing_errors = append(processing_errors, fmt.Sprintf("Error saving image %s to DB: %v", img_model.Name, err))
			continue // No asignamos el ID si no se pudo guardar
		}

		// *** CORRECCIÓN DEL BUG CRÍTICO ***
		// Asignar el ID al campo correcto (PicOne, PicTwo, PicThree)
		switch i {
		case 0:
			postgrade.PicOne = &img_model.ID
		case 1:
			postgrade.PicTwo = &img_model.ID // Asumiendo que tu modelo tiene PicTwo
		case 2:
			postgrade.PicThree = &img_model.ID // Asumiendo que tu modelo tiene PicThree
		}
	}

	// Finalmente, guardar el modelo principal del postgrado
	if err := psi_user_db.SavePostGradeModel(postgrade, db.DB); err != nil {
		return send_error_response(c, http.StatusInternalServerError, "error saving postgrade", err.Error())
	}

	// --- 6. Respuesta Exitosa ---
	return c.Status(http.StatusCreated).JSON(fiber.Map{
		"success":                   true,
		"data":                      postgrade,
		"image_processing_warnings": processing_errors, // Informar sobre errores no críticos
	})
}

//////////////////////////////////////////////////////
//////////////////////////////////////////////////////
//////////////////////////////////////////////////////

// update_postgrade_request define los campos que se pueden actualizar.
// Los campos son opcionales. El backend solo actualizará los que se envíen.
type update_postgrade_request struct {
	PostGradeTitle          string `form:"post_grade_title"`
	PostGradeUniversity     string `form:"post_grade_university"`
	PostGradeGraduationYear string `form:"post_grade_graduation_year"`
	PostGradeDescription    string `form:"post_grade_description"`
	// Para indicar si se quiere eliminar una imagen existente sin reemplazarla.
	DeletePicOne   bool `form:"delete_pic_1"`
	DeletePicTwo   bool `form:"delete_pic_2"`
	DeletePicThree bool `form:"delete_pic_3"`
}

// UpdatePsiUserPostgrade maneja la actualización de un registro de postgrado existente.
// UpdatePsiUserPostgrade maneja la actualización de un registro de postgrado existente.
func UpdatePsiUserPostgrade(c *fiber.Ctx, db db.StructDb) error {
	// --- 1. Obtener y Validar el ID del Postgrado ---
	postgrade_id_str := c.Params("id")
	postgrade_id, err := uuid.Parse(postgrade_id_str)
	if err != nil {
		return send_error_response(c, http.StatusBadRequest, "invalid postgrade id format", err.Error())
	}

	// --- 2. Buscar el Registro Existente en la DB ---
	var postgrade models.PisUserPostGrade
	if err := db.DB.First(&postgrade, "id = ?", postgrade_id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return send_error_response(c, http.StatusNotFound, "postgrade record not found", nil)
		}
		return send_error_response(c, http.StatusInternalServerError, "database error", err.Error())
	}

	// --- (Opcional pero recomendado) Autorización ---
	// ...

	// --- 3. Recopilar Datos del Formulario ---
	var req update_postgrade_request
	if err := c.BodyParser(&req); err != nil {
		return send_error_response(c, http.StatusBadRequest, "invalid form data", err.Error())
	}

	// --- 4. Actualizar Campos de Texto ---
	// (Esta parte queda igual)
	if req.PostGradeTitle != "" {
		postgrade.PostGradeTitle = req.PostGradeTitle
	}
	if req.PostGradeUniversity != "" {
		postgrade.PostGradeUniversity = req.PostGradeUniversity
	}
	if req.PostGradeDescription != "" {
		postgrade.PostGradeDescription = req.PostGradeDescription
	}
	if req.PostGradeGraduationYear != "" {
		formatted_date, err := formatearFecha(req.PostGradeGraduationYear)
		if err != nil {
			return send_error_response(c, http.StatusBadRequest, "invalid date format", err.Error())
		}
		postgrade.PostGradeGraduationYear = formatted_date
	}

	// --- 5. Procesar Imágenes (Lógica Corregida) ---

	// PASO CLAVE 1: Crear un slice para almacenar el estado FINAL de los IDs de las imágenes.
	// Lo inicializamos con los valores actuales.
	final_pic_ids := []*uuid.UUID{postgrade.PicOne, postgrade.PicTwo, postgrade.PicThree}

	delete_flags := []bool{req.DeletePicOne, req.DeletePicTwo, req.DeletePicThree}

	user, _ := psi_user_db.GetPsiUserById(db.DB, postgrade.PsiUserID)
	if user == nil {
		return send_error_response(c, http.StatusInternalServerError, "owner user not found", nil)
	}

	for i := 0; i < 3; i++ {
		slot_index := i + 1 // 1, 2, 3 para los nombres de campo
		file_field_name := fmt.Sprintf("image_%d", slot_index)

		// Caso A: El usuario quiere ELIMINAR la imagen de este slot
		if delete_flags[i] {
			// (Opcional) Borrar la imagen antigua de la tabla PostGradePic
			// if final_pic_ids[i] != nil {
			//     db_images.DeletePostGradePicModel(*final_pic_ids[i], db.Image)
			// }

			// PASO CLAVE 2: Modificamos nuestro slice temporal.
			final_pic_ids[i] = nil
			continue // Pasamos a la siguiente iteración
		}

		// Caso B: El usuario envía un NUEVO archivo para este slot
		file_header, err := c.FormFile(file_field_name)
		if err != nil {
			if err == http.ErrMissingFile {
				continue // No hay archivo nuevo, no hacemos nada, mantenemos la imagen existente.
			}
			return send_error_response(c, http.StatusBadRequest, fmt.Sprintf("error reading %s", file_field_name), err.Error())
		}

		// (Opcional) Borrar la imagen antigua antes de reemplazarla.
		// if final_pic_ids[i] != nil {
		//     db_images.DeletePostGradePicModel(*final_pic_ids[i], db.Image)
		// }

		// Procesamos la nueva imagen
		new_image_model, err := process_and_create_image_model(c, file_header, slot_index, user)
		if err != nil {
			return send_error_response(c, http.StatusInternalServerError, fmt.Sprintf("error processing %s", file_field_name), err.Error())
		}

		// Guardamos la nueva imagen en la DB
		if err := db_images.SavePostGradePicModel(new_image_model, db.Image); err != nil {
			return send_error_response(c, http.StatusInternalServerError, "error saving new image", err.Error())
		}

		// PASO CLAVE 3: Actualizamos nuestro slice temporal con el ID de la nueva imagen.
		final_pic_ids[i] = &new_image_model.ID
	}

	// PASO CLAVE 4: Asignar los valores del slice temporal de vuelta al modelo principal.
	// Esto se hace una sola vez, después de que toda la lógica del bucle ha terminado.
	postgrade.PicOne = final_pic_ids[0]
	postgrade.PicTwo = final_pic_ids[1]
	postgrade.PicThree = final_pic_ids[2]

	// --- 6. Guardar el Modelo Actualizado en la DB ---
	if err := db.DB.Save(&postgrade).Error; err != nil {
		return send_error_response(c, http.StatusInternalServerError, "failed to update postgrade", err.Error())
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"success": true,
		"data":    postgrade,
	})
}

//////////////////////////////////////////////////////
//////////////////////////////////////////////////////
//////////////////////////////////////////////////////

// SoftDeletePsiUserPostgrade realiza un borrado lógico del registro.
// GORM no eliminará la fila, sino que establecerá la fecha en el campo 'deleted_at'.
func SoftDeletePsiUserPostgrade(c *fiber.Ctx, db db.StructDb) error {
	// --- 1. Obtener y Validar el ID ---
	postgrade_id_str := c.Params("id")
	postgrade_id, err := uuid.Parse(postgrade_id_str)
	if err != nil {
		return send_error_response(c, http.StatusBadRequest, "invalid postgrade id format", err.Error())
	}

	// --- (Opcional pero recomendado) Autorización ---
	// Antes de borrar, deberías verificar que el usuario tiene permiso para hacerlo.

	// --- 2. Realizar el Soft Delete ---
	// GORM es lo suficientemente inteligente como para saber que debe hacer un soft delete
	// porque el modelo tiene el campo gorm.DeletedAt.
	// Pasamos un struct vacío con el ID para que GORM sepa qué tabla y qué registro afectar.
	result := db.DB.Delete(&models.PisUserPostGrade{}, "id = ?", postgrade_id)

	if result.Error != nil {
		return send_error_response(c, http.StatusInternalServerError, "database error during deletion", result.Error.Error())
	}

	// Si RowsAffected es 0, significa que no se encontró el registro para eliminar.
	if result.RowsAffected == 0 {
		return send_error_response(c, http.StatusNotFound, "postgrade record not found", nil)
	}

	// Opcional: Si además quieres marcar el campo `Active` como false.
	// Nota: db.Delete() solo afecta a `DeletedAt`. Para cambiar otros campos,
	// necesitarías un `Update` antes o en lugar del `Delete`.
	// Ejemplo: db.DB.Model(&models.PisUserPostGrade{}).Where("id = ?", postgrade_id).Update("active", false)

	// --- 3. Enviar Respuesta ---
	// El código de estado 204 No Content es el estándar para una eliminación exitosa.
	return c.SendStatus(http.StatusNoContent)
}

//////////////////////////////////////////////////////
//////////////////////////////////////////////////////
//////////////////////////////////////////////////////

// ////////////////////////////////////////////////////
// ////////////////////////////////////////////////////
// ////////////////////////////////////////////////////
var mesesEnEspanol = map[time.Month]string{
	time.January:   "enero",
	time.February:  "febrero",
	time.March:     "marzo",
	time.April:     "abril",
	time.May:       "mayo",
	time.June:      "junio",
	time.July:      "julio",
	time.August:    "agosto",
	time.September: "septiembre",
	time.October:   "octubre",
	time.November:  "noviembre",
	time.December:  "diciembre",
}

func formatearFecha(fechaStr string) (string, error) {
	// Limpiar el string de entrada por si acaso
	fechaStr = strings.TrimSpace(fechaStr)

	// Lista de layouts (formatos) a probar, en orden de prioridad.
	// La fecha de referencia en Go es: Mon Jan 2 15:04:05 MST 2006
	// Que corresponde a:           mes día h min seg zona año
	layouts := []string{
		// Formatos comunes con números y separadores
		"2006-01-02", // ISO 8601 (YYYY-MM-DD) - El más recomendado
		"02-01-2006", // Formato común en LATAM/Europa (DD-MM-YYYY)
		"2006/01/02", // YYYY/MM/DD
		"02/01/2006", // DD/MM/YYYY

		// Formatos que incluyen hora (la ignoraremos en la salida)
		"2006-01-02 15:04:05",
		"02-01-2006 15:04:05",
		time.RFC3339, // Formato estándar de internet (ej: "2023-10-27T10:00:00Z")
		time.RFC822,  // Otro formato estándar (ej: "27 Oct 23 10:00 UTC")

		// Formatos con mes en texto (en inglés, ya que Go lo parsea así por defecto)
		"02 January 2006",
		"January 02, 2006",
		"02-Jan-2006",

		// Formato ambiguo de EE.UU. (MM-DD-YYYY), lo ponemos al final para darle menos prioridad
		"01-02-2006",
		"01/02/2006",
	}

	var fechaParseada time.Time
	var err error

	// Iteramos sobre cada layout e intentamos parsear la fecha
	for _, layout := range layouts {
		fechaParseada, err = time.Parse(layout, fechaStr)
		// Si err es nil, significa que el parseo fue exitoso y podemos detenernos
		if err == nil {
			break
		}
	}

	// Si después de probar todos los layouts, 'err' todavía no es nil,
	// significa que no encontramos un formato válido.
	if err != nil {
		return "", errors.New("formato de fecha no reconocido: '" + fechaStr + "'")
	}

	// Si llegamos aquí, 'fechaParseada' contiene la fecha correcta.
	// Ahora la formateamos a nuestro gusto.
	dia := fechaParseada.Day()
	mes := mesesEnEspanol[fechaParseada.Month()]
	ano := fechaParseada.Year()

	// Usamos Sprintf para construir el string de salida, con el día formateado a dos dígitos (ej: 07)
	resultado := fmt.Sprintf("%02d, %s, %d", dia, mes, ano)

	return resultado, nil
}

// Función auxiliar para no repetir la lógica de procesamiento de imagen
func process_and_create_image_model(c *fiber.Ctx, file_header *multipart.FileHeader, index int, user *models.PsiUserModel) (models.PostGradePic, error) {
	file_name_field := fmt.Sprintf("file_name_%d", index)
	mime_type_field := fmt.Sprintf("mime_type_%d", index)

	file_name := c.FormValue(file_name_field)
	mime_type := c.FormValue(mime_type_field)
	if file_name == "" || mime_type == "" {
		return models.PostGradePic{}, errors.New("missing file_name or mime_type")
	}

	image_data, err := psi_user_controller.ImageProcesser(file_header)
	if err != nil {
		return models.PostGradePic{}, err
	}

	image_decoded, format, err := psi_user_controller.ImageDecoder(image_data, mime_type)
	if err != nil {
		return models.PostGradePic{}, err
	}

	if len(file_name) > 50 {
		file_name = file_name[:50]
	}

	return image_controller.CreatePostGradePicModel(file_name, format, user.Username, user.ID, image_data, image_decoded)
}
