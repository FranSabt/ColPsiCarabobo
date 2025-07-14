package psiuser_presenter

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"strings"

	"github.com/FranSabt/ColPsiCarabobo/db"
	image_controller "github.com/FranSabt/ColPsiCarabobo/src/images/controller"
	db_images "github.com/FranSabt/ColPsiCarabobo/src/images/db"
	psi_user_db "github.com/FranSabt/ColPsiCarabobo/src/psi-user/db"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type ImageRequest struct {
	UserID   string `json:"psi_user_id"`
	Filename string `json:"file_name"`
	MimeType string `json:"mime_type"`
	Data     []byte `json:"data"` // base64 si viene como JSON
}

func CreatePsiUserImage(c *fiber.Ctx, db db.StructDb) error {
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

	// Parsear UUID del usuario
	uid, err := uuid.Parse(userID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Invalid UUID format",
			"success": false,
		})
	}

	// Validar si existe el usuario
	fmt.Println("Buscando Usuario")
	user, err := psi_user_db.GetPsiUserById(db.DB, uid)
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
	// El usuario esta solvente
	if !user.Solvent {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "User is inactive",
			"success": false,
		})
	}

	// ---- CHECAR IMAGENES ---- //
	count_image, err := db_images.CheckProfilePicLimit(user.ID, db.Image)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Can't connect to DB",
			"success": false,
		})
	}

	can_put, message := canCreatePic(count_image)
	if !can_put {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Can't connect to DB",
			"success": false,
			"message": message,
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
	if err := db_images.SaveUserImage(*model, db.Image); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to save image data",
			"success": false,
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"message": "Image saved successfully",
		"format":  format,
		"size":    compressed.Len(),
	})
}

//////////////////////////////////////////////////////
//////////////////////////////////////////////////////
//////////////////////////////////////////////////////

func UpdatePsiUserImage(c *fiber.Ctx, db db.StructDb) error {
	// Obtener campos de texto del form-data
	userID := c.FormValue("psi_user_id")
	filename := c.FormValue("file_name")
	mimeType := c.FormValue("mime_type")
	image_id := c.FormValue("image_id")

	// Validar campos de texto
	if userID == "" || filename == "" || mimeType == "" || image_id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Missing required form fields",
			"success": false,
		})
	}

	// Obtener archivo
	fileHeader, err := c.FormFile("image")
	if err != nil || fileHeader == nil {
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

	// Parsear UUID del usuario
	uid, err := uuid.Parse(userID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Invalid UUID format",
			"success": false,
		})
	}

	// Validar si existe el usuario
	fmt.Println("Buscando Usuario")
	user, err := psi_user_db.GetPsiUserById(db.DB, uid)
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
	// El usuario esta solvente
	if !user.Solvent {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "User is inactive",
			"success": false,
		})
	}

	// ---- La Imagen Existe? ---- //
	uid_image, err := uuid.Parse(image_id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Invalid UUID format for image",
			"success": false,
		})
	}

	// Usar la nueva función GetImageById
	existingImage, err := db_images.GetImageById(uid_image, db.Image)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Can't connect to DB",
			"success": false,
		})
	}
	if existingImage == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error":   "Image does not exist",
			"success": false,
		})
	}

	// Verificar que la imagen pertenece al usuario (seguridad)
	if existingImage.UserID != uid {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error":   "You don't have permission to update this image",
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
	fmt.Println("Comprimiendo imagen")
	compressed, _, err := image_controller.CompressImages(&data, img)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Image compression failed",
			"success": false,
		})
	}

	// Actualizar imagen en la DB usando la nueva función
	fmt.Println("Actualizando en BD")
	success, err := db_images.UpdateImageById(uid_image, compressed.Bytes(), format, db.Image, filename)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to update image data",
			"success": false,
		})
	}
	if !success {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to update image",
			"success": false,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Image updated successfully",
		"format":  format,
		"size":    compressed.Len(),
	})
}

//////////////////////////////////////////////////////
//////////////////////////////////////////////////////
//////////////////////////////////////////////////////

