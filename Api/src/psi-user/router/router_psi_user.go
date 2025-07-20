package psi_user_router

import (
	"github.com/FranSabt/ColPsiCarabobo/db"
	"github.com/FranSabt/ColPsiCarabobo/src/middleware"
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

	group.Post("/login", func(c *fiber.Ctx) error {
		return psiuser_presenter.PsiUserLogin(c, db.DB)
	})

	//==================================================//
	//      RUTAS PRIVADAS (PARA USUARIOS LOGUEADOS)    //
	//==================================================//
	// Todas las rutas en este grupo requerirán un token válido.

	private := group.Group("/psi-user")                    // Puedes nombrar el grupo como quieras
	private.Use(middleware.ProtectedWithDynamicKey(db.DB)) // Aplicamos el middleware aquí

	// Get My info
	private.Get("/psi-user", func(c *fiber.Ctx) error {
		return psiuser_presenter.GetPsiUserSelfInfo(c, db.DB)
	})
	// Update my info
	private.Put("/psi-user", func(c *fiber.Ctx) error {
		return psiuser_presenter.UpdatePsiUserSelfInfo(c, db.DB)
	})
	// Get self profile pics
	private.Get("/psi-user/user-pic", func(c *fiber.Ctx) error {
		return psiuser_presenter.GetMyProfilePic(c, db)
	})
	// Create/Upload Profile Pic
	private.Post("/psi-user/user-pic", func(c *fiber.Ctx) error {
		return psiuser_presenter.CreatePsiUserImage(c, db)
	})
	// Update Profile Pic
	private.Put("/psi-user/user-pic", func(c *fiber.Ctx) error {
		return psiuser_presenter.UpdatePsiUserImage(c, db)
	})
	// Delete Profile Pic - NOTA: HTTP DELETE es más semántico que PUT para borrar.
	private.Delete("/psi-user/user-pic", func(c *fiber.Ctx) error {
		// Asumo que tienes un manejador para borrar, si no, puedes crearlo.
		// return psiuser_presenter.DeletePsiUserImage(c, db)
		return c.SendString("DELETE /psi-user/user-pic (protected)")
	})

	//==================================================//
	//       RUTAS DE ADMINISTRADOR (AÚN MÁS PRIVADAS)   //
	//==================================================//
	// Estas rutas también están en 'group', por lo que son públicas.
	// Deberían tener su propio grupo y, idealmente, un middleware adicional
	// que verifique el rol de "admin" en el token.

	admin := group.Group("/admin")
	// Primero, requiere que el usuario esté logueado.
	// admin.Use(middleware.ProtectedWithDynamicKey(db.DB))
	// Segundo, requiere que el usuario tenga el rol de admin.
	admin.Use(middleware.ProtectedAdminWithDynamicKey(db.DB)) // <-- Deberías crear este middleware

	// create PSI
	admin.Post("/psi-user", func(c *fiber.Ctx) error {
		return psi_user_admin_presenter.AdminCreatePsiUser(c, db.DB)
	})
	// upload CSV
	admin.Post("/upload-csv", func(c *fiber.Ctx) error {
		return psi_user_admin_presenter.UploadCsv(c, db.DB)
	})
	// Get user list as admin
	admin.Get("/psi-user-list", func(c *fiber.Ctx) error {
		return psi_user_admin_presenter.AdminGetPsiUserList(c, db.DB)
	})
	// Get specific user by ID as admin
	admin.Post("/psi-user-by-id", func(c *fiber.Ctx) error {
		return psi_user_admin_presenter.GetPsiUsersByID(c, db.DB)
	})
	// Patch user by ID as admin
	admin.Patch("/psi-user-by-id", func(c *fiber.Ctx) error {
		return psi_user_admin_presenter.PatchPsiUserByID(c, db.DB)
	})
}
