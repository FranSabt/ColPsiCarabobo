package psi_user_admin_presenter

import (
	"fmt"
	"net/http"
	"reflect"
	"regexp"
	"strconv"
	"time"
	"unicode"

	admin_controller "github.com/FranSabt/ColPsiCarabobo/src/admin/controller"
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

	// ---- EXISTE EL ADMIN? ---- //
	admin_id := c.FormValue("admin_id")
	if admin_id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Falta el ID del administrador",
			"details": "El campo admin_id es requerido",
		})
	}

	admin_uuid, err := uuid.Parse(admin_id)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "invalid id format",
			"details": err.Error(),
		})
	}

	admin, err := admin_controller.GetAdminById(admin_uuid, db)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "invalid id format",
			"details": err.Error(),
		})
	}

	if !admin.CanCreatePsi {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "user cant create PsiUser",
			"details": fmt.Sprintf("user %v cant create a PsiUser", admin.Username),
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
		err := db.Transaction(func(tx *gorm.DB) error {
			// Mapear a los modelos
			psi_model_mapped := psi_user_mapper.PsiUserCsv_To_PsiUserModel(r, admin)

			// Crear PsiUserModel
			if err := psi_user_db.CreatePsiUseDb(tx, psi_model_mapped); err != nil {
				return fmt.Errorf("error al crear PsiUserModel: %w", err)
			}

			// Mapear y crear PsiUserColData
			psi_col_data_mapped := psi_user_mapper.PsiUserCsv_To_PsiUserColData(r)
			psi_col_data_mapped.PsiUserModelID = psi_model_mapped.ID

			if err := psi_user_db.CreatePsiColDataDb(tx, psi_col_data_mapped); err != nil {
				return fmt.Errorf("error al crear PsiUserColData: %w", err)
			}

			// Actualizar PsiUserModel con el ID de ColData
			err := tx.Model(&psi_model_mapped).
				Update("psi_user_col_data_id", psi_col_data_mapped.ID).Error
			if err != nil {
				return fmt.Errorf("error al actualizar PsiUserModel con PsiUserColDataID: %w", err)
			}

			// Todo salió bien
			return nil
		})

		// Si algo falló, agregamos a los fallidos
		if err != nil {
			fail = append(fail, struct {
				Record interface{}
				Error  string
			}{
				Record: r,
				Error:  err.Error(),
			})
			continue
		}

		count++
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

	// ---- EXISTE EL ADMIN? ---- //
	admin_id := request.AdminId

	admin_uuid, err := uuid.Parse(admin_id)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "invalid id format",
			"details": err.Error(),
		})
	}

	admin, err := admin_controller.GetAdminById(admin_uuid, db)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "invalid id format",
			"details": err.Error(),
		})
	}

	if !admin.CanCreatePsi {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "user cant create PsiUser",
			"details": fmt.Sprintf("user %v cant create a PsiUser", admin.Username),
		})
	}

	// ------- CHECK DE CAMPOS DEL REQUEST ------- //
	errors_list := isOkPsiUserFields(request)
	if len(errors_list) > 0 {
		return c.Status(http.StatusCreated).JSON(fiber.Map{
			"success":   false,
			"message":   "errors in request fields",
			"conflicts": errors_list, // Opcional: devolver los datos creados
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
	psi_user, psi_user_col_data, err := psi_user_controller.CreateNewPsiUser(db, request, admin)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Error while creating the new User",
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

func PatchPsiUserByID(c *fiber.Ctx, db *gorm.DB) error {
	var request psi_user_request.PsiUserUpdateRequest

	if err := c.BodyParser(&request); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Cuerpo de solicitud inválido",
			"details": err.Error(),
		})
	}

	if !hasFieldsToUpdate(request) {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Se requiere al menos un campo adicional al ID para actualizar",
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

	err = psi_user_controller.UpdatePsiUserModelFields(psi_user, &request)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Error while trying to update the PsiUser",
			"error":   err.Error(),
		})
	}

	err = psi_user_controller.UpdatePsiUserColDataFields(psi_user_col_data, &request)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Error while trying to update the PsiUser",
			"error":   err.Error(),
		})
	}

	err = psi_user_db.SaveUpdatedPsiUser(db, psi_user, psi_user_col_data)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Error while trying to update the PsiUse in DBr",
			"error":   err.Error(),
		})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"success":           true,
		"message":           "user updated",
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

