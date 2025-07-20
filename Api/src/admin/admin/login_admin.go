package admin

import (
	"errors"
	"log"
	"strings"
	"time"

	db_admin "github.com/FranSabt/ColPsiCarabobo/src/admin/db"
	"github.com/FranSabt/ColPsiCarabobo/src/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

type AdminLoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func AdminLogin(c *fiber.Ctx, db *gorm.DB) error {
	var req AdminLoginRequest

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
	admin, err := db_admin.GetAdminByUsernameOrEmal(db, strings.ToLower(req.Username), searchField)
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
	if !utils.CheckPasswordHash(req.Password, admin.Password) {
		// Log de seguridad
		log.Printf("Failed login attempt for user: %s", admin.Username)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"error":   "Invalid credentials",
		})
	}

	admin.Key = utils.GenerateSecureRandomString(512)
	err = db_admin.SaveUpdatedAdminOnly(db, admin)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
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

	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["username"] = admin.Username
	claims["user_id"] = admin.ID
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

	t, err := token.SignedString([]byte(admin.Key))
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	// return c.JSON(fiber.Map{"status": "success", "message": "Success login", "data": t})

	return c.JSON(fiber.Map{
		"success": true,
		"data":    t,
	})
}

func getSearchField(username string) string {
	if strings.Contains(username, "@") {
		return "email = ?"
	}
	return "username = ?"
}
