package psi_user_controller

import (
	"bytes"
	"errors"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"mime/multipart"
	"strings"
)

func ImageProcesser(file_header *multipart.FileHeader) ([]byte, error) {
	// Abrir archivo
	file, err := file_header.Open()
	if err != nil {
		return nil, err
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	return data, err
}

// ImageDecoder decodes image data based on the provided MIME type.
func ImageDecoder(data []byte, mimeType string) (image.Image, string, error) {
	// Declare all variables that will be assigned inside the switch.
	var img image.Image
	var err error // NOTE: Declare err here.
	var format string

	reader := bytes.NewReader(data)

	switch strings.ToLower(mimeType) {
	case "image/jpeg", "image/jpg":
		format = "jpg"
		img, err = jpeg.Decode(reader)
	case "image/png":
		format = "png"
		img, err = png.Decode(reader)
	default:
		return nil, "", errors.New("unsupported image format (only JPEG or PNG)")
	}

	// function-scoped 'err' variable.
	if err != nil {
		return nil, "", err
	}

	// The function-scoped 'img' variable now holds the decoded image.
	return img, format, nil
}
