package repository

import (
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func init() {
	var err error

	dsn := "host=localhost user=postgres password=@Nicomanu3 dbname=gerenciador_conta port=5432 sslmode=disable"
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal("Failed to connect to the Database")
	}

	// Migration
	DB.AutoMigrate(&ContaModel{}, &Operacao{})

	fmt.Print("Migrations executadas\n")
}
