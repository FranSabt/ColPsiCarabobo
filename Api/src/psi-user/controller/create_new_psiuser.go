package psi_user_controller

import (
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/FranSabt/ColPsiCarabobo/src/models"
	psi_user_db "github.com/FranSabt/ColPsiCarabobo/src/psi-user/db"
	psi_user_request "github.com/FranSabt/ColPsiCarabobo/src/psi-user/request-structs"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func CreateNewPsiUser(db *gorm.DB, request psi_user_request.PsiUserCreateRequest) (*models.PsiUserModel, *models.PsiUserColData, error) {

	// create the psi user
	psi_user, err := createPsiuserModel(request)
	if err != nil {
		return nil, nil, err
	}

	fmt.Printf("%v", psi_user)
	// if ok, crete psi-user coldata
	psi_user_col_data, err := createPsiUserColDataModel(request, psi_user.ID)
	if err != nil {
		fmt.Printf("%s", err.Error())
		return nil, nil, err
	}
	// giver to psi_user the col_data ID
	psi_user.PsiUserColDataID = &psi_user_col_data.ID

	// save the psi user
	err = psi_user_db.CreatePsiUseDb2(db, psi_user)
	if err != nil {
		fmt.Printf("%s", err.Error())
		return nil, nil, err
	}
	err = psi_user_db.CreatePsiColDataDb2(db, psi_user_col_data)
	if err != nil {
		fmt.Printf("%s", err.Error())
		return psi_user, nil, err
	}

	// return both elements
	return psi_user, psi_user_col_data, nil
}

///////////////////////////////////////////////
///////////////////////////////////////////////
///////////////////////////////////////////////

// Auxiliar Functions //

///////////////////////////////////////////////
///////////////////////////////////////////////
///////////////////////////////////////////////

func createPsiuserModel(request psi_user_request.PsiUserCreateRequest) (*models.PsiUserModel, error) {
	// Parsear fecha de nacimiento
	bornDate, err := ParseDateString(request.BornDate)
	if err != nil {
		return nil, fmt.Errorf("invalid birth date: %v", err)
	}

	// Generar contraseña aleatoria segura y su hash
	password := generateRandomPassword(12) // Longitud de 12 caracteres
	hashedPassword, err := hashPassword(password)
	if err != nil {
		return nil, fmt.Errorf("could not hash password: %v", err)
	}

	// Generar nuevo UUID para el usuario
	user_id := uuid.New()

	// Crear modelo de usuario
	psiUser := models.PsiUserModel{
		ID:       user_id,
		Username: request.Username,
		Email:    request.Email,
		Password: hashedPassword,
		// IDENTITY
		FirstName:      request.FirstName,
		SecondName:     request.SecondName,
		LastName:       request.LastName,
		SecondLastName: request.SecondLastName,
		FPV:            request.FPV,
		CI:             request.CI,
		Nationality:    request.Nationality,
		BornDate:       bornDate,
		Genre:          request.Genre,
		// Contact
		ContactEmail:             request.ContactEmail,
		ShowContactEmail:         false,
		PublicPhone:              request.PublicPhone,
		ShowPublicPhone:          false,
		ServiceAddress:           request.ServiceAddress,
		ShowPublicServiceAddress: false,
		// PsiCol
		Solvent:     true,
		ProofOfLife: true,
		// Carabobo Direction
		MunicipalityCarabobo: request.MunicipalityCarabobo,
		PhoneCarabobo:        request.PhoneCarabobo,
		CelPhoneCarabobo:     request.CelPhoneCarabobo,
		// Outside Carabobo Direction
		StateOutside:                request.StateOutside,
		MunicipalityOutSideCarabobo: request.MunicipalityOutSideCarabobo,
		PhoneOutSideCarabobo:        request.PhoneOutSideCarabobo,
		CelPhoneOutSideCarabobo:     request.CelPhoneOutSideCarabobo,
		// Relación con PsiUserColData
		PsiUserColDataID: &uuid.Nil,
	}

	// Validación básica del modelo creado
	if psiUser.Username == "" || psiUser.Email == "" || psiUser.CI <= 0 || psiUser.FPV <= 0 {
		return nil, errors.New("required fields are missing")
	}

	return &psiUser, nil
}

// Función para generar contraseña aleatoria
func generateRandomPassword(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*()_+"
	password := make([]byte, length)
	for i := range password {
		password[i] = charset[rand.Intn(len(charset))]
	}
	return string(password)
}

// Función para hashear la contraseña
func hashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

///////////////////////////////////////////////
///////////////////////////////////////////////
///////////////////////////////////////////////

func createPsiUserColDataModel(request psi_user_request.PsiUserCreateRequest,
	psi_user_id uuid.UUID) (*models.PsiUserColData, error) {
	graduate_date, err := ParseDateString(request.BornDate)
	if err != nil {
		return nil, fmt.Errorf("invalid graduate date: %v", err)
	}

	register_title_date, err := ParseDateString(request.RegisterTitleDate)
	if err != nil {
		return nil, fmt.Errorf("invalid graduate date: %v", err)
	}

	date_last_solvency, err := ParseDateString(request.DateOfLastSolvency)
	if err != nil {
		return nil, fmt.Errorf("inavel date of last solvency: %v", err)
	}

	user_id := uuid.New()

	psi_user_col_data := models.PsiUserColData{
		ID: user_id,

		// Undergraduate Data
		UniversityUndergraduate: request.UniversityUndergraduate,
		GraduateDate:            graduate_date,
		MentionUndergraduate:    request.MentionUndergraduate,

		// Undergraduate Data Title Register
		RegisterTitleState: request.RegisterTitleState, // Puedes usar un enum si lo defines
		RegisterTitleDate:  register_title_date,
		RegisterNumber:     request.RegisterNumber,
		RegisterFolio:      request.RegisterFolio, // Puedes usar un enum si lo defines
		RegisterTome:       request.RegisterTome,

		// Professional Data
		GuildDirector:       request.GuildDirector,
		SixtyFiveOrPlus:     request.SixtyFiveOrPlus,
		GuildCollaborator:   request.GuildCollaborator,
		PublicEmployee:      request.PublicEmployee,
		UniversityProfessor: request.UniversityProfessor,

		// Otros campos
		DateOfLastSolvency: date_last_solvency,
		DoubleGuild:        request.DoubleGuild,
		CPSM:               request.CPSM,

		// Relación con PsiUserModel
		PsiUserModelID: psi_user_id, // Clave foránea
	}

	if psi_user_col_data.UniversityUndergraduate == "" || psi_user_col_data.RegisterTitleState == "" {
		return nil, errors.New("required fields are missing")
	}

	return &psi_user_col_data, nil
}

///////////////////////////////////////////////
///////////////////////////////////////////////
///////////////////////////////////////////////

func ParseDateString(dateString string) (time.Time, error) {
	// Lista de formatos de fecha comunes que intentaremos analizar
	formats := []string{
		"2006-01-02",   // Formato ISO (YYYY-MM-DD)
		"02/01/2006",   // Formato DD/MM/YYYY
		"01/02/2006",   // Formato MM/DD/YYYY
		"Jan 02, 2006", // Ej: "Dec 25, 2023"
		"02-Jan-2006",  // Ej: "25-Dec-2023"
		time.RFC3339,   // Formato ISO con zona horaria
	}

	var parsedTime time.Time
	var err error

	// Intentar parsear con cada formato hasta que uno funcione
	for _, format := range formats {
		parsedTime, err = time.Parse(format, dateString)
		if err == nil {
			return parsedTime, nil
		}
	}

	return time.Time{}, fmt.Errorf("no se pudo parsear la fecha: %v, formatos intentados: %v", dateString, formats)
}
