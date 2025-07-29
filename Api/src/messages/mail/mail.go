package messages_email

import (
	"fmt"

	"github.com/FranSabt/ColPsiCarabobo/config"
	"gopkg.in/gomail.v2"
)

// SendEmail envía un correo electrónico a la dirección de correo especificada.
// Se puede enviar cualquier formato de correo y cualquier numero de receptores.
// Las validaciones quedan de parte de cada grupo logico a implementar.
func SendEmail(mail, body, subject, destiny string) error {
	// Configuración del servidor SMTP
	smtpHost := config.EnvConfig("SMTP_HOST")
	smtpPort := 587 // Cambiado a un entero

	// Credenciales de autenticación
	sender := config.EnvConfig("SMTP_HOST_USER")
	password := config.EnvConfig("SMTP_HOST_PASSWORD")
	fmt.Println("SMTP_HOST:", smtpHost)
	fmt.Println("SMTP_USER:", sender)

	if smtpHost == "" || sender == "" || password == "" || destiny == "" {
		return fmt.Errorf("faltan credenciales de autenticación")
	}

	// Crear un nuevo mensaje
	m := gomail.NewMessage()
	m.SetHeader("From", sender) // El remitente debe ser el correo del usuario (no el host)
	// m.SetHeader("To", mail)     // Correo del destinatario
	m.SetHeader("To", destiny) // Correo del destinatario
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body) // Cuerpo del correo

	// Configuración del dialer con el host, puerto y autenticación
	d := gomail.NewDialer(smtpHost, smtpPort, sender, password)
	d.SSL = false // Asegurar que se usan SSL

	// Enviar el correo y manejar errores
	if err := d.DialAndSend(m); err != nil {
		fmt.Printf("Error al enviar el correo: %v\n", err)
	} else {
		fmt.Println("Correo enviado exitosamente")
	}

	return nil
}
