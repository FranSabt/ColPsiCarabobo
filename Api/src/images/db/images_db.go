package db_images

import (
	"errors"
	"fmt"
	"time"

	"github.com/FranSabt/ColPsiCarabobo/src/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func SaveUserImage(image models.ProfilePicModel, db *gorm.DB) error {
	return db.Create(&image).Error
}

func GetImageByIDWithAssociations(id string, db *gorm.DB) (models.ProfilePicModel, error) {
	var image models.ProfilePicModel
	result := db.Preload("RelatedModel").Where("id = ?", id).First(&image)
	return image, result.Error
}

func GetFirstImageByAssociatedID(associatedID uuid.UUID, db *gorm.DB) (*models.ProfilePicModel, error) {
	var image models.ProfilePicModel
	result := db.Where("user_id = ?", associatedID).First(&image)

	if result.Error != nil {
		return nil, result.Error
	}
	return &image, nil
}

func GetAllUserProfilePicsID(associatedID uuid.UUID, db *gorm.DB) ([]models.ProfilePicModel, error) {
	var images []models.ProfilePicModel

	result := db.Where("user_id = ?", associatedID).Find(&images)

	if result.Error != nil {
		return nil, result.Error
	}

	return images, nil
}

func CheckProfilePicLimit(associatedID uuid.UUID, db *gorm.DB) (int64, error) {
	var count int64
	result := db.Model(&models.ProfilePicModel{}).
		Where("user_id = ?", associatedID).
		Count(&count)

	if result.Error != nil {
		return -1, result.Error
	}

	return count, nil
}

func GetImageById(imageID uuid.UUID, db *gorm.DB) (*models.ProfilePicModel, error) {
	var image models.ProfilePicModel
	result := db.Where("id = ?", imageID).First(&image) // Usar First() y pasar &image

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil // O manejar como prefieras
		}
		return nil, result.Error
	}

	return &image, nil
}

func UpdateImageById(imageID, user_uuid uuid.UUID, newImageData []byte, format string, db *gorm.DB, filename, username string) (bool, error) {
	// Opción 1: Update directo (más eficiente)
	result := db.Model(&models.ProfilePicModel{}).
		Where("id = ?", imageID).
		Updates(map[string]interface{}{
			"image_data":   newImageData,
			"format":       format,
			"updated_at":   time.Now(),
			"name":         filename,
			"update_by":    username,
			"update_by_id": user_uuid,
		})

	if result.Error != nil {
		return false, result.Error
	}

	// Verificar si se actualizó algún registro
	if result.RowsAffected == 0 {
		return false, fmt.Errorf("no image found with id: %s", imageID)
	}

	return true, nil
}

func DeleteImageById(imageID uuid.UUID, db *gorm.DB) (bool, error) {
	result := db.Where("id = ?", imageID).Delete(&models.ProfilePicModel{})

	if result.Error != nil {
		return false, result.Error
	}

	// Verificar si se eliminó algún registro
	if result.RowsAffected == 0 {
		return false, fmt.Errorf("no image found with id: %s", imageID)
	}

	return true, nil
}
