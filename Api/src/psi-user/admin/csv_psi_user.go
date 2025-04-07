package psi_user_admin_presenter

import (
	"fmt"
	"net/http"
	"strconv"
	"unicode"

	psi_user_controller "github.com/FranSabt/ColPsiCarabobo/src/psi-user/controller"
	psi_user_db "github.com/FranSabt/ColPsiCarabobo/src/psi-user/db"
	psi_user_mapper "github.com/FranSabt/ColPsiCarabobo/src/psi-user/mapper"
	psi_user_request "github.com/FranSabt/ColPsiCarabobo/src/psi-user/request-structs"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
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
	result, err := psi_user_controller.ProcessCsv(src)
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

////////////////////////////////////////////
////////////////////////////////////////////
////////////////////////////////////////////

func AdminCreatePsiUser(c *fiber.Ctx, db *gorm.DB) error {
	var request psi_user_request.PsiUserCreateRequest

	// Parsear el cuerpo JSON
	if err := c.BodyParser(&request); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Cuerpo de solicitud inválido",
			"details": err.Error(),
		})
	}

	// ------- Funcion para hacer los check de los campos unicos ------- //
	can_pass, conflicts, err := psi_user_controller.CheckPsiUserUniqueFields(db, request)
	if err != nil {
		return c.Status(http.StatusCreated).JSON(fiber.Map{
			"success":   false,
			"message":   err.Error(),
			"conflicts": conflicts, // Opcional: devolver los datos creados
		})
	}
	fmt.Println("Can pass: ", can_pass)
	if !can_pass {
		return c.Status(http.StatusCreated).JSON(fiber.Map{
			"success":   false,
			"message":   "Conflict found",
			"conflicts": conflicts, // Opcional: devolver los datos creados
		})
	} // check error

	// ------- Create User in DB ------- //
	psi_user, psi_user_col_data, err := psi_user_controller.CreateNewPsiUser(db, request)
	if err != nil {
		return c.Status(http.StatusCreated).JSON(fiber.Map{
			"success": false,
			"message": "Error ehile creatin the new USer",
			"details": err.Error(),
		})
	}

	// ------- Send Email ------//
	// TODO: Send email

	// Respuesta exitosa
	return c.Status(http.StatusCreated).JSON(fiber.Map{
		"success":           true,
		"psi_user":          psi_user,
		"psi_user_col_data": psi_user_col_data, // Opcional: devolver los datos creados
	})
}

////////////////////////////////////////////
////////////////////////////////////////////
////////////////////////////////////////////

// En tu archivo de modelos o dtos
type PaginatedResponse struct {
	Data        interface{} `json:"data"`
	Total       int64       `json:"total"`
	CurrentPage int         `json:"current_page"`
	TotalPages  int         `json:"total_pages"`
	PerPage     int         `json:"per_page"`
}

// En psi_user_controller.go
func AdminGetPsiUserList(c *fiber.Ctx, db *gorm.DB) error {
	page := c.Query("page", "1")
	ci := c.Query("ci")
	fpv := c.Query("fpv")
	name := c.Query("name")

	// Validaciones
	page_valid := isOnlyPositiveNumbers(page)
	if !page_valid {
		page = "0"
	}
	ci_valid := isOnlyPositiveNumbers(ci)
	fpv_valid := isOnlyPositiveNumbers(fpv)
	name_valid := isValidName(name)

	// Construir query
	baseQuery, countQuery := psi_user_controller.CreateAdminPsiUserSearchQuery(ci, fpv, name, ci_valid, fpv_valid, name_valid, db)

	// Obtener datos paginados
	pageNum, _ := strconv.Atoi(page)
	if pageNum < 1 {
		pageNum = 1
	}
	pageSize := 10

	psiUsers, total, err := psi_user_db.SearchPsiUsersByQuery(db, baseQuery, countQuery, pageNum, pageSize)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Error searching users",
			"error":   err.Error(),
		})
	}

	// Calcular total de páginas
	totalPages := int(total) / pageSize
	if int(total)%pageSize != 0 {
		totalPages++
	}

	// Construir respuesta paginada
	response := PaginatedResponse{
		Data:        psiUsers,
		Total:       total,
		CurrentPage: pageNum,
		TotalPages:  totalPages,
		PerPage:     pageSize,
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"success": true,
		"result":  response,
	})
}

////////////////////////////////////////////
////////////////////////////////////////////
////////////////////////////////////////////

type RequestPsiUserDetails struct {
	ID string `json:"id"`
}

func GetPsiUsersByID(c *fiber.Ctx, db *gorm.DB) error {
	var request RequestPsiUserDetails

	if err := c.BodyParser(&request); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Cuerpo de solicitud inválido",
			"details": err.Error(),
		})
	}

	id, err := uuid.Parse(request.ID)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Not valid id",
			"error":   err.Error(),
		})
	}

	psi_user, psi_user_col_data, err := psi_user_db.GetPsiUserByIdDetails(db, id)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Error while trying to retrieve the PsiUser",
			"error":   err.Error(),
		})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"success":           true,
		"message":           "user found",
		"psi_user":          psi_user,
		"psi_user_col_data": psi_user_col_data,
	})

}

////////////////////////////////////////////
////////////////////////////////////////////
////////////////////////////////////////////

// 	FUNCIONES AUXILIARES //

////////////////////////////////////////////
////////////////////////////////////////////
////////////////////////////////////////////

func isOnlyPositiveNumbers(s string) bool {
	if s == "" {
		return false // String vacío no es válido
	}

	for _, char := range s {
		if char < '0' || char > '9' {
			return false // Si encuentra cualquier caracter que no sea dígito
		}
	}
	return true
}

////////////////////////////////////////////
////////////////////////////////////////////
////////////////////////////////////////////

func isValidName(s string) bool {
	for _, r := range s {
		if !unicode.IsLetter(r) {
			return false
		}
	}

	return true
}
