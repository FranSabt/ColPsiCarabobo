package psi_user_mapper

import (
	"time"

	"github.com/FranSabt/ColPsiCarabobo/src/models"
	"github.com/google/uuid"
)

type PsiUserPublicData struct {
	ID uuid.UUID `json:"id"`

	// Identity
	FirstName      string `json:"first_name"`
	SecondName     string `json:"second_name,omitempty"`
	LastName       string `json:"last_name"`
	SecondLastName string `json:"second_last_name,omitempty"`
	FPV            int    `json:"fpv"`
	CI             int    `json:"ci"`
	Nationality    string `json:"nationality"`
	BornDate       string `json:"born_date"`
	Genre          string `json:"genre"`

	// Contact
	ContactEmail   string `json:"contact_email"`
	PublicPhone    string `json:"public_phone,omitempty"`
	ServiceAddress string `json:"service_address,omitempty"`

	// Col Data
	UniversityUndergraduate string `json:"university_undergraduate"`
	GraduateDate            string `json:"graduate_year"` // Usar time.Time para fechas
	MentionUndergraduate    string `json:"mention_undergraduate"`
	PrimarySpecialty        string `json:"primary_specialty"`
	SecondarySpecialty      string `json:"secondary_specialty"`
}

func PsiUserDataToPublic(psi_user *models.PsiUserModel, col_data *models.PsiUserColData) *PsiUserPublicData {
	is_solvent := psi_user.Solvent

	psi_user_public := PsiUserPublicData{
		ID:          psi_user.ID,
		FirstName:   psi_user.FirstName,
		SecondName:  psi_user.SecondName,
		LastName:    psi_user.SecondLastName,
		FPV:         psi_user.FPV,
		CI:          psi_user.CI,
		Nationality: psi_user.Nationality,
		BornDate:    psi_user.BornDate.String(),
		Genre:       psi_user.Genre,
		/* --- No puede publicitarse si no esta solvente --- */
		ContactEmail:            isContactEmail(is_solvent, psi_user.ShowContactEmail, psi_user.Email),
		PublicPhone:             isPublicPhone(is_solvent, psi_user.ShowPublicPhone, psi_user.PublicPhone),
		ServiceAddress:          isServiceAddress(is_solvent, psi_user.ShowPublicServiceAddress, psi_user.ServiceAddress),
		UniversityUndergraduate: isUniversityUndergraduate(is_solvent, col_data.ShowUniversityUndergraduate, col_data.UniversityUndergraduate),
		GraduateDate:            isGraduateDate(is_solvent, col_data.ShowGraduateDate, col_data.GraduateDate),
		MentionUndergraduate:    isMentionUndergraduate(is_solvent, col_data.ShowMentionUndergraduate, col_data.MentionUndergraduate),
		PrimarySpecialty:        psi_user.PrimarySpecialty,
		SecondarySpecialty:      psi_user.SecondarySpecialty,
	}

	return &psi_user_public
}

func isContactEmail(solvent, show_contact_email bool, email string) string {
	if solvent && show_contact_email {
		return email
	}

	return ""
}

func isPublicPhone(solvent, show_public_phone bool, public_phone string) string {
	if solvent && show_public_phone {
		return public_phone
	}

	return ""
}

func isServiceAddress(solvent, show_service_address bool, service_address string) string {
	if solvent && show_service_address {
		return service_address
	}

	return ""
}

func isUniversityUndergraduate(solvent, show_university_undergraduate bool, university_undergraduate string) string {
	if solvent && show_university_undergraduate {
		return university_undergraduate
	}

	return ""
}

func isGraduateDate(solvent, show_graduate_date bool, graduate_date time.Time) string {
	if solvent && show_graduate_date {
		return graduate_date.String()
	}

	return ""
}

func isMentionUndergraduate(solvent, show_mention_undergraduate bool, mention_undergraduate string) string {
	if solvent && show_mention_undergraduate {
		return mention_undergraduate
	}

	return ""
}
