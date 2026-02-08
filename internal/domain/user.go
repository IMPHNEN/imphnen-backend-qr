package domain

import "time"

type User struct {
	ID         string    `json:"id"`
	Email      string    `json:"email"`
	Password   *string   `json:"-"`
	Name       string    `json:"name"`
	Role       string    `json:"role"`
	Provider   string    `json:"provider"`
	ProviderID *string   `json:"provider_id,omitempty"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type UserRepository interface {
	Create(user *User) error
	FindByEmail(email string) (*User, error)
	FindByID(id string) (*User, error)
	FindByProviderID(provider, providerID string) (*User, error)
	FindAll() ([]User, error)
	Update(user *User) error
	UpdateRole(id, role string) error
	Delete(id string) error
}
