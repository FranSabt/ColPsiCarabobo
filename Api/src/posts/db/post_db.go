package post_db

import (
	"errors"

	"github.com/FranSabt/ColPsiCarabobo/src/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func CreatePost(post_model models.Post, db *gorm.DB) error {
	// Intentar crear el registro en la base de datos
	result := db.Create(&post_model)
	if result.Error != nil {
		// Si hay un error, lo retornamos
		return result.Error
	}

	// Si todo está bien, retornamos nil (sin error)
	return nil
}

func GetPostById(id uuid.UUID, db *gorm.DB) (*models.Post, error) {
	psiUser := &models.Post{}
	err := db.Where("id = ?", id).First(psiUser).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("psi_user record not found")
		}
		return nil, err
	}

	return psiUser, nil
}

// PaginationResult es un struct genérico que puede ser útil para las respuestas JSON
type PaginationResult struct {
	Data     interface{} `json:"data"`
	Total    int64       `json:"total"`
	Page     int         `json:"page"`
	PageSize int         `json:"page_size"`
}

// GetActivePostsPaginated recupera una lista paginada de posts activos.
// Devuelve los posts, el conteo total de registros activos, y un error.
func GetActivePostsPaginated(post_type string, page, pageSize int, db *gorm.DB) (*PaginationResult, error) {
	if post_type == "" {
		post_type = "public"
	}

	// --- 1. Validar y sanear los parámetros de paginación ---
	if page <= 0 {
		page = 1
	}
	switch {
	case pageSize > 100:
		pageSize = 20
	case pageSize <= 0:
		pageSize = 20
	}
	offset := (page - 1) * pageSize

	// --- 2. Crear una consulta base con el filtro de estado ---
	// De esta forma, nos aseguramos de que ambas consultas (Count y Find)
	// operen sobre el mismo conjunto de datos.
	baseQuery := db.Model(&models.Post{}).Where("is_active = ?", true).Where("type = ?", post_type)

	// --- 3. Obtener el conteo total usando la consulta base ---
	var totalCount int64
	if err := baseQuery.Count(&totalCount).Error; err != nil {
		return nil, err
	}

	// --- 4. Obtener los registros de la página actual, también usando la consulta base ---
	var posts []models.Post

	// Partimos de la 'baseQuery' y le añadimos el resto de condiciones
	err := baseQuery.
		Preload("Creator").
		Preload("Updater").
		Order("created_at DESC").
		Limit(pageSize).
		Offset(offset).
		Find(&posts).Error

	if err != nil {
		return nil, err
	}

	result := PaginationResult{
		Data:     posts,
		Total:    totalCount,
		Page:     page,
		PageSize: pageSize,
	}

	return &result, nil
}

// UpdatePost busca un post por su ID y actualiza sus campos.
// Devuelve el post actualizado o un error.
func UpdatePost(db *gorm.DB, id uuid.UUID, updateData models.Post) (*models.Post, error) {
	var post models.Post

	if err := db.Model(&post).Updates(updateData).Error; err != nil {
		return nil, err
	}

	// (Opcional pero recomendado) Recargar el post con los datos del 'Updater'
	// para devolver el objeto completo en la respuesta.
	if err := db.Preload("Updater").First(&post, post.ID).Error; err != nil {
		return nil, err
	}

	return &post, nil
}
