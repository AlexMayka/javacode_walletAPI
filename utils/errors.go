package utils

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
)

// ErrorResponse defines the standard error response format returned by the API.
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
	Code    int    `json:"code,omitempty"`
}

// Predefined service-level errors used across handlers and services.
var (
	ErrInvalidRequest  = errors.New("invalid request")
	ErrInvalidAmount   = errors.New("amount must be greater than 0")
	ErrNegativeBalance = errors.New("the amount cannot be negative")
	ErrWalletNotFound  = errors.New("wallet not found")
	ErrDatabase        = errors.New("database error")
)

// HandleError maps internal errors to appropriate HTTP responses and sends them via Gin.
//
// It supports specific error types like ErrInvalidRequest, ErrWalletNotFound, etc.,
// and defaults to 500 Internal Server Error if the error is unknown.
func HandleError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, ErrInvalidRequest):
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "invalid_request",
			Message: "Request is invalid or missing required fields",
			Code:    400,
		})
	case errors.Is(err, ErrInvalidAmount):
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "invalid_amount",
			Message: "Amount must be greater than zero",
			Code:    400,
		})
	case errors.Is(err, ErrNegativeBalance):
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "negative_amount",
			Message: "The amount cannot be negative.",
			Code:    400,
		})
	case errors.Is(err, ErrWalletNotFound):
		c.JSON(http.StatusNotFound, ErrorResponse{
			Error:   "wallet_not_found",
			Message: "Wallet not found by uuid",
			Code:    404,
		})
	case errors.Is(err, ErrDatabase):
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "database_error",
			Message: "Database operation failed",
			Code:    500,
		})
	default:
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "internal_server_error",
			Message: "An unexpected error occurred",
			Code:    500,
		})
	}
}
