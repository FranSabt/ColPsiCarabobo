package psi_user_mapper

import (
	"fmt"
	"strconv"
	"time"

	"github.com/FranSabt/ColPsiCarabobo/src/models"
	psi_user_controller "github.com/FranSabt/ColPsiCarabobo/src/psi-user/controller"
	"github.com/FranSabt/ColPsiCarabobo/src/utils"

	"github.com/google/uuid"
)

func PsiUserCsv_To_PsiUserModel(csv psi_user_controller.PsiUserCsv, admin models.UserAdmin) models.PsiUserModel {
	// Convertir campos que requieren transformación
	fpv := 0
	if csv.FPV != "" {
		fpv = convertirAEntero(csv.FPV)
	}

	ci := 0
	if csv.CI != "" {
		ci = convertirAEntero(csv.CI)
	}

	bornDate, _ := time.Parse("2006-01-02", csv.BornDate) // Asume que la fecha está en formato YYYY-MM-DD
	hash, _ := utils.HashPassword("123456")

	// Crear y devolver el modelo
	return models.PsiUserModel{
		ID:                          uuid.New(), // Generar un nuevo UUID
		Username:                    csv.UserName,
		Email:                       csv.Email,
		Password:                    hash, //csv.Password,
		FirstName:                   csv.FirstName,
		SecondName:                  csv.SecondName,
		LastName:                    csv.LastName,
		SecondLastName:              csv.SecondLastName,
		FPV:                         fpv,
		CI:                          ci,
		Nationality:                 csv.Nationality,
		BornDate:                    bornDate,
		Genre:                       csv.Genre,
		ContactEmail:                csv.ContactEmail,
		ShowContactEmail:            convertirABool(csv.ShowContactEmail),
		PublicPhone:                 csv.PublicPhone,
		ShowPublicPhone:             convertirABool(csv.ShowPublicPhone),
		ServiceAddress:              csv.ServiceAddress,
		ShowPublicServiceAddress:    convertirABool(csv.ShowPublicServiceAddress),
		Solvent:                     convertirABool(csv.Solvent),
		ProofOfLife:                 false, // Valor por defecto
		MunicipalityCarabobo:        csv.MunicipalityCarabobo,
		PhoneCarabobo:               csv.PhoneCarabobo,
		CelPhoneCarabobo:            csv.CelPhoneCarabobo,
		StateOutside:                csv.State,
		MunicipalityOutSideCarabobo: csv.MunicipalityOutSideCarabobo,
		PhoneOutSideCarabobo:        csv.PhoneOutSideCarabobo,
		CelPhoneOutSideCarabobo:     csv.CelPhoneOutSideCarabobo,
		PsiUserColDataID:            nil, // Valor por defecto
		PrimarySpecialty:            "",
		SecondarySpecialty:          "",
		// Creation Fields
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
		CreateBy:   admin.Username,
		UpdateBy:   admin.Username,
		CreateById: &admin.ID,
		UpdateById: &admin.ID,
	}
}

////////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////

func PsiUserCsv_To_PsiUserColData(csv psi_user_controller.PsiUserCsv) models.PsiUserColData {
	// Convertir campos que requieren transformación
	graduateDate, _ := time.Parse("2006-01-02", csv.GraduateDate) // Asume que la fecha está en formato YYYY-MM-DD
	registerTitleDate, _ := time.Parse("2006-01-02", csv.RegisterTitleDate)
	dateOfLastSolvency, _ := time.Parse("2006-01-02", csv.DateOfLastSolvency)

	registerNumber := 0
	if csv.RegisterNumber != "" {
		registerNumber = convertirAEntero(csv.RegisterNumber) // Función auxiliar para convertir cadenas a enteros
	}

	// Crear y devolver el modelo
	return models.PsiUserColData{
		ID:                      uuid.New(), // Generar un nuevo UUID
		UniversityUndergraduate: csv.UniversityUndergraduate,
		GraduateDate:            graduateDate,
		MentionUndergraduate:    csv.MentionUndergraduate,
		RegisterTitleState:      csv.RegisterTitleState,
		RegisterTitleDate:       registerTitleDate,
		RegisterNumber:          registerNumber,
		RegisterFolio:           csv.RegisterFolio,
		RegisterTome:            csv.RegisterTome,
		GuildDirector:           convertirABool(csv.GuildDirector),
		SixtyFiveOrPlus:         convertirABool(csv.SixtyFiveOrPlus),
		GuildCollaborator:       convertirABool(csv.GuildCollaborator),
		PublicEmployee:          convertirABool(csv.PublicEmployee),
		UniversityProfessor:     convertirABool(csv.UniversityProfessor),
		DateOfLastSolvency:      dateOfLastSolvency,
		DoubleGuild:             convertirABool(csv.DoubleGuild),
		CPSM:                    convertirABool(csv.CPSM),
		PsiUserModelID:          uuid.Nil, // Inicialmente vacío, se debe asignar después
	}
}

////////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////

// Función auxiliar para convertir cadenas a enteros
func convertirAEntero(value string) int {
	// Implementa la lógica de conversión (puedes usar strconv.Atoi)
	// Aquí se asume que la cadena es un número válido
	num, err := strconv.Atoi(value)
	if err != nil {
		fmt.Println("Error al convertir la cadena a entero:", err)
		return 0
	}
	return num
}

// Función auxiliar para convertir cadenas a booleanos
func convertirABool(valor string) bool {
	return valor == "true" || valor == "1" || valor == "True" // Ajusta según tus necesidades
}
