package psiuser_presenter

import (
	"errors"
	"log"
	"strings"

	psi_user_db "github.com/FranSabt/ColPsiCarabobo/src/psi-user/db"
	"github.com/FranSabt/ColPsiCarabobo/src/utils"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type PsiUserLoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func PsiUserLogin(c *fiber.Ctx, db *gorm.DB) error {
	var req PsiUserLoginRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Invalid request format",
		})
	}

	// Validar campos requeridos
	req.Username = strings.TrimSpace(req.Username)
	req.Password = strings.TrimSpace(req.Password)
	if req.Username == "" || req.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Username and password are required",
		})
	}

	// Determinar el tipo de búsqueda (email, SUB, o username)
	searchField := getSearchField(req.Username)
	log.Println("Search field:", searchField)

	// Buscar usuario
	psiUser, err := psi_user_db.GetPsiUserByUsernameOrEmal(db, strings.ToLower(req.Username), searchField)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Log de seguridad
			log.Printf("Login attempt for non-existent user: %s", req.Username)
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"success": false,
				"error":   "Invalid credentials",
			})
		}

		log.Printf("Database error: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   "Internal server error",
		})
	}

	// Verificar contraseña con manejo de errores
	if !utils.CheckPasswordHash(req.Password, psiUser.Password) {
		// Log de seguridad
		log.Printf("Failed login attempt for user: %s", psiUser.Username)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"error":   "Invalid credentials2",
		})
	}

	// // Generar token JWT
	// token, err := generateJWTToken(psiUser)
	// if err != nil {
	//     log.Printf("JWT generation error: %v", err)
	//     return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
	//         "success": false,
	//         "error":   "Internal server error",
	//     })
	// }

	return c.JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"user_id":  psiUser.ID,
			"username": psiUser.Username,
			// "token":    token,
		},
	})
}

func getSearchField(username string) string {
	if strings.Contains(username, "@") {
		return "email = ?"
	}
	return "username = ?"
}
