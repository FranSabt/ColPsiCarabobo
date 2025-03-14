package config

import (
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
)

func ResponseLogger(c *fiber.Ctx) error {
	// Continuar con el siguiente middleware o manejador de ruta
	err := c.Next()

	// Registrar la información de la respuesta
	log.Printf(
		"[%s] %s %s - %d",
		time.Now().Format("2006-01-02 15:04:05"),
		c.Method(),
		c.Path(),
		c.Response().StatusCode(),
	)

	return err
}
