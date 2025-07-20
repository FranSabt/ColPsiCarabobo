package main

import (
	"log"
	"time"

	"github.com/FranSabt/ColPsiCarabobo/config"
	"github.com/FranSabt/ColPsiCarabobo/db"
	router "github.com/FranSabt/ColPsiCarabobo/src/routes"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// connectWithRetries es una función auxiliar para manejar la conexión a la base de datos con reintentos.
// Toma como argumentos:
// - connectFunc: La función de conexión real (ej. db.Connect, db.ConnectImage).
// - dbName: Un nombre descriptivo de la base de datos para los logs.
// - retries: El número de reintentos a intentar.
// - delay: El tiempo de espera entre reintentos.
func connectWithRetries(connectFunc func() (*gorm.DB, error), dbName string, retries int, delay time.Duration) *gorm.DB {
	var database *gorm.DB
	var err error

	for i := 0; i < retries; i++ {
		database, err = connectFunc()
		if err == nil {
			log.Printf("Successfully connected to the %s database!", dbName)
			return database
		}
		log.Printf("Failed to connect to %s database, retrying in %v... (Attempt %d/%d). Error: %v", dbName, delay, i+1, retries, err)
		time.Sleep(delay)
	}

	// Si después de todos los reintentos no se pudo conectar, el programa se detiene.
	log.Fatalf("Could not connect to the %s database after %d attempts. Last error: %v", dbName, retries, err)
	return nil // Esta línea nunca se alcanzará debido a log.Fatalf
}

func main() {
	app := fiber.New()

	// Logger
	app.Use(config.ResponseLogger)

	// Cargar las variables de entorno (esto ya lo tenías bien)
	adminUsername := config.EnvConfig("ADMIN_USERNAME")
	log.Println("ADMIN_USERNAME: ", adminUsername)

	// --- Conexión Robusta a las Bases de Datos ---

	// Parámetros para la lógica de reintentos
	const (
		maxRetries = 5
		retryDelay = 5 * time.Second
	)

	// Conectar a la base de datos principal
	database := connectWithRetries(db.Connect, "main", maxRetries, retryDelay)

	// Conectar a la base de datos de imágenes
	database_image := connectWithRetries(db.ConnectImage, "images", maxRetries, retryDelay)

	// Conectar a la base de datos de texto
	// NOTA: He corregido los mensajes de log que tenías con copy-paste.
	database_text := connectWithRetries(db.ConnectText, "text", maxRetries, retryDelay)

	// Crear el contenedor de conexiones de BD
	Db_Conteiner := db.StructDb{
		Image: database_image,
		DB:    database,
		Text:  database_text,
	}

	// --- Configuración de Rutas y Servidor ---

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
	log.Println("Starting server on port :5000")
	if err := app.Listen(":5000"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
