package main

import (
	"log"

	"github.com/IMPHNEN/imphnen-backend-qr/internal/config"
	"github.com/IMPHNEN/imphnen-backend-qr/internal/seeder"
	"github.com/IMPHNEN/imphnen-backend-qr/pkg/database"
)

func main() {
	cfg := config.Load()
	db := database.NewPostgres(cfg.DatabaseURL)
	defer db.Close()

	log.Println("Running seeder...")
	seeder.Run(db)
	log.Println("Seeder completed")
}
