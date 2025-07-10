package middleware

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/FranSabt/ColPsiCarabobo/src/models"
	psi_user_db "github.com/FranSabt/ColPsiCarabobo/src/psi-user/db"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Protected protect routes
// func Protected() fiber.Handler {
// 	return jwtware.New(jwtware.Config{
// 		SigningKey:   jwtware.SigningKey{Key: []byte(config.EnvConfig("SECRET"))},
// 		ErrorHandler: jwtError,
// 	})
// }

// jwtError es un helper para crear respuestas de error JSON consistentes.
func jwtError(c *fiber.Ctx, status int, message string) error {
	return c.Status(status).JSON(fiber.Map{
		"status":  "error",
		"message": message,
	})
}

func ProtectedWithDynamicKey(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// 1. Extraer el token del header "Authorization: Bearer <token>"
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return jwtError(c, fiber.StatusUnauthorized, "Missing or malformed JWT")
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return jwtError(c, fiber.StatusUnauthorized, "Missing or malformed JWT")
		}
		tokenString := parts[1]

		// 2. Parsear el token usando una "Keyfunc". Esta función mágica
		// se encargará de buscar la clave dinámicamente.
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Primero, verificar que el método de firma es el que esperas (HS256)
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}

			// 3. Extraer los claims para obtener el user_id
			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				return nil, errors.New("invalid token claims")
			}

			// En tu función de login usas "user_id", así que lo leemos aquí
			userID, ok := claims["user_id"].(string)
			if !ok {
				return nil, errors.New("user_id not found in token")
			}

			uuid_parsed, err := uuid.Parse(userID) // <-- CORRECCIÓN 1: Usar 'userID'
			if err != nil {
				// Si el ID en el token no es un UUID válido, es un token inválido.
				log.Printf("Invalid UUID format in token for user_id: %s", userID)
				return nil, errors.New("invalid user identifier in token")
			}

			// 4b. Ahora que tenemos un UUID válido, buscar en la base de datos.
			//    Usamos '=' en lugar de ':=' para asignar a la variable 'err' ya existente.
			var psiUser *models.PsiUserModel // Asegúrate de que el tipo sea el correcto
			psiUser, err = psi_user_db.GetPsiUserById(db, uuid_parsed)
			if err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					log.Printf("Auth attempt for user ID %s, but user not found.", userID)
					return nil, errors.New("user not found")
				}
				// Otro error de base de datos
				log.Printf("Database error fetching user %s: %v", userID, err)
				return nil, err
			}

			// 5. Devolver la clave específica de este usuario para la verificación
			return []byte(psiUser.Key), nil
		})

		// 6. Manejar errores del parseo final (firma inválida, token expirado, etc.)
		if err != nil {
			log.Printf("JWT validation error: %v", err)
			return jwtError(c, fiber.StatusUnauthorized, "Invalid or expired JWT")
		}

		if !token.Valid {
			return jwtError(c, fiber.StatusUnauthorized, "Invalid or expired JWT")
		}

		// 7. Si todo está bien, guardar el token validado en el contexto y continuar.
		c.Locals("user", token)
		return c.Next()
	}
}
