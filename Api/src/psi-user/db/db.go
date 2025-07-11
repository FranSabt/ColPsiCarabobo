package psi_user_db

import (
	"errors"
	"fmt"
	"log"

	"github.com/FranSabt/ColPsiCarabobo/src/models"
	"github.com/google/uuid"
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

func GetPaginatedPsiUsers(db *gorm.DB, page, pageSize int, ci *int, fpv *int, name, location, specialty string) ([]PsiUserResponse, int64, error) {
	var psiUsers []PsiUserResponse
	var totalRecords int64

	offset := (page - 1) * pageSize

	// Crear la consulta base. Concatenamos nombres para el campo 'name'.
	// NOTA: La sintaxis de concatenación puede variar entre DBs. '||' es para PostgreSQL/SQLite. Para MySQL usa CONCAT().
	query := db.Model(&models.PsiUserModel{}).
		Select("id, first_name || ' ' || last_name as name, fpv, ci, nationality")

	// --- APLICAR FILTROS ---

	// Filtro por CI (Cédula de Identidad)
	if ci != nil {
		query = query.Where("ci = ?", *ci)
	}

	// Filtro por FPV
	if fpv != nil {
		query = query.Where("fpv = ?", *fpv)
	}

	// Filtro por Nombre (búsqueda parcial, insensible a mayúsculas)
	if name != "" {
		// Busca en el nombre completo concatenado
		fullNameSearch := "%" + name + "%"
		query = query.Where("first_name || ' ' || last_name ILIKE ?", fullNameSearch)
	}

	// Filtro por Ubicación (búsqueda parcial en varios campos)
	if location != "" {
		locationSearch := "%" + location + "%"
		query = query.Where(
			"service_address ILIKE ? OR municipality_carabobo ILIKE ? OR state_outside ILIKE ? OR municipality_outside_carabobo ILIKE ?",
			locationSearch, locationSearch, locationSearch, locationSearch,
		)
	}

	// Filtro por Especialidad (búsqueda parcial en primaria o secundaria)
	if specialty != "" {
		specialtySearch := "%" + specialty + "%"
		query = query.Where(
			"primary_specialty ILIKE ? OR secondary_specialty ILIKE ?",
			specialtySearch, specialtySearch,
		)
	}

	// Contar el total de registros que coinciden con los filtros (antes de paginar)
	if err := query.Count(&totalRecords).Error; err != nil {
		return nil, 0, err
	}

	// Aplicar paginación y ejecutar la consulta para obtener los registros
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

// En psi_user_db.go
func SearchPsiUsersByQuery(db *gorm.DB, baseQuery, countQuery *gorm.DB, pageNum, pageSize int) ([]models.PsiUserModel, int64, error) {
	var users []models.PsiUserModel
	var total int64

	// Ejecutar query de conteo primero
	if err := countQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Aplicar paginación y ejecutar query principal
	offset := (pageNum - 1) * pageSize
	if err := baseQuery.Limit(pageSize).Offset(offset).Find(&users).Error; err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

func GetPsiUserById(db *gorm.DB, id uuid.UUID) (*models.PsiUserModel, error) {
	psiUser := &models.PsiUserModel{}
	err := db.Where("id = ?", id).First(psiUser).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("psi_user record not found")
		}
		return nil, err
	}

	return psiUser, nil
}

func GetPsiUserByIdDetails(db *gorm.DB, id uuid.UUID) (*models.PsiUserModel, *models.PsiUserColData, error) {
	psiUser := &models.PsiUserModel{}
	err := db.Where("id = ?", id).First(psiUser).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil, errors.New("psi_user record not found")
		}
		return nil, nil, err
	}

	psiUserColData := &models.PsiUserColData{}
	err = db.Where("id = ?", psiUser.PsiUserColDataID).First(psiUserColData).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil, errors.New("psi_user_col_data record not found")
		}
		return nil, nil, err
	}

	return psiUser, psiUserColData, nil
}

func SaveUpdatedPsiUserOnly(db *gorm.DB, psiUser *models.PsiUserModel) error {
	// Iniciar transacción
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Guardar PsiUserModel
	if err := tx.Save(psiUser).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Commit si todo fue bien
	return tx.Commit().Error
}

func SaveUpdatedPsiUser(db *gorm.DB, psiUser *models.PsiUserModel, colData *models.PsiUserColData) error {
	// Iniciar transacción
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Guardar PsiUserModel
	if err := tx.Save(psiUser).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Guardar PsiUserColData
	if err := tx.Save(colData).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Commit si todo fue bien
	return tx.Commit().Error
}

func GetPsiUserByUsernameOrEmal(db *gorm.DB, username, query string) (*models.PsiUserModel, error) {
	var psi_user models.PsiUserModel
	err := db.Where(query, username).First(&psi_user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &psi_user, nil
}
