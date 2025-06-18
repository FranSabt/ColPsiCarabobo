package admin_router

import (
	"github.com/FranSabt/ColPsiCarabobo/db"
	"github.com/FranSabt/ColPsiCarabobo/src/admin/admin"
	"github.com/gofiber/fiber/v2"
)

func AdminRouter(group fiber.Router, db db.StructDb) {
	group.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Admin")
	})

	group.Get("/get-all", func(c *fiber.Ctx) error {
		return admin.GetAdmins(c, db.DB)
	})
}
