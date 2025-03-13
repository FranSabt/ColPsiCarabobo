package router

import (
	psi_user_routes "github.com/FranSabt/ColPsiCarabobo/src/psi-user/router"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func Router(group fiber.Router, db *gorm.DB) {
	psiUser := group.Group("/psi-user")
	psi_user_routes.PsiUserRouter(psiUser)

}
