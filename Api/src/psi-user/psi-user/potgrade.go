package psiuser_presenter

import (
	"errors"
	"fmt"
	"image"
	"net/http"
	"strings"
	"time"

	psi_user_controller "github.com/FranSabt/ColPsiCarabobo/src/psi-user/controller"
	psi_user_db "github.com/FranSabt/ColPsiCarabobo/src/psi-user/db"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func CreatePsiUserPotgrade(c *fiber.Ctx, db *gorm.DB) error {
	// --- Recopilación de datos del formulario ---
	psi_user_id := c.FormValue("psi_user_id")
	post_grade_title := c.FormValue("post_grade_title")
	post_grade_university := c.FormValue("post_grade_university")
	post_grade_graduation_year := c.FormValue("post_grade_graduation_year")
	post_grade_description := c.FormValue("post_grade_description")

	if psi_user_id == "" || post_grade_title == "" || post_grade_university == "" || post_grade_graduation_year == "" || post_grade_description == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "empty field on the request",
		})
	}

	// ---- Validar fecha del titulo
	post_grade_graduation_year_clean, err := formatearFecha(post_grade_graduation_year)

	// --- Validación del Usuario ---
	uuid, err := uuid.Parse(psi_user_id)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "invalid id format",
			"details": err.Error(),
		})
	}

	user, err := psi_user_db.GetPsiUserById(db, uuid)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Can't connect to DB",
			"success": false,
		})
	}
	if user == nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "User doesn't exist",
			"success": false,
		})
	}
	if !user.Solvent {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "User is inactive",
			"success": false,
		})
	}

	// --- Procesamiento de Imágenes ---
	var image_data_array [3][]byte
	var err_1, err_2, err_3 error // Declarar errores para reutilizar

	file_header_1, errFormFile1 := c.FormFile("image_1")
	if errFormFile1 == nil && file_header_1 != nil {
		image_data_array[0], err_1 = psi_user_controller.ImageProcesser(file_header_1)
	}

	file_header_2, errFormFile2 := c.FormFile("image_2")
	if errFormFile2 == nil && file_header_2 != nil {
		image_data_array[1], err_2 = psi_user_controller.ImageProcesser(file_header_2)
	}

	file_header_3, errFormFile3 := c.FormFile("image_3")
	if errFormFile3 == nil && file_header_3 != nil {
		image_data_array[2], err_3 = psi_user_controller.ImageProcesser(file_header_3)
	}

	// --- Decodificación de Imágenes ---

	// Declarar variables fuera de los 'if' para que su scope sea toda la función
	var image_1, image_2, image_3 image.Image
	var format_1, format_2, format_3 string
	var filename_1, filename_2, filename_3 string

	if image_data_array[0] != nil && err_1 == nil {
		filename_1 = c.FormValue("file_name_1")
		mime_type_1 := c.FormValue("mime_type_1")
		if filename_1 == "" || mime_type_1 == "" {
			err_1 = errors.New("invalid mime_type or filename for image 1")
		} else {
			// CORRECTO: Usar '=' para asignar a las variables ya declaradas
			image_1, format_1, err_1 = psi_user_controller.ImageDecoder(image_data_array[0], mime_type_1)
			if len(filename_1) > 50 {
				filename_1 = filename_1[:50]
			}
		}
	}

	if image_data_array[1] != nil && err_2 == nil {
		filename_2 = c.FormValue("file_name_2")
		// CORRECCIÓN: Usar mime_type_2 en lugar de mime_type_1
		mime_type_2 := c.FormValue("mime_type_2")
		if filename_2 == "" || mime_type_2 == "" {
			err_2 = errors.New("invalid mime_type or filename for image 2")
		} else {
			// CORRECTO: Usar '=' para asignar
			image_2, format_2, err_2 = psi_user_controller.ImageDecoder(image_data_array[1], mime_type_2)
			if len(filename_2) > 50 {
				filename_2 = filename_2[:50]
			}
		}
	}

	if image_data_array[2] != nil && err_3 == nil {
		filename_3 = c.FormValue("file_name_3")
		// CORRECCIÓN: Usar mime_type_3 en lugar de mime_type_1
		mime_type_3 := c.FormValue("mime_type_3")
		if filename_3 == "" || mime_type_3 == "" {
			err_3 = errors.New("invalid mime_type or filename for image 3")
		} else {
			// CORRECTO: Usar '=' para asignar
			image_3, format_3, err_3 = psi_user_controller.ImageDecoder(image_data_array[2], mime_type_3)
			if len(filename_3) > 50 {
				filename_3 = filename_3[:50]
			}
		}
	}

	// --- Verificación final de errores ---
	var error_messages []string
	if err_1 != nil {
		error_messages = append(error_messages, "Image 1: "+err_1.Error())
	}
	if err_2 != nil {
		error_messages = append(error_messages, "Image 2: "+err_2.Error())
	}
	if err_3 != nil {
		error_messages = append(error_messages, "Image 3: "+err_3.Error())
	}

	if len(error_messages) > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "error processing one or more images",
			// CORRECCIÓN: La sintaxis para un array de strings en JSON es esta
			"details": error_messages,
		})
	}

	// Crear el modelo
	postgraduate_model := psi_user_controller.CreatePsiUserPostGradeModel(
		uuid,
		post_grade_title,
		post_grade_university,
		post_grade_graduation_year_clean,
		post_grade_description,
	)

	if image_1 != nil && err_1 == nil {
		// crear el PostPicModel
		// guardar el PostPicModel
		// llevar el id a postgraduate_model
	}
	// Las variables image_1, format_1, filename_1 (y las demás) ya están
	// disponibles para ser usadas.

	return nil // Placeholder
}

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
