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
	database_image, err := db.ConnectImage()
	if err != nil {
		log.Fatalf("Could not connect to the images database: %v", err)
	}
	log.Println("Successfully connected to the images database!")

	Db_Conteiner := db.StructDb{
		Image: database_image,
		DB:    database,
	}
	// Root route
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, ColPsiCarabobo!")
	})

	// Register routes from the router package
	api := app.Group("/api")
	v1 := api.Group("/v1")
	router.Router(v1, Db_Conteiner)

	// 404 Handler
	app.Use(func(c *fiber.Ctx) error {
		return c.SendStatus(404) // => 404 "Not Found"
	})

	// Start the server
	app.Listen(":5000")
}
