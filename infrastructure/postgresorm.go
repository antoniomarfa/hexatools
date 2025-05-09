package infrastructure

import (
	"context"
	"fmt"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectPostgresOrm(ctx context.Context, dsn string) (*gorm.DB, error) {
	db, _ := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	return db, pingSqlOrm(db)
}

func pingSqlOrm(db *gorm.DB) (err error) {
	// wait until db is ready
	for start := time.Now(); time.Since(start) < (5 * time.Second); {

		dbSQL, err := db.DB()
		if err != nil {
			break
		}

		// Hacemos un ping para asegurarnos de que la base de datos está viva
		err = dbSQL.Ping()
		if err != nil {
			break
		}

		time.Sleep(1 * time.Second)
	}
	return err
}

type PostgresRepositoryOrm struct {
	DB *gorm.DB
}

var SharedTables = map[string]string{
	"company": "travel",
}

func GetDBWithSchema(ctx context.Context, db *gorm.DB) *gorm.DB {
	// Obtener el esquema desde el contexto
	schemaVal, exists := ctx.Value("schema").(string)
	if !exists || schemaVal == "" {
		return db // Si no hay esquema, devolver la conexión por defecto
	}

	newDB := db.Session(&gorm.Session{NewDB: true})

	// Ejecutar SET search_path en la nueva sesión, pero sin devolver el resultado de Exec
	if err := newDB.Exec("SET search_path TO " + schemaVal).Error; err != nil {
		// Puedes manejar el error aquí si lo deseas, o dejarlo pasar silenciosamente
		fmt.Println("error al aplicar search_path:", err)
	}

	return newDB
}
