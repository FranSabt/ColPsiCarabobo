package psi_user_controller

import (
	"math/rand"
	"time"
)

func RandomPass() string {
	// Caracteres permitidos
	caracteres := "ABCDEFGHIJKLMNOPQRSTUVlongitudWXYZabcdefghijklmnopqrstuvwxyz0123456789"

	// Crear un generador de números aleatorios local
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	// Crear un slice de bytes para almacenar la cadena aleatoria
	cadena := make([]byte, 8)

	// Llenar el slice con caracteres aleatorios
	for i := range cadena {
		cadena[i] = caracteres[r.Intn(len(caracteres))]
	}

	// Convertir el slice de bytes a string y devolverlo
	return string(cadena)
}
