package handler

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"

	"github.com/IMPHNEN/imphnen-backend-qr/internal/service"
	"github.com/IMPHNEN/imphnen-backend-qr/internal/utils"
	"github.com/labstack/echo/v4"
)

type AuthHandler struct {
	authService *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

func (h *AuthHandler) Register(c echo.Context) error {
	var input service.RegisterInput
	if err := c.Bind(&input); err != nil {
		return utils.ErrorResponse(c, http.StatusBadRequest, "invalid request body", "bad_request")
	}

	if input.Email == "" || input.Password == "" || input.Name == "" {
		return utils.ErrorResponse(c, http.StatusBadRequest, "email, password, and name are required", "validation_error")
	}

	user, tokens, err := h.authService.Register(input)
	if err != nil {
		return utils.ErrorResponse(c, http.StatusConflict, err.Error(), "register_failed")
	}

	return utils.SuccessResponse(c, http.StatusCreated, "registration successful", map[string]interface{}{
		"user":   user,
		"tokens": tokens,
	})
}

func (h *AuthHandler) Login(c echo.Context) error {
	var input service.LoginInput
	if err := c.Bind(&input); err != nil {
		return utils.ErrorResponse(c, http.StatusBadRequest, "invalid request body", "bad_request")
	}

	if input.Email == "" || input.Password == "" {
		return utils.ErrorResponse(c, http.StatusBadRequest, "email and password are required", "validation_error")
	}

	user, tokens, err := h.authService.Login(input)
	if err != nil {
		return utils.ErrorResponse(c, http.StatusUnauthorized, err.Error(), "login_failed")
	}

	return utils.SuccessResponse(c, http.StatusOK, "login successful", map[string]interface{}{
		"user":   user,
		"tokens": tokens,
	})
}

func (h *AuthHandler) GoogleAuth(c echo.Context) error {
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	state := hex.EncodeToString(b)

	url, err := h.authService.GoogleAuthURL(state)
	if err != nil {
		return utils.ErrorResponse(c, http.StatusInternalServerError, err.Error(), "oauth_error")
	}

	return c.Redirect(http.StatusTemporaryRedirect, url)
}

func (h *AuthHandler) GoogleCallback(c echo.Context) error {
	code := c.QueryParam("code")
	if code == "" {
		return utils.ErrorResponse(c, http.StatusBadRequest, "missing code parameter", "bad_request")
	}

	user, tokens, err := h.authService.GoogleCallback(code)
	if err != nil {
		return utils.ErrorResponse(c, http.StatusUnauthorized, err.Error(), "oauth_failed")
	}

	return utils.SuccessResponse(c, http.StatusOK, "google login successful", map[string]interface{}{
		"user":   user,
		"tokens": tokens,
	})
}

func (h *AuthHandler) RefreshToken(c echo.Context) error {
	var body struct {
		RefreshToken string `json:"refresh_token"`
	}
	if err := c.Bind(&body); err != nil {
		return utils.ErrorResponse(c, http.StatusBadRequest, "invalid request body", "bad_request")
	}

	if body.RefreshToken == "" {
		return utils.ErrorResponse(c, http.StatusBadRequest, "refresh_token is required", "validation_error")
	}

	tokens, err := h.authService.RefreshToken(body.RefreshToken)
	if err != nil {
		return utils.ErrorResponse(c, http.StatusUnauthorized, err.Error(), "refresh_failed")
	}

	return utils.SuccessResponse(c, http.StatusOK, "token refreshed", tokens)
}
