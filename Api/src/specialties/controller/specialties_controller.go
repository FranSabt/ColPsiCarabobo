package specialties_controller

import (
	"errors"
	"regexp"

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

func GetPsiSpecialtiesCountController(db *gorm.DB) (int64, int64, error) {
	count, last_id, err := specialties_db.CountAndGetLastSpecialtyID(db)
	if err != nil {
		return 0, 0, err
	}

	return count, last_id, nil
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

// ////////////////////////////////////////////////////
// ////////////////////////////////////////////////////
// ////////////////////////////////////////////////////

func UpdatePsiSpecialtyController(request *specialties_structs.SpecialtyUpdate, db *gorm.DB) error {
	// 1. Validar los datos de entrada primero. Si fallan, no tocamos la BD.
	if err := checkUpdateFieldSpecialty(request); err != nil {
		return err // Error de validación (ej. 400 Bad Request)
	}

	// 2. Recuperar el modelo existente de la base de datos.
	//    Esto también confirma que el 'id' es válido y el registro está activo.
	id := uint(request.ID)
	model_to_update, err := specialties_db.GetSpecialtyById(db, id)
	if err != nil {
		return err // Error de "no encontrado" (ej. 404 Not Found) o error de BD (500)
	}

	// 3. Aplicar los cambios al modelo recuperado (si se proporcionaron).
	if request.Name != "" {
		model_to_update.Name = request.Name
	}

	if request.Description != "" {
		model_to_update.Description = request.Description
	}

	// 4. ¡PASO CRÍTICO FALTANTE! Guardar el modelo actualizado en la base de datos.
	//    Usamos la función `UpdateSpecialty` que corregimos antes.
	if err := specialties_db.UpdateSpecialty(db, model_to_update); err != nil {
		return err // Error durante la operación de guardado (ej. 500 Internal Server Error)
	}

	// 5. Si todo fue exitoso, devolver nil.
	return nil
}

//////////////////////////////////////////////////////
//////////////////////////////////////////////////////
//////////////////////////////////////////////////////

func DeleteSpecialtyController(id int64, db *gorm.DB) error {
	if id <= 0 {
		return errors.New("invalid id")
	}

	err := specialties_db.DeleteSpecialty(db, uint(id))
	if err != nil {
		return err
	}

	return nil
}

//////////////////////////////////////////////////////
//////////////////////////////////////////////////////
//////////////////////////////////////////////////////

// ------       Auxiliary Functions      ------     //

//////////////////////////////////////////////////////
//////////////////////////////////////////////////////
//////////////////////////////////////////////////////

var unicodeLetterNumberRegex = regexp.MustCompile(`^[\pL\pM\pN\s]+$`)

func checkUpdateFieldSpecialty(request *specialties_structs.SpecialtyUpdate) error {
	// 1. Validar Name si fue proporcionado
	if request.Name != "" { // Comprobamos si el puntero no es nulo
		name := request.Name // Obtenemos el valor
		if len(name) < 4 {
			return errors.New("specialty name must be at least 4 characters long")
		}
		if !unicodeLetterNumberRegex.MatchString(name) {
			return errors.New("specialty name must contain only letters, numbers and spaces")
		}
	}

	// 2. Validar Description si fue proporcionada
	if request.Description != "" {
		description := request.Description
		if len(description) < 10 { // Ejemplo: longitud mínima si se proporciona
			return errors.New("specialty description must be at least 10 characters long")
		}
		if len(description) > 1000 { // Ejemplo: longitud máxima
			return errors.New("specialty description cannot be longer than 1000 characters")
		}
	}

	return nil
}
