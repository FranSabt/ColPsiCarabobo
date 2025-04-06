package psi_user_controller

import (
	"fmt"

	psi_user_db "github.com/FranSabt/ColPsiCarabobo/src/psi-user/db"
	psi_user_request "github.com/FranSabt/ColPsiCarabobo/src/psi-user/request-structs"
	"gorm.io/gorm"
)

func CheckPsiUserUniqueFields(db *gorm.DB, request psi_user_request.PsiUserCreateRequest) (bool, [4]string, error) {

	var flag = true
	var conflicts = [4]string{}

	exist, err := psi_user_db.CheckIfExistPsiUser(db, "fpv", fmt.Sprint(request.FPV))
	if err != nil {
		return false, conflicts, err
	}
	if exist {
		flag = false
		conflicts[0] = "FPV already in use"
	}

	exist, err = psi_user_db.CheckIfExistPsiUser(db, "ci", fmt.Sprint(request.CI))
	if err != nil {
		return false, conflicts, err
	}
	if exist {
		flag = false
		conflicts[1] = "CI already in use"
	}

	exist, err = psi_user_db.CheckIfExistPsiUser(db, "username", request.Username)
	if err != nil {
		return false, conflicts, err
	}
	if exist {
		flag = false
		conflicts[2] = "Username already in use"
	}

	exist, err = psi_user_db.CheckIfExistPsiUser(db, "email", request.Email)
	if err != nil {
		flag = false
		return false, conflicts, err
	}
	if exist {
		flag = false
		conflicts[3] = "Username already in use"
	}

	return flag, conflicts, nil
}
