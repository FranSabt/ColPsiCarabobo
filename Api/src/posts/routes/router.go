package posts_routes

import (
	"github.com/FranSabt/ColPsiCarabobo/db"
	"github.com/FranSabt/ColPsiCarabobo/src/middleware"
	post_presenter "github.com/FranSabt/ColPsiCarabobo/src/posts/presenter"
	"github.com/gofiber/fiber/v2"
)

func PostRouter(group fiber.Router, db db.StructDb) {
	group.Get("/test-post", func(c *fiber.Ctx) error {
		return c.SendString("Post")
	})

	group.Get("/", func(c *fiber.Ctx) error {
		return post_presenter.GetPosts(c, db.DB)
	})

	group.Get("/get-text", func(c *fiber.Ctx) error {
		return post_presenter.GetPosts(c, db.DB)
	})

	admin := group.Group("/admin")
	admin.Use(middleware.ProtectedAdminWithDynamicKey(db.DB))

	admin.Post("/", func(c *fiber.Ctx) error {
		return post_presenter.CreatePostAdmin(c, db.DB, db.Text)
	})

	admin.Put("/", func(c *fiber.Ctx) error {
		return post_presenter.UpdatePost(c, db.DB, db.Text)
	})

}
