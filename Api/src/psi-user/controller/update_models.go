package psi_user_controller

import (
	"github.com/FranSabt/ColPsiCarabobo/src/models"
	psi_user_request "github.com/FranSabt/ColPsiCarabobo/src/psi-user/request-structs"
	"github.com/FranSabt/ColPsiCarabobo/src/utils"
)

// updatePsiUserModelFields actualiza los campos de PsiUserModel con los valores no nulos del request
func UpdatePsiUserModelFields(psiUser *models.PsiUserModel, request *psi_user_request.PsiUserUpdateRequest) error {
	if request.Username != nil {
		psiUser.Username = *request.Username
	}
	if request.FirstName != nil {
		psiUser.FirstName = *request.FirstName
	}
	if request.SecondName != nil {
		psiUser.SecondName = *request.SecondName
	}
	if request.LastName != nil {
		psiUser.LastName = *request.LastName
	}
	if request.SecondLastName != nil {
		psiUser.SecondLastName = *request.SecondLastName
	}
	if request.Email != nil {
		psiUser.Email = *request.Email
	}
	if request.FPV != nil {
		psiUser.FPV = *request.FPV
	}
	if request.CI != nil {
		psiUser.CI = *request.CI
	}
	if request.Nationality != nil {
		psiUser.Nationality = *request.Nationality
	}
	if request.BornDate != nil {
		date, err := utils.ParseDateString(*request.BornDate)
		if err != nil {
			return err
		}
		psiUser.BornDate = date
	}
	if request.Genre != nil {
		psiUser.Genre = *request.Genre
	}
	if request.ContactEmail != nil {
		psiUser.ContactEmail = *request.ContactEmail
	}
	if request.PublicPhone != nil {
		psiUser.PublicPhone = *request.PublicPhone
	}
	if request.ServiceAddress != nil {
		psiUser.ServiceAddress = *request.ServiceAddress
	}
	if request.MunicipalityCarabobo != nil {
		psiUser.MunicipalityCarabobo = *request.MunicipalityCarabobo
	}
	if request.PhoneCarabobo != nil {
		psiUser.PhoneCarabobo = *request.PhoneCarabobo
	}
	if request.CelPhoneCarabobo != nil {
		psiUser.CelPhoneCarabobo = *request.CelPhoneCarabobo
	}
	if request.StateOutside != nil {
		psiUser.StateOutside = *request.StateOutside
	}
	if request.MunicipalityOutSideCarabobo != nil {
		psiUser.MunicipalityOutSideCarabobo = *request.MunicipalityOutSideCarabobo
	}
	if request.PhoneOutSideCarabobo != nil {
		psiUser.PhoneOutSideCarabobo = *request.PhoneOutSideCarabobo
	}
	if request.CelPhoneOutSideCarabobo != nil {
		psiUser.CelPhoneOutSideCarabobo = *request.CelPhoneOutSideCarabobo
	}

	return nil
}

// updatePsiUserColDataFields actualiza los campos de PsiUserColData con los valores no nulos del request
func UpdatePsiUserColDataFields(colData *models.PsiUserColData, request *psi_user_request.PsiUserUpdateRequest) error {
	if request.UniversityUndergraduate != nil {
		colData.UniversityUndergraduate = *request.UniversityUndergraduate
	}
	if request.GraduateDate != nil {
		date, err := utils.ParseDateString(*request.GraduateDate)
		if err != nil {
			return err
		}
		colData.GraduateDate = date
	}
	if request.MentionUndergraduate != nil {
		colData.MentionUndergraduate = *request.MentionUndergraduate
	}
	if request.RegisterTitleState != nil {
		colData.RegisterTitleState = *request.RegisterTitleState
	}
	if request.RegisterTitleDate != nil {
		date, err := utils.ParseDateString(*request.RegisterTitleDate)
		if err != nil {
			return err
		}
		colData.RegisterTitleDate = date
	}
	if request.RegisterNumber != nil {
		colData.RegisterNumber = *request.RegisterNumber
	}
	if request.RegisterFolio != nil {
		colData.RegisterFolio = *request.RegisterFolio
	}
	if request.RegisterTome != nil {
		colData.RegisterTome = *request.RegisterTome
	}
	if request.GuildDirector != nil {
		colData.GuildDirector = *request.GuildDirector
	}
	if request.SixtyFiveOrPlus != nil {
		colData.SixtyFiveOrPlus = *request.SixtyFiveOrPlus
	}
	if request.GuildCollaborator != nil {
		colData.GuildCollaborator = *request.GuildCollaborator
	}
	if request.PublicEmployee != nil {
		colData.PublicEmployee = *request.PublicEmployee
	}
	if request.UniversityProfessor != nil {
		colData.UniversityProfessor = *request.UniversityProfessor
	}
	if request.DateOfLastSolvency != nil {
		date, err := utils.ParseDateString(*request.RegisterTitleDate)
		if err != nil {
			return err
		}
		colData.DateOfLastSolvency = date
	}
	if request.DoubleGuild != nil {
		colData.DoubleGuild = *request.DoubleGuild
	}
	if request.CPSM != nil {
		colData.CPSM = *request.CPSM
	}

	return nil
}
