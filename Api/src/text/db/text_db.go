package text_db // O el nombre de tu paquete de base de datos

import (
	"errors"

	"github.com/FranSabt/ColPsiCarabobo/src/models"
	"gorm.io/gorm"
)

// --- CREATE ---

// CreateTextDb inserta un nuevo TextModel en la base de datos.
// Devuelve el TextModel creado (con su ID) y un error si algo falla.
func CreateTextDb(db *gorm.DB, textData models.TextModel) (*models.TextModel, error) {
	// GORM llenará automáticamente el ID, CreatedAt y UpdatedAt en el struct 'textData'
	result := db.Create(&textData)
	if result.Error != nil {
		return nil, result.Error
	}
	// Devolvemos el struct completo con los campos actualizados
	return &textData, nil
}

// --- READ ---

// GetTextByIDDb busca un TextModel por su ID.
// Devuelve el TextModel encontrado o un error si no existe.
func GetTextByIDDb(db *gorm.DB, id uint) (*models.TextModel, error) {
	var text models.TextModel

	// Usamos First, que devuelve un error gorm.ErrRecordNotFound si no lo encuentra.
	result := db.First(&text, id)

	if result.Error != nil {
		// Verificamos específicamente si el error es "registro no encontrado"
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			// Devolvemos un error personalizado y más limpio para el servicio/handler
			return nil, errors.New("text not found")
		}
		// Para cualquier otro error de base de datos, lo devolvemos tal cual
		return nil, result.Error
	}

	return &text, nil
}

// --- UPDATE ---

// UpdateTextDb actualiza un TextModel existente por su ID.
// Recibe el ID y un mapa de los campos a actualizar.
// Devuelve un error si algo falla.
func UpdateTextDb(db *gorm.DB, id uint, updateData map[string]interface{}) error {
	// Primero, verificamos que el registro exista para evitar actualizaciones "fantasma"
	var text models.TextModel
	if err := db.First(&text, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("text not found, cannot update")
		}
		return err
	}

	// Usamos Model(&models.TextModel{}).Where("id = ?", id) para actualizar
	// de forma segura por ID.
	result := db.Model(&models.TextModel{}).Where("id = ?", id).Updates(updateData)
	if result.Error != nil {
		return result.Error
	}

	// Opcional: Verificar si realmente se actualizó algo
	if result.RowsAffected == 0 {
		// Esto puede pasar si los datos enviados son idénticos a los existentes
		// o si el registro fue borrado por otro proceso entre el First y el Updates.
		// En general, no es un error, así que podemos ignorarlo o loguearlo.
	}

	return nil
}

// --- DELETE ---

// DeleteTextDb elimina permanentemente un TextModel de la base de datos por su ID.
// Devuelve un error si algo falla o si el registro no existía.
func DeleteTextDb(db *gorm.DB, id uint) error {
	// Este modelo no tiene gorm.DeletedAt, por lo que `Delete` hará un borrado físico.
	result := db.Delete(&models.TextModel{}, id)

	if result.Error != nil {
		return result.Error
	}

	// Si `RowsAffected` es 0, significa que no se encontró un registro con ese ID.
	// Es bueno notificar esto para que el handler pueda devolver un 404.
	if result.RowsAffected == 0 {
		return errors.New("text not found, cannot delete")
	}

	return nil
}
