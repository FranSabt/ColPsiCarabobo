package specialties_router

import (
	"github.com/FranSabt/ColPsiCarabobo/db"
	"github.com/FranSabt/ColPsiCarabobo/src/middleware"
	"github.com/FranSabt/ColPsiCarabobo/src/specialties/specialties"
	"github.com/gofiber/fiber/v2"
)

func SpecialtiesRouter(group fiber.Router, db db.StructDb) {
	group.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Specialties")
	})

	// --- RUTAS DE CONSULTA (GET) ---

	// Ruta: GET /specialties/count
	// Devuelve el número total de especialidades.
	group.Get("/count", func(c *fiber.Ctx) error {
		// Usamos un cierre para pasar la conexión 'db.Db' (que es *gorm.DB) al handler.
		return specialties.GetPsiSpecialtiesCount(c, db.DB)
	})

	// Ruta: GET /specialties/names
	// Devuelve una lista de todas las especialidades con su ID y nombre.
	group.Get("/names", func(c *fiber.Ctx) error {
		return specialties.GetPsiSpecialtiesNames(c, db.DB)
	})

	// Ruta: GET /specialties/:id/description
	// Devuelve la descripción de una especialidad específica.
	group.Get("/:id/description", func(c *fiber.Ctx) error {
		return specialties.GetPsiSpecialtiesDescription(c, db.DB)
	})

	admin := group.Group("")
	admin.Use(middleware.ProtectedWithDynamicKey(db.DB))

	admin.Put("/:id", func(c *fiber.Ctx) error {
		return specialties.UpdatePsiSepecialty(c, db.DB)
	})

	admin.Delete("/:id", func(c *fiber.Ctx) error {
		return specialties.DeletePsiSpecialty(c, db.DB)
	})

	// --- RUTA DE CREACIÓN (POST) ---

	// Ruta: POST /specialties/
	// Crea una nueva especialidad.
	admin.Post("/", func(c *fiber.Ctx) error {
		return specialties.CreatePsiSpecialty(c, db.DB)
	})
}
