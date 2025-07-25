package post_presenter

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	admin_controller "github.com/FranSabt/ColPsiCarabobo/src/admin/controller"
	"github.com/FranSabt/ColPsiCarabobo/src/models"
	post_db "github.com/FranSabt/ColPsiCarabobo/src/posts/db"
	text_db "github.com/FranSabt/ColPsiCarabobo/src/text/db"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func GetPosts(c *fiber.Ctx, db *gorm.DB) error {

	post_type := c.Query("post_type", "")
	page := c.QueryInt("page", 1)
	page_size := c.QueryInt("page_size", 20)

	posts, err := post_db.GetActivePostsPaginated(post_type, page, page_size, db)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"succes":  false,
			"error":   "Error al obtener los psicólogos",
			"details": err.Error(),
		})
	}

	return c.JSON(posts)
}

func GetPostText(c *fiber.Ctx, db *gorm.DB, database_text *gorm.DB) error {

	id := c.QueryInt("id", 0)
	if id < 1 {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"succes":  false,
			"error":   "invalid id int",
			"details": "id can be 0 or less",
		})
	}

	uint_id := uint(id)

	posts, err := post_db.GetPostById(uint_id, db)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"succes":  false,
			"error":   fmt.Sprintf("error obtainint the post with id: %v", id),
			"details": err.Error(),
		})
	}

	post_text, err := text_db.GetTextByIDDb(database_text, posts.TextID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"succes":  false,
			"error":   fmt.Sprintf("error obtainint the post TEXT with text id: %v", posts.TextID),
			"details": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"succes":  true,
		"error":   nil,
		"details": nil,
		"data":    post_text,
	})
}

type CreatePostRequest struct {
	Name             string `json:"name" validate:"required,min=3,max=50"`
	Type             string `json:"type" validate:"required,oneof=public psi"` // Asumo que PostType es "public" o "psi"
	ShortDescription string `json:"short_description" validate:"max=250"`
	Text             string `json:"text" validate:"required"`

	// Campos opcionales que se pueden establecer en la creación
	IsActive *bool `json:"is_active"` // Usamos puntero para poder distinguir entre 'false' y 'no proporcionado'

	// Este campo no viene en el JSON, se añade en el handler desde el token
	CreatedBy string `json:"create_by"`
}

func CreatePostAdmin(c *fiber.Ctx, db *gorm.DB, db_text *gorm.DB) error {
	var request CreatePostRequest

	if err := c.BodyParser(&request); err != nil {
		return c.JSON(fiber.Map{
			"success": false,
			"error":   "Internal server error",
		})
	}

	admin_id := request.CreatedBy

	admin_uuid, err := uuid.Parse(admin_id)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "invalid id format",
			"details": err.Error(),
		})
	}

	admin, err := admin_controller.GetAdminById(admin_uuid, db)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "invalid id format",
			"details": err.Error(),
		})
	}

	if admin.Username == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "invalid id format",
			"details": "admin does not exist",
		})
	}

	err = createPost(request, admin, db, db_text)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "cant create post",
			"details": err.Error(),
		})
	}

	return c.JSON("ok")
}

func createPost(request CreatePostRequest, admin models.UserAdmin, db *gorm.DB, db_text *gorm.DB) error {
	// TODO: Funcion para limpiar el texto
	text_model := models.TextModel{
		Text:   request.Text,
		Active: true,
	}

	text_model_saved, err := text_db.CreateTextDb(db_text, text_model)
	if err != nil {
		return err
	}

	// --- COMIENZA LA CORRECCIÓN ---

	// 1. Convertir y validar el string del request a models.PostType
	postType, ok := models.StringToPostType[request.Type]
	if !ok {
		// Si 'ok' es false, el tipo de post no es válido.
		return fmt.Errorf("invalid post type: %q. valid types are 'public' or 'psi'", request.Type)
	}

	// --- FIN DE LA CORRECCIÓN ---

	post_model := models.Post{
		Name:             request.Name,
		Type:             postType, // <-- Ahora usas la variable convertida
		ShortDescription: request.ShortDescription,
		TextID:           text_model_saved.ID,
		IsActive:         true,
		CreatedAt:        time.Now(),
		CreateBy:         admin.Username,
		CreateById:       admin.ID,
		UpdateBy:         admin.Username,
		UpdateById:       admin.ID,
		UpdatedAt:        time.Now(),
	}

	err = post_db.CreatePost(post_model, db)
	if err != nil {
		// Aquí podrías querer borrar el texto que creaste si la creación del post falla (rollback manual)
		// text_db.DeleteTextDb(db_text, text_model_saved.ID)
		return err
	}

	return nil
}

