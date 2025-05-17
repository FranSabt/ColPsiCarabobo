package psiuser_presenter

import (
	"errors"
	"log"
	"regexp"

	"github.com/FranSabt/ColPsiCarabobo/src/models"
	psi_user_db "github.com/FranSabt/ColPsiCarabobo/src/psi-user/db"
	psi_user_request "github.com/FranSabt/ColPsiCarabobo/src/psi-user/request-structs"
	"github.com/FranSabt/ColPsiCarabobo/src/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func GetPsiUserSelfInfo(c *fiber.Ctx, db *gorm.DB) error {
	id := c.Query("id")

	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "No Id",
		})
	}

	uuid_pased, err := uuid.Parse(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "not valid id",
		})
	}

	psiuser_model, psiuser_data, err := psi_user_db.GetPsiUserByIdDetails(db, uuid_pased)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "not valid id",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success":      true,
		"error":        "",
		"psiuser":      psiuser_model,
		"psiuser_data": psiuser_data,
	})
}

//////////////////////////////////////////////////////
//////////////////////////////////////////////////////
//////////////////////////////////////////////////////

func UpdatePsiUserSelfInfo(c *fiber.Ctx, db *gorm.DB) error {
	var request psi_user_request.PsiUserUpdateRequestSelf

	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Invalid request format",
		})
	}

	// ------- Check de los campos ------- //
	can_pass, err := checkUpdateselfField(request)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success":       false,
			"error":         "No field to update",
			"error_message": err.Error(),
		})
	}
	if !can_pass {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "No field to update",
		})
	}

	// ------- Verificar el usuario ------ //
	uuid_pased, err := uuid.Parse(request.ID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "not valid id",
		})
	}

	psiuser_model, err := psi_user_db.GetPsiUserById(db, uuid_pased)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}
	if psiuser_model == nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "user or do not exist",
		})
	}

	// ------ Check Pass ------ //
	if !utils.CheckPasswordHash(request.Password, psiuser_model.Password) {
		// Log de seguridad
		log.Printf("Failed login attempt for user: %s", psiuser_model.Username)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"error":   "Invalid credentials2",
		})
	}
	// ------ Ajustar modelo ------ //
	model_updated := modifiPsiUSerModel(psiuser_model, request)

	// ------ Update ----- //
	err = psi_user_db.SaveUpdatedPsiUserOnly(db, model_updated)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})

	}

	return c.Status(fiber.StatusAccepted).JSON(fiber.Map{
		"success": true,
		"error":   "",
	})

}

///////////////////////////////////////////////////
///////////////////////////////////////////////////
///////////////////////////////////////////////////

func checkUpdateselfField(psi_user psi_user_request.PsiUserUpdateRequestSelf) (bool, error) {
	can_pass := false

	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	venezuelaPhoneRegex := regexp.MustCompile(`^0(412|414|424|416|426|212)[0-9]{7}$`)

	if psi_user.Email != nil {
		if !emailRegex.MatchString(*psi_user.Email) {
			return false, errors.New("email no válido")
		}
		can_pass = true
	}

	if psi_user.ContactEmail != nil {
		if !emailRegex.MatchString(*psi_user.ContactEmail) {
			return false, errors.New("email de contacto no válido")
		}
		can_pass = true
	}

	if psi_user.CelPhoneCarabobo != nil {
		if !venezuelaPhoneRegex.MatchString(*psi_user.CelPhoneCarabobo) {
			return false, errors.New("número de celular de Carabobo no válido")
		}
		can_pass = true
	}

	if psi_user.CelPhoneOutSideCarabobo != nil {
		if !venezuelaPhoneRegex.MatchString(*psi_user.CelPhoneOutSideCarabobo) {
			return false, errors.New("número de celular fuera de Carabobo no válido")
		}
		can_pass = true
	}

	if psi_user.PhoneCarabobo != nil {
		if !venezuelaPhoneRegex.MatchString(*psi_user.PhoneCarabobo) {
			return false, errors.New("número de teléfono de Carabobo no válido")
		}
		can_pass = true
	}

	if psi_user.NewPassword1 != nil {
		if psi_user.NewPassword2 == nil || *psi_user.NewPassword1 != *psi_user.NewPassword2 {
			return false, errors.New("las contraseñas nuevas no coinciden")
		}
		can_pass = true
	}

	if psi_user.ServiceAddress != nil {
		can_pass = true
	}

	// Verificación de booleanos (al menos uno presente)
	if psi_user.ShowContacEmail != nil || psi_user.ShowPublicPhone != nil || psi_user.ShowServiceAddress != nil {
		can_pass = true
	}

	if !can_pass {
		return false, errors.New("no se proporcionó ningún campo válido para actualizar")
	}

	return true, nil
}

