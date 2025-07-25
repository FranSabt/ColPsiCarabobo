package image_controller

import (
	"bytes"
	"image"
	"image/jpeg"
	"strings"
	"time"

	"github.com/FranSabt/ColPsiCarabobo/src/models"
	"github.com/google/uuid"
)

type ImageRequest struct {
	ID       uuid.UUID `json:"id"`
	UserID   uuid.UUID `json:"psi_user_id"`
	Filename string    `json:"file_name"`
	MimeType string    `json:"mime_type"`
	Data     []byte    `gorm:"type:bytea"`
}

const maxImageSize = 5 * 1024 * 1024 // 5MB

func CreateImageModel(user models.PsiUserModel, name, format string) models.ProfilePicModel {
	id := uuid.New()

	model := models.ProfilePicModel{
		ID:         id,
		UserID:     user.ID,
		Name:       sanitizeFileName(name),
		Format:     format,
		CreatedAt:  time.Now(),
		CreateById: &user.ID,
		CreateBy:   user.Username,
		UpdatedAt:  time.Now(),
		UpdateById: &user.ID,
		UpdateBy:   user.Username,
	}

	return model
}

func sanitizeFileName(name string) string {
	name = strings.ReplaceAll(name, "..", "")
	name = strings.ReplaceAll(name, "=", "")
	return strings.ReplaceAll(name, "/", "")
}

func CompressImages(data *[]byte, img image.Image) (*bytes.Buffer, image.Image, error) {
	var compressed bytes.Buffer

	if len(*data) > maxImageSize {
		quality := 90
		for {
			compressed.Reset()
			if err := jpeg.Encode(&compressed, img, &jpeg.Options{Quality: quality}); err != nil {
				return nil, nil, err
			}
			if compressed.Len() <= maxImageSize || quality <= 30 {
				break
			}
			quality -= 5
		}
	} else {
		compressed = *bytes.NewBuffer(*data)
	}

	return &compressed, img, nil
}
