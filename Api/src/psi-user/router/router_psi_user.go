package psi_user_router

import (
	psi_user_admin_presenter "github.com/FranSabt/ColPsiCarabobo/src/psi-user/admin"
	"github.com/gofiber/fiber/v2"
)

func PsiUserRouter(group fiber.Router) {
	group.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Psi User")
	})

	// get all /public
	// get details /public

	// get details self
	// update self

	//* ADMIND ROUTES **/
	// create PSI
	// update PSI
	// create massive
	group.Post("/upload-csv", func(c *fiber.Ctx) error {
		return psi_user_admin_presenter.UploadCsv(c)
	})

}
