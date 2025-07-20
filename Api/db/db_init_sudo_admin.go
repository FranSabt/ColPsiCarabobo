package db

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/FranSabt/ColPsiCarabobo/src/models"
	"golang.org/x/crypto/argon2"
	"gorm.io/gorm"
)

// HasSudoAdmin, verifica si ya existe un superusuario
// Versión optimizada y corregida.
func hasSudoAdmin(db *gorm.DB) (bool, error) {
	var count int64
	// Usamos Model y Count para mayor eficiencia.
	// Solo contamos cuántos usuarios cumplen la condición.
	err := db.Model(&models.UserAdmin{}).Where("sudo = ?", true).Count(&count).Error
	if err != nil {
		// No es necesario chequear ErrRecordNotFound con Count,
		// ya que si no hay resultados, count será 0 y err será nil.
		return false, nil
	}

	// Si count > 0, significa que existe al menos un super admin.
	return count > 0, nil
}

// CreateSudoAdmin crea un super usuario de forma segura.
func createSudoAdmin(db *gorm.DB) error {
	// 1. Verificamos si ya existe un super admin para evitar duplicados.
	exists, err := hasSudoAdmin(db)
	if err != nil {
		return fmt.Errorf("error al verificar la existencia de super admin: %w", err)
	}
	if exists {
		// Es mejor devolver un error conocido que fallar silenciosamente.
		fmt.Println("ya existe un super administrador en la base de datos")
		return nil
	}

	// 2. Obtenemos y validamos las variables de entorno.
	adminUsername := os.Getenv("ADMIN_USERNAME")
	password := os.Getenv("ADMIN_PASSWORD")
	email := os.Getenv("ADMIN_EMAIL") // Corregí el nombre de la variable de "ADMIN_Email"

	if adminUsername == "" || password == "" || email == "" {
		return errors.New("las variables de entorno ADMIN_USERNAME, ADMIN_PASSWORD y ADMIN_EMAIL son obligatorias")
	}

	// 3. Hasheamos la contraseña antes de guardarla.
	hashedPassword, err := hashPassword(password)
	if err != nil {
		return fmt.Errorf("could not hash password: %v", err)
	}

	// Creamos la instancia del admin.
	adminModel := models.UserAdmin{
		Username:       adminUsername,
		Password:       hashedPassword, // Guardamos el hash, no la contraseña original
		Email:          email,
		Sudo:           true, // Importante: ¡No olvidemos marcarlo como superusuario!
		CanCreateAdmin: true,
		CanUpdateAdmin: true,
		CanDeleteAdmin: true,
		// publicaciones
		CanPublish:       true,
		CanUpdatePublish: true,
		CanDeletePublish: true,
		// notificaciones
		CanSendNotifications:   true,
		CanManageNotifications: true,
		CanReadNotifications:   true,
		// tags
		CanCreateTags: true,
		CanDeleteTags: true,
		CanEditTags:   true,
	}

	// La transacción se mantiene igual, es una buena práctica.
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 4. Usamos Create con un puntero.
	if err := tx.Create(&adminModel).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func hashPassword(password string) (string, error) {
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
