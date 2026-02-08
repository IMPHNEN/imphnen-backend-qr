package handler

import (
	"errors"
	"log"
	"net/http"

	"github.com/IMPHNEN/imphnen-backend-qr/internal/service"
	"github.com/IMPHNEN/imphnen-backend-qr/internal/utils"
	"github.com/labstack/echo/v4"
)

type UserHandler struct {
	userService *service.UserService
}

func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

func (h *UserHandler) GetProfile(c echo.Context) error {
	userID := c.Get("user_id").(string)

	user, err := h.userService.GetProfile(userID)
	if err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			log.Printf("[WARN] GetProfile: user not found for user_id = %q", userID)
			return utils.ErrorResponse(c, http.StatusNotFound, "user not found", "user_not_found")
		}
		log.Printf("[ERROR] GetProfile: database error for user_id = %q: %v", userID, err)
		return utils.ErrorResponse(c, http.StatusInternalServerError, "internal server error", "internal_error")
	}

	return utils.SuccessResponse(c, http.StatusOK, "profile retrieved", user)
}

func (h *UserHandler) UpdateProfile(c echo.Context) error {
	userID := c.Get("user_id").(string)
	var input service.UpdateProfileInput
	if err := c.Bind(&input); err != nil {
		return utils.ErrorResponse(c, http.StatusBadRequest, "invalid request body", "bad_request")
	}

	user, err := h.userService.UpdateProfile(userID, input)
	if err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			log.Printf("[WARN] UpdateProfile: user not found for user_id = %q", userID)
			return utils.ErrorResponse(c, http.StatusNotFound, "user not found", "user_not_found")
		}
		log.Printf("[ERROR] UpdateProfile: error for user_id = %q: %v", userID, err)
		return utils.ErrorResponse(c, http.StatusBadRequest, err.Error(), "update_failed")
	}

	return utils.SuccessResponse(c, http.StatusOK, "profile updated", user)
}

func (h *UserHandler) GetAllUsers(c echo.Context) error {
	users, err := h.userService.GetAllUsers()
	if err != nil {
		return utils.ErrorResponse(c, http.StatusInternalServerError, "failed to fetch users", "internal_error")
	}

	return utils.SuccessResponse(c, http.StatusOK, "users retrieved", users)
}

func (h *UserHandler) UpdateUserRole(c echo.Context) error {
	id := c.Param("id")

	var body struct {
		Role string `json:"role"`
	}
	if err := c.Bind(&body); err != nil {
		return utils.ErrorResponse(c, http.StatusBadRequest, "invalid request body", "bad_request")
	}

	if err := h.userService.UpdateUserRole(id, body.Role); err != nil {
		return utils.ErrorResponse(c, http.StatusBadRequest, err.Error(), "update_role_failed")
	}

	return utils.SuccessResponse(c, http.StatusOK, "user role updated", nil)
}

func (h *UserHandler) DeleteUser(c echo.Context) error {
	id := c.Param("id")

	if err := h.userService.DeleteUser(id); err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			log.Printf("[WARN] DeleteUser: user not found for id = %q", id)
			return utils.ErrorResponse(c, http.StatusNotFound, "user not found", "user_not_found")
		}
		log.Printf("[ERROR] DeleteUser: %v", err)
		return utils.ErrorResponse(c, http.StatusInternalServerError, "failed to delete user", "internal_error")
	}

	return utils.SuccessResponse(c, http.StatusOK, "user deleted", nil)
}
