package db_images

import (
	db "github.com/FranSabt/ColPsiCarabobo/db"
	"github.com/FranSabt/ColPsiCarabobo/src/models"
)

func SaveUserImage(image models.ImagesModel) error {
	return db.DB_Images.Db.Create(&image).Error
}
