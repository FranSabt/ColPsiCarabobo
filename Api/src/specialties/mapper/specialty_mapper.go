package specialties_mapper

import (
	"time"

	"github.com/FranSabt/ColPsiCarabobo/src/models"
	specialties_structs "github.com/FranSabt/ColPsiCarabobo/src/specialties/request-structs"
)

func SpecialtyRequestToSpecialtyModel(specialtyRequest specialties_structs.SpecialtiesRequest) models.PsiSpecialty {
	specialtyModel := models.PsiSpecialty{
		Name:        specialtyRequest.Name,
		Description: specialtyRequest.Description,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	return specialtyModel
}
