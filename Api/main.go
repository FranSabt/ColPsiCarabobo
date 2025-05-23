package main

import (
	"log"

	"github.com/FranSabt/ColPsiCarabobo/config"
	"github.com/FranSabt/ColPsiCarabobo/db"
	router "github.com/FranSabt/ColPsiCarabobo/src/routes"
	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	// logger
	app.Use(config.ResponseLogger)

	// Cargar las variables de entorno desde el archivo .env
	adminUsername := config.EnvConfig("ADMIN_USERNAME")
	log.Println("ADMIN_USERNAME: ", adminUsername)

	// Conectar a la base de datos
	database, err := db.Connect()
	if err != nil {
		log.Fatalf("Could not connect to the database: %v", err)
	}
	log.Println("Successfully connected to the database!")

	// Conectar a la base de datos
	_, err = db.ConnectImage()
	if err != nil {
		log.Fatalf("Could not connect to the images database: %v", err)
	}
	log.Println("Successfully connected to the images database!")

	// Root route
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, ColPsiCarabobo!")
	})

	// Register routes from the router package
	api := app.Group("/api")
	v1 := api.Group("/v1")
	router.Router(v1, database)

	// Start the server
	app.Listen(":5000")
}
