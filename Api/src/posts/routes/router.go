package posts_routes

import (
	"github.com/FranSabt/ColPsiCarabobo/db"
	"github.com/gofiber/fiber/v2"
)

func PostRouter(group fiber.Router, db db.StructDb) {
	group.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Post")
	})
}
