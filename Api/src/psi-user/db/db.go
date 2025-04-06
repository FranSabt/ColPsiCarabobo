package psi_user_db

import (
	"fmt"
	"log"

	"github.com/FranSabt/ColPsiCarabobo/src/models"
	"gorm.io/gorm"
)

func CreatePsiUseDb(db *gorm.DB, psiUserModel models.PsiUserModel) error {
	// Intentar crear el registro en la base de datos
	result := db.Create(&psiUserModel)
	if result.Error != nil {
		// Si hay un error, lo retornamos
		return result.Error
	}

	// Si todo está bien, retornamos nil (sin error)
	return nil
}

func CreatePsiColDataDb(db *gorm.DB, psiUserColData models.PsiUserColData) error {
	// Intentar crear el registro en la base de datos
	result := db.Create(&psiUserColData)
	if result.Error != nil {
		// Si hay un error, lo retornamos
		return result.Error
	}

	// Si todo está bien, retornamos nil (sin error)
	return nil
}

func CreatePsiUseDb2(db *gorm.DB, psiUserModel *models.PsiUserModel) error {
	// Intentar crear el registro en la base de datos
	result := db.Create(&psiUserModel)
	if result.Error != nil {
		// Si hay un error, lo retornamos
		return result.Error
	}

	// Si todo está bien, retornamos nil (sin error)
	return nil
}

func CreatePsiColDataDb2(db *gorm.DB, psiUserColData *models.PsiUserColData) error {
	// Intentar crear el registro en la base de datos
	result := db.Create(&psiUserColData)
	if result.Error != nil {
		// Si hay un error, lo retornamos
		return result.Error
	}

	// Si todo está bien, retornamos nil (sin error)
	return nil
}

type PsiUserResponse struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	FPV         int    `json:"fpv"`
	CI          int    `json:"ci"`
	Nationality string `json:"nationality"`
}

func GetPaginatedPsiUsers(db *gorm.DB, page int, pageSize int, ci *int, fpv *int) ([]PsiUserResponse, int64, error) {
	var psiUsers []PsiUserResponse
	var totalRecords int64

	// Calcular el offset para la paginación
	offset := (page - 1) * pageSize

	// Crear la consulta base
	query := db.Model(&models.PsiUserModel{}).
		Select("id, first_name || ' ' || last_name as name, fpv, ci, nationality")

	// Aplicar filtros si se proporcionan
	if ci != nil {
		query = query.Where("ci = ?", *ci)
	}
	if fpv != nil {
		query = query.Where("fpv = ?", *fpv)
	}

	// Contar el total de registros (con los filtros aplicados)
	if err := query.Count(&totalRecords).Error; err != nil {
		return nil, 0, err
	}

	// Aplicar paginación y obtener los registros
	if err := query.Offset(offset).Limit(pageSize).Find(&psiUsers).Error; err != nil {
		return nil, 0, err
	}

	return psiUsers, totalRecords, nil
}

func CheckIfExistPsiUser(db *gorm.DB, column, value string) (bool, error) {
	var count int64

	// Ejecutar la consulta contando registros que coincidan
	err := db.Model(&models.PsiUserModel{}).
		Where(column+" = ?", value).
		Count(&count).
		Error

	if err != nil {
		log.Printf("Error searching in database: %v", err)
		return false, fmt.Errorf("database error: %v", err)
	}

	// Si count > 0, el usuario existe
	return count > 0, nil
}
