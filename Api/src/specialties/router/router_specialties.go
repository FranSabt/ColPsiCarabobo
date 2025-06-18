package specialties_router

import (
	"github.com/FranSabt/ColPsiCarabobo/db"
	"github.com/gofiber/fiber/v2"
)

func SpecialtiesRouter(group fiber.Router, db db.StructDb) {
	group.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Specialties")
	})
}
