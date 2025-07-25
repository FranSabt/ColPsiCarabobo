package specialties_mapper

import (
	"time"

	"github.com/FranSabt/ColPsiCarabobo/src/models"
	specialties_structs "github.com/FranSabt/ColPsiCarabobo/src/specialties/request-structs"
)

func SpecialtyRequestToSpecialtyModel(specialtyRequest specialties_structs.SpecialtiesRequest, admin models.UserAdmin) models.PsiSpecialty {
	specialtyModel := models.PsiSpecialty{
		Name:        specialtyRequest.Name,
		Description: specialtyRequest.Description,
		Active:      true,
		CreatedAt:   time.Now(),
		CreateBy:    admin.Username,
		CreateById:  admin.ID,
		UpdatedAt:   time.Now(),
		UpdateBy:    admin.Username,
		UpdateById:  admin.ID,
	}

	return specialtyModel
}
