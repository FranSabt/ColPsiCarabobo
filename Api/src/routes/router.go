package router

import (
	"github.com/FranSabt/ColPsiCarabobo/db"
	psi_user_routes "github.com/FranSabt/ColPsiCarabobo/src/psi-user/router"
	"github.com/gofiber/fiber/v2"
)

func Router(group fiber.Router, db db.StructDb) {
	psiUser := group.Group("/psi-user")
	psi_user_routes.PsiUserRouter(psiUser, db)

}
