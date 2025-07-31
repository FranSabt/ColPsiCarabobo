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

	// Get self profile pics // publico
	group.Get("/user-pic", func(c *fiber.Ctx) error {
		return psiuser_presenter.GetMyProfilePic(c, db)
	})

	// --- RUTAS PÚBLICAS PARA POSTGRADOS ---
	// Obtener los detalles de un postgrado por su ID
	group.Get("/postgrade/:id", func(c *fiber.Ctx) error {
		return psiuser_presenter.GetPsiUserPostgradeById(c, db.DB)
	})

	// Obtener el archivo de imagen de un postgrado por el ID DE LA IMAGEN
	group.Get("/postgrade/image/:id", func(c *fiber.Ctx) error {
		return psiuser_presenter.GetPostgradeImageById(c, db)
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

	// --- RUTAS PARA FOTOS DE PERFIL ---
	// Create/Upload Profile Pic
	private.Post("/user-pic", func(c *fiber.Ctx) error {
		return psiuser_presenter.CreatePsiUserImage(c, db)
	})
	// Update Profile Pic
	private.Put("/user-pic", func(c *fiber.Ctx) error {
		return psiuser_presenter.UpdatePsiUserImage(c, db)
	})
	// Delete Profile Pic
	private.Delete("/user-pic", func(c *fiber.Ctx) error {
		// return psiuser_presenter.DeletePsiUserImage(c, db)
		return c.SendString("DELETE /psi-user/user-pic (protected)")
	})

	//==================================================//
	//       RUTAS PARA POSTGRADOS (PRIVADAS)           //
	//==================================================//
	// Un usuario logueado puede gestionar sus propios postgrados.

	// Crear un nuevo postgrado
	private.Post("/postgrade", func(c *fiber.Ctx) error {
		// Asumiendo que CreatePsiUserPotgradeRefactored es el nombre final de la función de creación.
		return psiuser_presenter.CreatePsiUserPotgradeRefactored(c, db)
	})

	// Actualizar un postgrado existente por su ID
	// El :id es un parámetro de ruta que capturaremos en el handler con c.Params("id")
	private.Put("/postgrade/:id", func(c *fiber.Ctx) error {
		return psiuser_presenter.UpdatePsiUserPostgrade(c, db)
	})

	// Borrar (soft delete) un postgrado por su ID
	private.Delete("/postgrade/:id", func(c *fiber.Ctx) error {
		return psiuser_presenter.SoftDeletePsiUserPostgrade(c, db)
	})

	//==================================================//
	//       RUTAS DE ADMINISTRADOR (AÚN MÁS PRIVADAS)   //
	//==================================================//

	admin := group.Group("/admin")
	admin.Use(middleware.ProtectedAdminWithDynamicKey(db.DB))

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
