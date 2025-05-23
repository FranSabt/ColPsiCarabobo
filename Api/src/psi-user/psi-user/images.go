package psiuser_presenter

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"strings"

	image_controller "github.com/FranSabt/ColPsiCarabobo/src/images/controller"
	db_images "github.com/FranSabt/ColPsiCarabobo/src/images/db"
	psi_user_db "github.com/FranSabt/ColPsiCarabobo/src/psi-user/db"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ImageRequest struct {
	UserID   string `json:"psi_user_id"`
	Filename string `json:"file_name"`
	MimeType string `json:"mime_type"`
	Data     []byte `json:"data"` // base64 si viene como JSON
}

func CreatePsiUserImage(c *fiber.Ctx, db *gorm.DB) error {
	// Obtener campos de texto del form-data
	userID := c.FormValue("psi_user_id")
	filename := c.FormValue("file_name")
	mimeType := c.FormValue("mime_type")

	// Validar campos de texto
	if userID == "" || filename == "" || mimeType == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Missing required form fields",
			"success": false,
		})
	}

	// Obtener archivo
	fileHeader, err := c.FormFile("image")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Image file is required",
			"success": false,
		})
	}

	// Abrir archivo
	file, err := fileHeader.Open()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to open image file",
			"success": false,
		})
	}
	defer file.Close()

	// Leer archivo en []byte
	data, err := io.ReadAll(file)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to read image data",
			"success": false,
		})
	}

	// Parsear UUID
	uid, err := uuid.Parse(userID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Invalid UUID format",
			"success": false,
		})
	}

	// Validar si existe el usuario
	fmt.Println("Buscando Usuario")
	user, err := psi_user_db.GetPsiUserById(db, uid)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Can't connect to DB",
			"success": false,
		})
	}
	if user == nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "User doesn't exist",
			"success": false,
		})
	}

	// Truncar el nombre si es muy largo
	if len(filename) > 50 {
		filename = filename[:50]
	}

	// Decodificar imagen según MIME
	var img image.Image
	reader := bytes.NewReader(data)
	var format string

	fmt.Println("Formato")
	switch strings.ToLower(mimeType) {
	case "image/jpeg", "image/jpg":
		format = "jpg"
		img, err = jpeg.Decode(reader)
	case "image/png":
		format = "png"
		img, err = png.Decode(reader)
	default:
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Unsupported image format (only JPEG or PNG)",
			"success": false,
		})
	}
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Failed to decode image",
			"success": false,
		})
	}

	// Comprimir imagen
	fmt.Println("Imagen")
	compressed, _, err := image_controller.CompressImages(&data, img)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Image compression failed",
			"success": false,
		})
	}

	// Crear modelo para guardar
	fmt.Println("Modelo")
	model := image_controller.CreateImageModel(uid, filename, format)
	if model == nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to format data to save",
			"success": false,
		})
	}

	model.ImageData = compressed.Bytes()

	// Guardar en la DB
	fmt.Println("BD")
	if err := db_images.SaveUserImage(*model); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to save image data",
			"success": false,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Image saved successfully",
		"format":  format,
		"size":    compressed.Len(),
	})
}
