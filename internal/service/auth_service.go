package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/IMPHNEN/imphnen-backend-qr/internal/config"
	"github.com/IMPHNEN/imphnen-backend-qr/internal/domain"
	"github.com/IMPHNEN/imphnen-backend-qr/internal/utils"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type AuthService struct {
	userRepo    domain.UserRepository
	cfg         *config.Config
	googleOAuth *oauth2.Config
}

func NewAuthService(userRepo domain.UserRepository, cfg *config.Config) *AuthService {
	var googleOAuth *oauth2.Config
	if cfg.GoogleClientID != "" && cfg.GoogleClientSecret != "" {
		googleOAuth = &oauth2.Config{
			ClientID:     cfg.GoogleClientID,
			ClientSecret: cfg.GoogleClientSecret,
			RedirectURL:  cfg.GoogleRedirectURL,
			Scopes:       []string{"openid", "email", "profile"},
			Endpoint:     google.Endpoint,
		}
	}

	return &AuthService{
		userRepo:    userRepo,
		cfg:         cfg,
		googleOAuth: googleOAuth,
	}
}

type RegisterInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
}

type LoginInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (s *AuthService) Register(input RegisterInput) (*domain.User, *domain.TokenPair, error) {
	existing, err := s.userRepo.FindByEmail(input.Email)
	if err != nil {
		return nil, nil, err
	}
	if existing != nil {
		return nil, nil, errors.New("email already registered")
	}

	hashed, err := utils.HashPassword(input.Password)
	if err != nil {
		return nil, nil, err
	}

	user := &domain.User{
		Email:    input.Email,
		Password: &hashed,
		Name:     input.Name,
		Role:     "user",
		Provider: "local",
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, nil, err
	}

	accessToken, refreshToken, err := utils.GenerateTokenPair(s.cfg.JWTSecret, user.ID, user.Email, user.Role)
	if err != nil {
		return nil, nil, err
	}

	return user, &domain.TokenPair{AccessToken: accessToken, RefreshToken: refreshToken}, nil
}

func (s *AuthService) Login(input LoginInput) (*domain.User, *domain.TokenPair, error) {
	user, err := s.userRepo.FindByEmail(input.Email)
	if err != nil {
		return nil, nil, err
	}
	if user == nil {
		return nil, nil, errors.New("invalid credentials")
	}

	if user.Password == nil {
		return nil, nil, errors.New("this account uses Google login")
	}

	if err := utils.ComparePassword(*user.Password, input.Password); err != nil {
		return nil, nil, errors.New("invalid credentials")
	}

	accessToken, refreshToken, err := utils.GenerateTokenPair(s.cfg.JWTSecret, user.ID, user.Email, user.Role)
	if err != nil {
		return nil, nil, err
	}

	return user, &domain.TokenPair{AccessToken: accessToken, RefreshToken: refreshToken}, nil
}

func (s *AuthService) GoogleAuthURL(state string) (string, error) {
	if s.googleOAuth == nil {
		return "", errors.New("google oauth not configured")
	}
	return s.googleOAuth.AuthCodeURL(state, oauth2.AccessTypeOffline), nil
}

type googleUserInfo struct {
	ID    string `json:"sub"`
	Email string `json:"email"`
	Name  string `json:"name"`
}

func (s *AuthService) GoogleCallback(code string) (*domain.User, *domain.TokenPair, error) {
	if s.googleOAuth == nil {
		return nil, nil, errors.New("google oauth not configured")
	}

	token, err := s.googleOAuth.Exchange(context.Background(), code)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to exchange code: %w", err)
	}

	client := s.googleOAuth.Client(context.Background(), token)
	resp, err := client.Get("https://openidconnect.googleapis.com/v1/userinfo")
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get user info: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, nil, errors.New("failed to get user info from google")
	}

	var gUser googleUserInfo
	if err := json.NewDecoder(resp.Body).Decode(&gUser); err != nil {
		return nil, nil, err
	}

	// Check if user exists by provider ID
	user, err := s.userRepo.FindByProviderID("google", gUser.ID)
	if err != nil {
		return nil, nil, err
	}

	if user == nil {
		// Check if email already used with local account
		user, err = s.userRepo.FindByEmail(gUser.Email)
		if err != nil {
			return nil, nil, err
		}

		if user != nil {
			return nil, nil, errors.New("email already registered with local account")
		}

		// Create new user
		providerID := gUser.ID
		user = &domain.User{
			Email:      gUser.Email,
			Name:       gUser.Name,
			Role:       "user",
			Provider:   "google",
			ProviderID: &providerID,
		}
		if err := s.userRepo.Create(user); err != nil {
			return nil, nil, err
		}
	}

	accessToken, refreshToken, err := utils.GenerateTokenPair(s.cfg.JWTSecret, user.ID, user.Email, user.Role)
	if err != nil {
		return nil, nil, err
	}

	return user, &domain.TokenPair{AccessToken: accessToken, RefreshToken: refreshToken}, nil
}

func (s *AuthService) RefreshToken(refreshTokenStr string) (*domain.TokenPair, error) {
	claims, err := utils.ValidateToken(s.cfg.JWTSecret, refreshTokenStr)
	if err != nil {
		return nil, errors.New("invalid refresh token")
	}

	accessToken, refreshToken, err := utils.GenerateTokenPair(s.cfg.JWTSecret, claims.UserID, claims.Email, claims.Role)
	if err != nil {
		return nil, err
	}

	return &domain.TokenPair{AccessToken: accessToken, RefreshToken: refreshToken}, nil
}
