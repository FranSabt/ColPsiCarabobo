package utils

import (
	"crypto/subtle"
	"encoding/base64"
	"log"

	"golang.org/x/crypto/argon2"
)

func CheckPasswordHash(password, encodedHash string) bool {
	const saltLength = 16
	const keyLength = 32

	// Verificar si el hash está vacío
	if encodedHash == "" {
		log.Println("Empty encoded hash")
		return false
	}

	// Usar DecodeString estándar que maneja padding
	hashSalt, err := base64.StdEncoding.DecodeString(encodedHash)
	if err != nil {
		log.Printf("Error decoding hash: %v | Hash: %s", err, encodedHash)
		return false
	}

	// Verificar longitud del hash decodificado
	if len(hashSalt) != saltLength+keyLength {
		log.Printf("Invalid hash length: %d", len(hashSalt))
		return false
	}

	salt := hashSalt[:saltLength]
	storedHash := hashSalt[saltLength:]

	newHash := argon2.IDKey([]byte(password), salt, 1, 64*1024, 4, keyLength)

	return subtle.ConstantTimeCompare(storedHash, newHash) == 1
}
