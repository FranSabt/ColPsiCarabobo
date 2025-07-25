package admin_controller

import (
	"errors"

	db_admin "github.com/FranSabt/ColPsiCarabobo/src/admin/db"
	"github.com/FranSabt/ColPsiCarabobo/src/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func GetAdminById(id uuid.UUID, db *gorm.DB) (models.UserAdmin, error) {
	admin, err := db_admin.GetAdminById(id, db)
	if err != nil {
		return models.UserAdmin{}, err
	}
	if admin.Username == "" {
		return models.UserAdmin{}, errors.New("no admin found")
	}

	return *admin, nil

}