///////////////////////////////////////////////////
///////////////////////////////////////////////////
///////////////////////////////////////////////////

func modifiPsiUSerModel(
	psiuser_model *models.PsiUserModel,
	psi_user_request psi_user_request.PsiUserUpdateRequestSelf,
) *models.PsiUserModel {

	// ---- Campos de PsiUser ---- //
	//
	if psi_user_request.Username != nil {
		psiuser_model.Username = *psi_user_request.Username
	}

	if psi_user_request.Email != nil {
		psiuser_model.Email = *psi_user_request.Email
	}

	if psi_user_request.NewPassword1 != nil {

		if psi_user_request.NewPassword1 == psi_user_request.NewPassword2 {
			hash, err := utils.HashPassword(*psi_user_request.NewPassword1)

			if err == nil {
				psiuser_model.Password = hash
			}

		}
	}

	// Contact Information
	if psi_user_request.ContactEmail != nil {
		psiuser_model.ContactEmail = *psi_user_request.ContactEmail
	}

	if psi_user_request.ShowContacEmail != nil {
		psiuser_model.ShowContactEmail = *psi_user_request.ShowContacEmail
	}

	if psi_user_request.PublicPhone != nil {
		psiuser_model.PublicPhone = *psi_user_request.PublicPhone
	}

	if psi_user_request.ShowPublicPhone != nil {
		psiuser_model.ShowPublicPhone = *psi_user_request.ShowPublicPhone
	}

	if psi_user_request.ServiceAddress != nil {
		psiuser_model.ServiceAddress = *psi_user_request.ServiceAddress
	}

	if psi_user_request.ShowServiceAddress != nil {
		psiuser_model.ShowPublicServiceAddress = *psi_user_request.ShowServiceAddress
	}

	// Address Information
	if psi_user_request.MunicipalityCarabobo != nil {
		psiuser_model.MunicipalityCarabobo = *psi_user_request.MunicipalityCarabobo
	}

	if psi_user_request.PhoneCarabobo != nil {
		psiuser_model.PhoneCarabobo = *psi_user_request.PhoneCarabobo
	}

	if psi_user_request.CelPhoneCarabobo != nil {
		psiuser_model.CelPhoneCarabobo = *psi_user_request.CelPhoneCarabobo
	}

	// Outside Carabobo Address
	if psi_user_request.StateOutside != nil {
		psiuser_model.StateOutside = *psi_user_request.StateOutside
	}

	if psi_user_request.MunicipalityOutSideCarabobo != nil {
		psiuser_model.MunicipalityOutSideCarabobo = *psi_user_request.MunicipalityOutSideCarabobo
	}

	if psi_user_request.PhoneOutSideCarabobo != nil {
		psiuser_model.PhoneOutSideCarabobo = *psi_user_request.PhoneOutSideCarabobo
	}

	if psi_user_request.CelPhoneOutSideCarabobo != nil {
		psiuser_model.CelPhoneOutSideCarabobo = *psi_user_request.CelPhoneOutSideCarabobo
	}

	return psiuser_model
}
