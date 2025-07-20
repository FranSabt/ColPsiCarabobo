package db

import (
	"fmt"
	"log"
	"os"

	"github.com/FranSabt/ColPsiCarabobo/src/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type StructDb struct {
	Image *gorm.DB
	DB    *gorm.DB
}

type Dbinstance struct {
	Db *gorm.DB
}

var DB Dbinstance

func Connect() (*gorm.DB, error) {
	// Obtener las variables de entorno
	host := os.Getenv("DB_HOST")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")
	port := os.Getenv("DB_PORT")
	timezone := os.Getenv("DB_TIMEZONE")
	// automigrate := os.Getenv("AUTOMIGRATE")
	// development := os.Getenv("DEVELOPMENT")

	// Crear el DSN (Data Source Name)
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=%s",
		host, user, password, dbname, port, timezone)

	// Configuración de Gorm (puedes ajustar el logger según tus necesidades)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent), // Puedes cambiar el nivel del logger si es necesario
	})

	if err != nil {
		return nil, fmt.Errorf("could not connect to the database: %w", err)
	}

	// Opcionalmente, configurar la conexión a la base de datos (por ejemplo, conexión máxima, tiempo de espera, etc.)
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("could not configure the database connection: %w", err)
	}

	log.Println("connected")
	db.Logger = logger.Default.LogMode(logger.Info)
	log.Println("running migrations")

	// if development == "true" && automigrate == "true" {
	// 	AutoMigrateDB(db)
	// }

	err = createSudoAdmin(db)
	if err != nil {
		return nil, fmt.Errorf("error checking the sudo admin: %w", err)
	}

	// Configurar los parámetros de la conexión, como máximo número de conexiones abiertas, etc.
	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetMaxIdleConns(25)
	sqlDB.SetConnMaxLifetime(5 * 60) // Ejemplo de 5 minutos

	return db, nil
}

func AutoMigrateDB(db *gorm.DB) error {
	// Elimina primero la tabla intermedia, si existe
	db.Migrator().DropTable("psi_user")

	// Elimina las tablas principales
	db.Migrator().DropTable(&models.PsiUserModel{})
	// db.Migrator().DropTable(&models.SpellsModel{})

	// Crea las tablas principales primero
	err := db.AutoMigrate(&models.PsiUserModel{})
	if err != nil {
		return fmt.Errorf("error al migrar las tablas: %w", err)
	}

	db.Migrator().DropTable("psi_user_col_data")
	db.Migrator().DropTable(&models.PsiUserColData{})

	err = db.AutoMigrate(&models.PsiUserColData{})
	if err != nil {
		return fmt.Errorf("error al migrar las tablas: %w", err)
	}

	db.Migrator().DropTable("psi_specialties")
	db.Migrator().DropTable(&models.PsiSpecialty{})

	err = db.AutoMigrate(&models.PsiSpecialty{})
	if err != nil {
		return fmt.Errorf("error al migrar las tablas: %w", err)
	}

	db.Migrator().DropTable("user_admin")
	db.Migrator().DropTable(&models.UserAdmin{})
	err = db.AutoMigrate(&models.UserAdmin{})
	if err != nil {
		return fmt.Errorf("error al migrar las tablas: %w", err)
	}

	return nil
}
