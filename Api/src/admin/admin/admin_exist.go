package admin

import (
	db_admin "github.com/FranSabt/ColPsiCarabobo/src/admin/db"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// AdminExists
// Usa solo para verificar que existe el admin
func AdminExists(uuid uuid.UUID, db *gorm.DB) (bool, error) {
	admin, err := db_admin.GetAdminById(uuid, db)
	if err != nil {
		return false, err
	}

	if admin == nil {
		return false, nil
	}

	return true, nil
}
