package router

import (
	"github.com/FranSabt/ColPsiCarabobo/db"
	admin_router "github.com/FranSabt/ColPsiCarabobo/src/admin/router"
	psi_user_routes "github.com/FranSabt/ColPsiCarabobo/src/psi-user/router"
	specialties_router "github.com/FranSabt/ColPsiCarabobo/src/specialties/router"
	"github.com/gofiber/fiber/v2"
)

func Router(group fiber.Router, db db.StructDb) {
	psiUser := group.Group("/psi-user")
	specialties := group.Group("/specialties")
	admin := group.Group("/admin")

	psi_user_routes.PsiUserRouter(psiUser, db)
	specialties_router.SpecialtiesRouter(specialties, db)
	admin_router.AdminRouter(admin, db)

}
