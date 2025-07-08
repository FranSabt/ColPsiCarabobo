package specialties_db

import (
	"errors"
	"fmt"
	"time"

	"github.com/FranSabt/ColPsiCarabobo/src/models"
	specialties_structs "github.com/FranSabt/ColPsiCarabobo/src/specialties/request-structs"
	"gorm.io/gorm"
)

func CreateSpecialty(db *gorm.DB, specialty models.PsiSpecialty) error {
	// Intentar crear el registro en la base de datos
	result := db.Create(&specialty)
	if result.Error != nil {
		// Si hay un error, lo retornamos
		return result.Error
	}

	// Si todo está bien, retornamos nil (sin error)
	return nil
}

//////////////////////////////////////////////////////
//////////////////////////////////////////////////////
//////////////////////////////////////////////////////

func GetPsiSpecialtiesNames(db *gorm.DB) ([]specialties_structs.SpecialtyName, error) {
	// 1. Declarar un slice vacío para almacenar los resultados.
	//    El nombre "specialties" (plural) indica que esperamos múltiples resultados.
	var specialties []specialties_structs.SpecialtyName

	// 2. Usar Model() para especificar la tabla y luego Select() y Find().
	//    GORM mapeará los campos "id" y "name" de la tabla "psi_specialties"
	//    a los campos del struct "SpecialtyName".
	result := db.Model(&models.PsiSpecialty{}).
		Select("id", "name").
		Where("active = ?", true).
		Find(&specialties)

	// 3. Comprobar si hubo un error durante la consulta a la base de datos.
	if result.Error != nil {
		// Si hay un error, devolvemos el slice vacío y el error.
		return nil, result.Error
	}

	// Si no hubo errores, devolvemos el slice con los datos y un error nulo.
	return specialties, nil
}

//////////////////////////////////////////////////////
//////////////////////////////////////////////////////
//////////////////////////////////////////////////////

func CountSpecialties(db *gorm.DB) (int64, error) {
	// 1. Declarar una variable para almacenar el resultado del conteo.
	//    GORM usa int64 para los conteos para evitar desbordamientos en tablas muy grandes.
	var count int64

	// 2. Especificar el modelo para indicar en qué tabla contar y ejecutar Count().
	//    GORM generará una consulta SQL similar a: SELECT COUNT(*) FROM "psi_specialties";
	//    El resultado se guarda en la variable `count` que pasamos por referencia.
	result := db.Model(&models.PsiSpecialty{}).
		Where("active = ?", true).
		Count(&count)

	// 3. Comprobar si hubo un error durante la consulta.
	if result.Error != nil {
		// En caso de error, devolvemos 0 y el error.
		return 0, result.Error
	}

	if count < 1 {
		// En caso de error, devolvemos 0 y el error.
		return 0, errors.New("retorno un valo negativo")
	}

	// 4. Si la consulta fue exitosa, devolvemos el conteo y un error nulo.
	return count, nil
}

//////////////////////////////////////////////////////
//////////////////////////////////////////////////////
//////////////////////////////////////////////////////