// Este es el DTO correcto para una petición de actualización.
// Lo que se actualiza viene en el body. El ID del post viene en la URL.
type UpdatePostRequest struct {
	Id               int64   `json:"post_id"`
	Name             *string `json:"name,omitempty"`
	Type             *string `json:"type,omitempty"`
	ShortDescription *string `json:"short_description,omitempty"`
	Text             *string `json:"text,omitempty"`
	IsActive         *bool   `json:"is_active,omitempty"`
	UpdateBy         string  `json:"update_by"`
}

func UpdatePost(c *fiber.Ctx, db *gorm.DB, db_text *gorm.DB) error {
	// --- 1. Obtener IDs de fuentes seguras ---

	// --- 2. Parsear el Body y Validar ---

	var request UpdatePostRequest
	if err := c.BodyParser(&request); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Invalid request body",
		})
	}

	// Obtener el ID del POST desde la URL (ej: /posts/123)
	postID := request.Id
	if postID <= 0 {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Invalid post ID provided in URL",
		})
	}

	admin_id := request.UpdateBy

	admin_uuid, err := uuid.Parse(admin_id)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "invalid id format",
			"details": err.Error(),
		})
	}

	admin, err := admin_controller.GetAdminById(admin_uuid, db)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "invalid id format",
			"details": err.Error(),
		})
	}

	if admin.Username == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "admin does not exist",
			"details": "admin does not exist",
		})
	}

	uint_id := uint(postID)
	post_to_update, err := post_db.GetPostById(uint_id, db)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return c.Status(http.StatusNotFound).JSON(fiber.Map{
				"success": false,
				"error":   "Post not found",
			})
		}
		// Para otros errores de DB
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Database error fetching post"})
	}

	// --- 4. Lógica de Actualización ---

	// a) Actualizar el texto en la base de datos de texto, si se proporcionó
	if request.Text != nil {
		textUpdateData := map[string]interface{}{"Text": *request.Text}
		if err := text_db.UpdateTextDb(db_text, post_to_update.TextID, textUpdateData); err != nil {
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update post text"})
		}
	}

	// b) Preparar los datos para actualizar el Post en la base de datos principal
	postUpdateData := models.Post{}
	update := false

	if request.Name != nil {
		postUpdateData.Name = *request.Name
		update = true
	}
	if request.ShortDescription != nil {
		postUpdateData.ShortDescription = *request.ShortDescription
		update = true

	}
	if request.IsActive != nil {
		postUpdateData.IsActive = *request.IsActive
		update = true

	}
	if request.Type != nil {
		// Convertir y validar el string a models.PostType
		postType, ok := models.StringToPostType[*request.Type]
		if !ok {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": fmt.Sprintf("Invalid post type: %s", *request.Type)})
		}
		postUpdateData.Type = postType
		update = true

	}

	// Añadir siempre el ID del usuario que actualiza
	postUpdateData.UpdateBy = admin.Username
	postUpdateData.UpdateById = admin.ID
	postUpdateData.UpdatedAt = time.Now()

	// c) Ejecutar la actualización del Post si hay algo que actualizar
	if update { // Mayor que 1 porque siempre incluimos UpdatedBy
		if _, err := post_db.UpdatePost(db, uint_id, postUpdateData); err != nil {
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update post metadata"})
		}
	}

	// --- 5. Devolver Respuesta Exitosa ---

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Post updated successfully",
	})
}
