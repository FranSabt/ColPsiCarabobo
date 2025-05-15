package utils

import (
	"crypto/rand"
	"encoding/base64"
	"log"

	"golang.org/x/crypto/argon2"
)

func HashPassword(password string) (string, error) {
	const saltLength = 16
	const keyLength = 32

	// Generar salt
	salt := make([]byte, saltLength)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}

	// Generar hash usando argon2 y el salt generado
	hash := argon2.IDKey([]byte(password), salt, 1, 64*1024, 4, keyLength)

	// Combinar salt y hash, y codificarlos en base64
	hashSalt := append(salt, hash...)
	encoded := base64.RawStdEncoding.EncodeToString(hashSalt)

	log.Printf("Pass: %v\nHash generado: %v", password, encoded) // Log para depuración
	return encoded, nil
}
