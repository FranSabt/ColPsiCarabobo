package psi_user_controller

import (
	"github.com/FranSabt/ColPsiCarabobo/src/models"
	"gorm.io/gorm"
)

// En psi_user_controller.go
func CreateAdminPsiUserSearchQuery(ci, fpv, name string, ci_valid, fpv_valid, name_valid bool, db *gorm.DB) (baseQuery *gorm.DB, countQuery *gorm.DB) {
	// Query base para los datos
	baseQuery = db.Model(&models.PsiUserModel{}).
		Select("id, username, first_name, last_name, ci, fpv, email")

	// Query para el conteo (mismos filtros pero sin SELECT específico)
	countQuery = db.Model(&models.PsiUserModel{})

	// Aplicar filtros comunes a ambas queries
	if ci_valid && ci != "" {
		baseQuery = baseQuery.Where("ci = ?", ci)
		countQuery = countQuery.Where("ci = ?", ci)
	}

	if fpv_valid && fpv != "" {
		baseQuery = baseQuery.Where("fpv = ?", fpv)
		countQuery = countQuery.Where("fpv = ?", fpv)
	}

	if name_valid && name != "" {
		likeName := "%" + name + "%"
		baseQuery = baseQuery.Where("first_name LIKE ? OR last_name LIKE ?", likeName, likeName)
		countQuery = countQuery.Where("first_name LIKE ? OR last_name LIKE ?", likeName, likeName)
	}

	return baseQuery, countQuery
}
