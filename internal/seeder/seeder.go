package seeder

import (
	"database/sql"
	"log"

	"github.com/IMPHNEN/imphnen-backend-qr/internal/domain"
	"github.com/IMPHNEN/imphnen-backend-qr/internal/repository"
	"github.com/IMPHNEN/imphnen-backend-qr/internal/utils"
)

type seedUser struct {
	Email    string
	Password string
	Name     string
	Role     string
}

func Run(db *sql.DB) {
	userRepo := repository.NewUserRepository(db)

	users := []seedUser{
		{
			Email:    "admin@imphnen.dev",
			Password: "admin123",
			Name:     "Admin Demo",
			Role:     "admin",
		},
		{
			Email:    "user@imphnen.dev",
			Password: "user123",
			Name:     "User Demo",
			Role:     "user",
		},
	}

	for _, u := range users {
		existing, err := userRepo.FindByEmail(u.Email)
		if err != nil {
			log.Printf("[SEEDER] Error checking user %s: %v", u.Email, err)
			continue
		}
		if existing != nil {
			log.Printf("[SEEDER] User %s already exists, skipping", u.Email)
			continue
		}

		hashed, err := utils.HashPassword(u.Password)
		if err != nil {
			log.Printf("[SEEDER] Error hashing password for %s: %v", u.Email, err)
			continue
		}

		user := &domain.User{
			Email:    u.Email,
			Password: &hashed,
			Name:     u.Name,
			Role:     u.Role,
			Provider: "local",
		}

		if err := userRepo.Create(user); err != nil {
			log.Printf("[SEEDER] Error creating user %s: %v", u.Email, err)
			continue
		}

		log.Printf("[SEEDER] Created user %s (role: %s)", u.Email, u.Role)
	}
}
