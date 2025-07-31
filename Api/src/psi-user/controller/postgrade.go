package psi_user_controller

import (
	"time"

	"github.com/FranSabt/ColPsiCarabobo/src/models"
	"github.com/google/uuid"
)

func CreatePsiUserPostGradeModel(psi_user_id uuid.UUID, post_grade_title, post_grade_university, post_grade_graduation_year, post_grade_description string) models.PisUserPostGrade {
	id := uuid.New()
	postgrade := models.PisUserPostGrade{
		ID:        id,
		PsiUserID: psi_user_id,

		PostGradeTitle:          post_grade_title,
		PostGradeUniversity:     post_grade_university,
		PostGradeGraduationYear: post_grade_graduation_year,
		PostGradeDescription:    post_grade_description,

		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	return postgrade
}
