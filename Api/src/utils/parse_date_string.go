package utils

import (
	"fmt"
	"time"
)

func ParseDateString(dateString string) (time.Time, error) {
	// Lista de formatos de fecha comunes que intentaremos analizar
	formats := []string{
		"2006-01-02",   // Formato ISO (YYYY-MM-DD)
		"02/01/2006",   // Formato DD/MM/YYYY
		"01/02/2006",   // Formato MM/DD/YYYY
		"Jan 02, 2006", // Ej: "Dec 25, 2023"
		"02-Jan-2006",  // Ej: "25-Dec-2023"
		time.RFC3339,   // Formato ISO con zona horaria
	}

	var parsedTime time.Time
	var err error

	// Intentar parsear con cada formato hasta que uno funcione
	for _, format := range formats {
		parsedTime, err = time.Parse(format, dateString)
		if err == nil {
			return parsedTime, nil
		}
	}

	return time.Time{}, fmt.Errorf("no se pudo parsear la fecha: %v, formatos intentados: %v", dateString, formats)
}
