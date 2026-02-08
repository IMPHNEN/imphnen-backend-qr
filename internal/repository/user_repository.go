package repository

import (
	"database/sql"
	"time"

	"github.com/IMPHNEN/imphnen-backend-qr/internal/domain"
	"github.com/google/uuid"
)

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) domain.UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(user *domain.User) error {
	if user.ID == "" {
		user.ID = uuid.New().String()
	}
	now := time.Now()
	user.CreatedAt = now
	user.UpdatedAt = now

	_, err := r.db.Exec(
		`INSERT INTO users (id, email, password, name, role, provider, provider_id, created_at, updated_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`,
		user.ID, user.Email, user.Password, user.Name, user.Role, user.Provider, user.ProviderID, user.CreatedAt, user.UpdatedAt,
	)
	return err
}

func (r *userRepository) FindByEmail(email string) (*domain.User, error) {
	user := &domain.User{}
	err := r.db.QueryRow(
		`SELECT id, email, password, name, role, provider, provider_id, created_at, updated_at
		 FROM users WHERE email = $1`, email,
	).Scan(&user.ID, &user.Email, &user.Password, &user.Name, &user.Role, &user.Provider, &user.ProviderID, &user.CreatedAt, &user.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return user, err
}

func (r *userRepository) FindByID(id string) (*domain.User, error) {
	if _, err := uuid.Parse(id); err != nil {
		return nil, nil
	}

	user := &domain.User{}
	err := r.db.QueryRow(
		`SELECT id, email, password, name, role, provider, provider_id, created_at, updated_at
		 FROM users WHERE id = $1::uuid`, id,
	).Scan(&user.ID, &user.Email, &user.Password, &user.Name, &user.Role, &user.Provider, &user.ProviderID, &user.CreatedAt, &user.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return user, err
}

func (r *userRepository) FindByProviderID(provider, providerID string) (*domain.User, error) {
	user := &domain.User{}
	err := r.db.QueryRow(
		`SELECT id, email, password, name, role, provider, provider_id, created_at, updated_at
		 FROM users WHERE provider = $1 AND provider_id = $2`, provider, providerID,
	).Scan(&user.ID, &user.Email, &user.Password, &user.Name, &user.Role, &user.Provider, &user.ProviderID, &user.CreatedAt, &user.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return user, err
}

func (r *userRepository) FindAll() ([]domain.User, error) {
	rows, err := r.db.Query(
		`SELECT id, email, password, name, role, provider, provider_id, created_at, updated_at
		 FROM users ORDER BY created_at DESC`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []domain.User
	for rows.Next() {
		var user domain.User
		if err := rows.Scan(&user.ID, &user.Email, &user.Password, &user.Name, &user.Role, &user.Provider, &user.ProviderID, &user.CreatedAt, &user.UpdatedAt); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, rows.Err()
}

func (r *userRepository) Update(user *domain.User) error {
	user.UpdatedAt = time.Now()
	_, err := r.db.Exec(
		`UPDATE users SET email = $1, name = $2, updated_at = $3 WHERE id = $4::uuid`,
		user.Email, user.Name, user.UpdatedAt, user.ID,
	)
	return err
}

func (r *userRepository) UpdateRole(id, role string) error {
	_, err := r.db.Exec(
		`UPDATE users SET role = $1, updated_at = $2 WHERE id = $3::uuid`,
		role, time.Now(), id,
	)
	return err
}

func (r *userRepository) Delete(id string) error {
	if _, err := uuid.Parse(id); err != nil {
		return err
	}

	_, err := r.db.Exec(`DELETE FROM users WHERE id = $1::uuid`, id)
	return err
}
