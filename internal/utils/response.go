package utils

import "github.com/labstack/echo/v4"

type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

func SuccessResponse(c echo.Context, code int, message string, data interface{}) error {
	return c.JSON(code, Response{
		Success: true,
		Message: message,
		Data:    data,
	})
}

func ErrorResponse(c echo.Context, code int, message, errCode string) error {
	return c.JSON(code, Response{
		Success: false,
		Message: message,
		Error:   errCode,
	})
}