func GetSpecialtyDescriptionByID(db *gorm.DB, id uint) (string, error) {
	// 1. Declarar una variable string para almacenar el resultado de la descripción.
	var description string

	// 2. Construir y ejecutar la consulta.
	// - Model(&models.PsiSpecialty{}): Especifica que la consulta es en la tabla "psi_specialties".
	// - Where("id = ?", id): Filtra por el ID proporcionado.
	// - Pluck("description", &description): Extrae el valor de la columna "description"
	//   y lo guarda en la variable `description`.
	result := db.Model(&models.PsiSpecialty{}).
		Where("active = ?", true).
		Where("id = ?", id).
		Pluck("description", &description)

	// 3. Manejar los posibles errores.
	if result.Error != nil {
		// Es una buena práctica comprobar específicamente el error "registro no encontrado".
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			// Devolvemos un string vacío y un error claro indicando que no se encontró el ID.
			return "", fmt.Errorf("specialty with id %d not found", id) // Mensaje mejorado
		}
		// Para cualquier otro error de base de datos (ej. conexión perdida).
		return "", result.Error
	}

	// 4. Si la consulta fue exitosa, el `result.RowsAffected` será 1.
	// El chequeo `result.RowsAffected == 0` es redundante aquí, ya que si no se encuentra
	// un registro, GORM devolverá `gorm.ErrRecordNotFound`, que ya se maneja arriba.
	// Lo podemos eliminar para simplificar.
	/*
		if result.RowsAffected == 0 {
			// Este es otro chequeo de seguridad en caso de que Pluck no devuelva ErrRecordNotFound.
			return "", fmt.Errorf("specialty with id %d not found", id)
		}
	*/

	// 5. Si todo salió bien, devolvemos la descripción encontrada y un error nulo.
	return description, nil
}

//////////////////////////////////////////////////////
//////////////////////////////////////////////////////
//////////////////////////////////////////////////////

func GetSpecialtyById(db *gorm.DB, id uint) (*models.PsiSpecialty, error) {
	var specialty models.PsiSpecialty

	result := db.Model(&models.PsiSpecialty{}).
		Where("active = ?", true).
		First(&specialty, id) // .First() puede tomar la clave primaria directamente como argumento

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			// El mensaje de error es más preciso ahora.
			return nil, fmt.Errorf("specialty with id %d not found or is not active", id)
		}
		// Para cualquier otro error de base de datos.
		return nil, result.Error
	}

	//    Devolvemos un puntero a ella.
	return &specialty, nil
}

//////////////////////////////////////////////////////
//////////////////////////////////////////////////////
//////////////////////////////////////////////////////

func UpdateSpecialty(db *gorm.DB, specialty_model *models.PsiSpecialty) error {
	// `db.Save` actualiza el registro si la clave primaria (ID) no es cero.
	// La variable `specialty_model` ya es un puntero, por lo que se pasa directamente.
	result := db.Save(specialty_model)

	// Comprobamos el campo .Error del resultado para ver si hubo un problema
	// a nivel de base de datos (ej. violación de la restricción 'unique' en 'name').
	if result.Error != nil {
		return result.Error // Devolvemos el error de GORM directamente.
	}

	// Debemos verificar si alguna fila fue realmente actualizada. Si RowsAffected es 0,
	// significa que no se encontró un registro con ese ID.
	if result.RowsAffected == 0 {
		return fmt.Errorf("specialty with id %d not found, no records were updated", specialty_model.ID)
	}

	// Si todo salió bien, devolvemos nil para indicar que no hubo error.
	return nil
}

//////////////////////////////////////////////////////
//////////////////////////////////////////////////////
//////////////////////////////////////////////////////

func DeleteSpecialty(db *gorm.DB, id uint) error {
	// 1. Buscar la especialidad que se va a "borrar".
	//    Esto asegura que exista y esté activa antes de continuar.
	var specialty models.PsiSpecialty
	if err := db.Where("active = ?", true).First(&specialty, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("specialty with id %d not found or is already inactive", id)
		}
		// Para cualquier otro error durante la búsqueda.
		return err
	}

	// 2. Cambiar el estado a inactivo.
	//    Ser explícito con `false` es a menudo más claro que invertir el estado.
	specialty.Active = !specialty.Active
	specialty.UpdatedAt = time.Now()

	// 3. Guardar el cambio en la base de datos.
	//    Pasamos un puntero a `specialty` (`&specialty`).
	//    `Save` actualizará todos los campos, incluido `UpdatedAt` automáticamente.
	if err := db.Save(&specialty).Error; err != nil {
		return fmt.Errorf("failed to delete specialty with id %d: %w", id, err)
	}

	// 4. Si todo salió bien, devolver nil.
	return nil
}