func DeletePsiUserImage(c *fiber.Ctx, db db.StructDb) error {
	// Obtener parámetros (pueden venir del body JSON o como parámetros de ruta)
	userID := c.FormValue("psi_user_id")
	image_id := c.FormValue("image_id")

	// Si no vienen por form-data, intentar obtenerlos del JSON body
	if userID == "" || image_id == "" {
		var requestBody struct {
			UserID  string `json:"psi_user_id"`
			ImageID string `json:"image_id"`
		}

		if err := c.BodyParser(&requestBody); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error":   "Invalid request format",
				"success": false,
			})
		}

		if requestBody.UserID != "" {
			userID = requestBody.UserID
		}
		if requestBody.ImageID != "" {
			image_id = requestBody.ImageID
		}
	}

	// Validar campos requeridos
	if userID == "" || image_id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Missing required fields: psi_user_id and image_id",
			"success": false,
		})
	}

	// Parsear UUID del usuario
	uid, err := uuid.Parse(userID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Invalid UUID format for user",
			"success": false,
		})
	}

	// Parsear UUID de la imagen
	uid_image, err := uuid.Parse(image_id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Invalid UUID format for image",
			"success": false,
		})
	}

	// Validar si existe el usuario
	fmt.Println("Buscando Usuario")
	user, err := psi_user_db.GetPsiUserById(db.DB, uid)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Can't connect to DB",
			"success": false,
		})
	}
	if user == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error":   "User doesn't exist",
			"success": false,
		})
	}

	// El usuario esta solvente
	if !user.Solvent {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "User is inactive",
			"success": false,
		})
	}

	// ---- Verificar que la imagen existe ---- //
	fmt.Println("Buscando Imagen")
	existingImage, err := db_images.GetImageById(uid_image, db.Image)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Can't connect to DB",
			"success": false,
		})
	}
	if existingImage == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error":   "Image does not exist",
			"success": false,
		})
	}

	// Verificar que la imagen pertenece al usuario (seguridad)
	if existingImage.UserID != uid {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error":   "You don't have permission to delete this image",
			"success": false,
		})
	}

	// Eliminar imagen de la DB
	fmt.Println("Eliminando imagen de BD")
	success, err := db_images.DeleteImageById(uid_image, db.Image)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to delete image",
			"success": false,
		})
	}
	if !success {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to delete image",
			"success": false,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Image deleted successfully",
		"deleted_image": fiber.Map{
			"id":      existingImage.ID,
			"name":    existingImage.Name,
			"format":  existingImage.Format,
			"user_id": existingImage.UserID,
		},
	})
}

// ////////////////////////////////////////////////////
// ////////////////////////////////////////////////////
// ////////////////////////////////////////////////////
// TODO: Modificar para que busque por id
func GetMyProfilePic(c *fiber.Ctx, db db.StructDb) error {

	userID := c.Query("id")
	if userID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Empty param id",
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
	user, err := psi_user_db.GetPsiUserById(db.DB, uid)
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

	image, err := db_images.GetFirstImageByAssociatedID(uid, db.Image)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Can't connect to DB",
			"success": false,
		})
	}
	if image == nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "User doesn't exist",
			"success": false,
		})
	}

	c.Set("Content-Type", "image/"+image.Format) // Ej: "image/png", "image/jpeg"
	c.Set("Content-Disposition", "inline; filename="+image.Name)

	// Enviar los datos binarios de la imagen
	return c.Status(fiber.StatusOK).Send(image.ImageData)
}

//////////////////////////////////////////////////////
//////////////////////////////////////////////////////
//////////////////////////////////////////////////////

func GetMyProfilePicList(c *fiber.Ctx, db db.StructDb) error {
	userID := c.Query("id")
	if userID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Empty param id",
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
	user, err := psi_user_db.GetPsiUserById(db.DB, uid)
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

	images_ids, err := db_images.GetAllUserProfilePicsID(uid, db.Image)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Can't connect to DB",
			"success": false,
		})
	}
	if images_ids == nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "User doesn't exist",
			"success": false,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success":    true,
		"message":    "Image deleted successfully",
		"images:ids": images_ids,
	})
}

//////////////////////////////////////////////////////
//////////////////////////////////////////////////////
//////////////////////////////////////////////////////

// ------       Auxiliary Functions      ------     //

//////////////////////////////////////////////////////
//////////////////////////////////////////////////////
//////////////////////////////////////////////////////

func canCreatePic(currentImageCount int64) (allowed bool, mesage string) {
	switch {
	case currentImageCount < 0:
		return false, "se colo un error -1"
	case currentImageCount < 3:
		return true, ""
	case currentImageCount >= 3:
		return false, "limite de imgenes alcanzado"
	default:
		return false, "error desconocido"
	}
}
