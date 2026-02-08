package service

import (
	"errors"

	"github.com/IMPHNEN/imphnen-backend-qr/internal/domain"
)

var ErrUserNotFound = errors.New("user not found")

type UserService struct {
	userRepo domain.UserRepository
}

func NewUserService(userRepo domain.UserRepository) *UserService {
	return &UserService{userRepo: userRepo}
}

func (s *UserService) GetProfile(userID string) (*domain.User, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrUserNotFound
	}
	return user, nil
}

type UpdateProfileInput struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

func (s *UserService) UpdateProfile(userID string, input UpdateProfileInput) (*domain.User, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrUserNotFound
	}

	if input.Name != "" {
		user.Name = input.Name
	}
	if input.Email != "" {
		user.Email = input.Email
	}

	if err := s.userRepo.Update(user); err != nil {
		return nil, err
	}
	return user, nil
}

func (s *UserService) GetAllUsers() ([]domain.User, error) {
	return s.userRepo.FindAll()
}

func (s *UserService) UpdateUserRole(userID, role string) error {
	if role != "admin" && role != "user" {
		return errors.New("invalid role, must be 'admin' or 'user'")
	}

	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return err
	}
	if user == nil {
		return ErrUserNotFound
	}

	return s.userRepo.UpdateRole(userID, role)
}

func (s *UserService) DeleteUser(userID string) error {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return err
	}
	if user == nil {
		return ErrUserNotFound
	}

	return s.userRepo.Delete(userID)
}
