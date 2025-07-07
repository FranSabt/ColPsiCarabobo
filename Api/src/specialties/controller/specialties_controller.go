package specialties_controller

import (
	"errors"
	// Asegúrate de tener este import si no estaba
	"github.com/FranSabt/ColPsiCarabobo/src/models"
	specialties_db "github.com/FranSabt/ColPsiCarabobo/src/specialties/db"
	specialties_structs "github.com/FranSabt/ColPsiCarabobo/src/specialties/request-structs"
	"gorm.io/gorm"
)

// Convención en Go: Los nombres de las funciones no suelen incluir "Controller".
// Y se usa camelCase, ej: SaveNewSpecialty.
func SaveNewSpecialty(db *gorm.DB, newSpecialty models.PsiSpecialty) error {
	// Validación de longitud del nombre
	if len(newSpecialty.Name) < 4 {
		// Corregido: "character"
		return errors.New("invalid name length: must be at least 4 characters")
	}

	// Validación de longitud de la descripción
	if len(newSpecialty.Description) < 150 {
		return errors.New("invalid description length: must be at least 150 characters")
	}

	// ----- LA CORRECCIÓN CLAVE ESTÁ AQUÍ -----
	// Comprueba si las fechas tienen su valor cero, lo que indica que no se asignaron.
	if newSpecialty.CreatedAt.IsZero() || newSpecialty.UpdatedAt.IsZero() {
		return errors.New("create or update dates must be set")
	}

	err := specialties_db.CreateSpecialty(db, newSpecialty)
	if err != nil {
		return err
	}

	// Si todas las validaciones pasan, no devuelvas ningún error.
	return nil
}

//////////////////////////////////////////////////////
//////////////////////////////////////////////////////
//////////////////////////////////////////////////////

func GetPsiSpecialtiesCountController(db *gorm.DB) (int64, error) {
	count, err := specialties_db.CountSpecialties(db)
	if err != nil {
		return 0, err
	}

	return count, nil
}

//////////////////////////////////////////////////////
//////////////////////////////////////////////////////
//////////////////////////////////////////////////////

func GetPsiSpecialtiesNamesController(db *gorm.DB) ([]specialties_structs.SpecialtyName, error) {
	specialties, err := specialties_db.GetPsiSpecialtiesNames(db)
	if err != nil {
		return nil, err
	}

	return specialties, nil
}

//////////////////////////////////////////////////////
//////////////////////////////////////////////////////
//////////////////////////////////////////////////////

func GetPsiSpecialtiesDescriptionController(db *gorm.DB, id uint) (string, error) {
	description, err := specialties_db.GetSpecialtyDescriptionByID(db, id)
	if err != nil {
		return "", err
	}

	return description, nil
}
