package admin_router

import (
	"github.com/FranSabt/ColPsiCarabobo/db"
	"github.com/FranSabt/ColPsiCarabobo/src/admin/admin"
	"github.com/FranSabt/ColPsiCarabobo/src/middleware"
	"github.com/gofiber/fiber/v2"
)

func AdminRouter(group fiber.Router, db db.StructDb) {
	group.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Admin")
	})

	group.Post("/login", func(c *fiber.Ctx) error {
		return admin.AdminLogin(c, db.DB)
	})

	private := group.Group("/protected")                        // Puedes nombrar el grupo como quieras
	private.Use(middleware.ProtectedAdminWithDynamicKey(db.DB)) // Aplicamos el middleware aquí

	private.Get("/get-all", func(c *fiber.Ctx) error {
		return admin.GetAdmins(c, db.DB)
	})

	private.Post("/admin", func(c *fiber.Ctx) error {
		return admin.CreateOrUpdateAdminHandler(c, db.DB)
	})

	private.Put("/admin", func(c *fiber.Ctx) error {
		return admin.CreateOrUpdateAdminHandler(c, db.DB)
	})

}
