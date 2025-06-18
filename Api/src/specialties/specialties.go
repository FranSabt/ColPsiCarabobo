package specialties

import (
	"fmt"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type specialtiesRequest struct {
	Name     string
	AdmindId string
}

func GetPsiUsers(c *fiber.Ctx, db *gorm.DB) error {
	var request specialtiesRequest

	if err := c.BodyParser(request); err != nil {
		return c.JSON(fiber.Map{
			"success": false,
			"error":   "Internal server error",
		})
	}

	id, err := uuid.Parse(request.AdmindId)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "invalid id format",
			"details": err.Error(),
		})
	}

	fmt.Println(id)

	return c.JSON(fiber.Map{
		"success": true,
		"data":    nil,
	})
}
