package psi_user_router

import (
	psi_user_admin_presenter "github.com/FranSabt/ColPsiCarabobo/src/psi-user/admin"
	psiuser_presenter "github.com/FranSabt/ColPsiCarabobo/src/psi-user/psi-user"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func PsiUserRouter(group fiber.Router, db *gorm.DB) {
	group.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Psi User")
	})

	//* ---- PUBLIC ROUTES ---- *//
	// get all /public
	group.Get("/get-all", func(c *fiber.Ctx) error {
		return psiuser_presenter.GetPsiUsers(c, db)
	})
	// get details /public
	group.Post("/get-by-id", func(c *fiber.Ctx) error {
		return psiuser_presenter.GetPsiUserById(c, db)
	})

	//* ---- PRIVATE ROUTES ---- *//
	// Login
	group.Post("/login", func(c *fiber.Ctx) error {
		return psiuser_presenter.PsiUserLogin(c, db)

	})
	// Get My info
	group.Get("/psi-user", func(c *fiber.Ctx) error {
		return psiuser_presenter.GetPsiUserSelfInfo(c, db)

	})
	// Update my info
	// Update Profile Pic

	// get details self
	// update self

	//* ADMIND ROUTES **/
	// create PSI
	// update PSI
	// create massive
	group.Post("/upload-csv", func(c *fiber.Ctx) error {
		return psi_user_admin_presenter.UploadCsv(c, db)
	})

	// TODO: minimisar la cantidad de informacion enviada en la respuesta
	group.Get("/psi-user", func(c *fiber.Ctx) error {
		return psi_user_admin_presenter.AdminGetPsiUserList(c, db)
	})

	group.Post("/psi-user", func(c *fiber.Ctx) error {
		return psi_user_admin_presenter.AdminCreatePsiUser(c, db)
	})

	group.Post("/psi-user-by-id", func(c *fiber.Ctx) error {
		return psi_user_admin_presenter.GetPsiUsersByID(c, db)
	})

	group.Patch("/psi-user-by-id", func(c *fiber.Ctx) error {
		return psi_user_admin_presenter.PatchPsiUserByID(c, db)
	})
}