////////////////////////////////////////////
////////////////////////////////////////////
////////////////////////////////////////////

func hasFieldsToUpdate(request psi_user_request.PsiUserUpdateRequest) bool {
	// Usamos reflect para verificar dinámicamente todos los campos
	v := reflect.ValueOf(request)
	typeOfS := v.Type()

	for i := 0; i < v.NumField(); i++ {
		fieldName := typeOfS.Field(i).Name
		fieldValue := v.Field(i).Interface()

		// Saltamos el campo ID
		if fieldName == "ID" {
			continue
		}

		// Verificamos si el campo es un puntero y no es nil
		if reflect.ValueOf(fieldValue).Kind() == reflect.Ptr && !reflect.ValueOf(fieldValue).IsNil() {
			return true
		}
	}

	return false
}

////////////////////////////////////////////
////////////////////////////////////////////
////////////////////////////////////////////

func isOkPsiUserFields(request psi_user_request.PsiUserCreateRequest) []string {
	var errors []string

	if request.Nationality != "v" && request.Nationality != "e" {
		errors = append(errors, "invalid nationality")
	}

	if request.Genre != "m" && request.Genre != "f" {
		errors = append(errors, "invalid genre")
	}

	if request.PhoneCarabobo != "" && !isValidVenezuelanPhoneNumber(request.PhoneCarabobo) {
		errors = append(errors, "invalid phone Carabobo")

	}

	if request.CelPhoneCarabobo != "" && !isValidVenezuelanPhoneNumber(request.CelPhoneCarabobo) {
		errors = append(errors, "invalid cell phone Carabobo")

	}

	if request.PublicPhone != "" && !isValidVenezuelanPhoneNumber(request.PublicPhone) {
		errors = append(errors, "invalid public phone Carabobo")

	}

	if request.CelPhoneOutSideCarabobo != "" && !isValidVenezuelanPhoneNumber(request.CelPhoneOutSideCarabobo) {
		errors = append(errors, "invalid cel phone outside Carabobo")

	}

	if request.PhoneOutSideCarabobo != "" && !isValidVenezuelanPhoneNumber(request.PhoneOutSideCarabobo) {
		errors = append(errors, "invalid phone outside Carabobo")

	}

	if !isValidDateYYYYMMDD(request.BornDate) {
		errors = append(errors, "invalid born date")
	}

	if !isValidDateYYYYMMDD(request.GraduateDate) {
		errors = append(errors, "invalid graduate date")
	}

	if !isValidDateYYYYMMDD(request.DateOfLastSolvency) {
		errors = append(errors, "invalid last solvency date")
	}

	return errors
}

////////////////////////////////////////////
////////////////////////////////////////////
////////////////////////////////////////////

var nonDigitsRegex = regexp.MustCompile(`\D`)
var venezuelanPhoneRegex = regexp.MustCompile(`^(0|58)(412|414|424|416|426|2\d{2})\d{7}$`)

// IsValidVenezuelanPhoneNumber recibe un número de teléfono en cualquier formato y retorna true si es válido.
func isValidVenezuelanPhoneNumber(phone string) bool {

	cleanedPhone := nonDigitsRegex.ReplaceAllString(phone, "")

	return venezuelanPhoneRegex.MatchString(cleanedPhone)
}

////////////////////////////////////////////
////////////////////////////////////////////
////////////////////////////////////////////

// IsValidDateYYYYMMDD verifica si un string es una fecha válida en formato "YYYY-MM-DD".
// Retorna 'true' si el formato y la fecha son válidos, de lo contrario 'false'.
func isValidDateYYYYMMDD(dateString string) bool {

	const layout = "2006-01-02"

	_, err := time.Parse(layout, dateString)

	// Si no hay error (err == nil), la fecha es válida.
	return err == nil
}
