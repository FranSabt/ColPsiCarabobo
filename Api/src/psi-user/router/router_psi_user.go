package psi_user_router

import (
	"github.com/FranSabt/ColPsiCarabobo/db"
	psi_user_admin_presenter "github.com/FranSabt/ColPsiCarabobo/src/psi-user/admin"
	psiuser_presenter "github.com/FranSabt/ColPsiCarabobo/src/psi-user/psi-user"
	"github.com/gofiber/fiber/v2"
)

func PsiUserRouter(group fiber.Router, db db.StructDb) {
	group.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Psi User")
	})

	//* ---- PUBLIC ROUTES ---- *//
	// get all /public
	group.Get("/get-all", func(c *fiber.Ctx) error {
		return psiuser_presenter.GetPsiUsers(c, db.DB)
	})
	// get details /public
	group.Post("/get-by-id", func(c *fiber.Ctx) error {
		return psiuser_presenter.GetPsiUserById(c, db.DB)
	})

	//* ---- PRIVATE ROUTES ---- *//
	// Login
	group.Post("/login", func(c *fiber.Ctx) error {
		return psiuser_presenter.PsiUserLogin(c, db.DB)
	})
	// Get My info
	group.Get("/psi-user", func(c *fiber.Ctx) error {
		return psiuser_presenter.GetPsiUserSelfInfo(c, db.DB)
	})
	// Update my info
	group.Put("/psi-user", func(c *fiber.Ctx) error {
		return psiuser_presenter.UpdatePsiUserSelfInfo(c, db.DB)
	})
	// Get self profile pics
	group.Get("/psi-user/user-pic", func(c *fiber.Ctx) error {
		return psiuser_presenter.GetMyProfilePic(c, db)
	})
	// Update Profile Pic
	group.Post("/psi-user/user-pic", func(c *fiber.Ctx) error {
		return psiuser_presenter.CreatePsiUserImage(c, db)
	})
	// update Profile Pic
	group.Put("/psi-user/user-pic", func(c *fiber.Ctx) error {
		return psiuser_presenter.UpdatePsiUserImage(c, db)
	})
	// Delete Profile Pic
	group.Put("/psi-user/user-pic", func(c *fiber.Ctx) error {
		return psiuser_presenter.UpdatePsiUserImage(c, db)
	})

	//* ADMIND ROUTES **/
	// create PSI
	// update PSI
	// create massive
	group.Post("/upload-csv", func(c *fiber.Ctx) error {
		return psi_user_admin_presenter.UploadCsv(c, db.DB)
	})

	// TODO: minimisar la cantidad de informacion enviada en la respuesta
	group.Get("/psi-user", func(c *fiber.Ctx) error {
		return psi_user_admin_presenter.AdminGetPsiUserList(c, db.DB)
	})

	group.Post("/psi-user", func(c *fiber.Ctx) error {
		return psi_user_admin_presenter.AdminCreatePsiUser(c, db.DB)
	})

	group.Post("/psi-user-by-id", func(c *fiber.Ctx) error {
		return psi_user_admin_presenter.GetPsiUsersByID(c, db.DB)
	})

	group.Patch("/psi-user-by-id", func(c *fiber.Ctx) error {
		return psi_user_admin_presenter.PatchPsiUserByID(c, db.DB)
	})
}
