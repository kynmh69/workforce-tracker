package handlers

import (
	"database/sql"
	"net/http"
	"workforce-tracker-backend/internal/database"
	"workforce-tracker-backend/internal/models"
	"workforce-tracker-backend/pkg/auth"

	"github.com/labstack/echo/v4"
)

type AuthHandler struct {
	DB        *database.DB
	JWTSecret string
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type LoginResponse struct {
	Token string       `json:"token"`
	User  *models.User `json:"user"`
}

func (h *AuthHandler) Login(c echo.Context) error {
	var req LoginRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request body")
	}

	// Get user from database
	var user models.User
	query := `SELECT id, email, password_hash, name, role, created_at, updated_at 
			  FROM users 
			  WHERE email = $1 AND deleted_at IS NULL`
	
	err := h.DB.QueryRow(query, req.Email).Scan(
		&user.ID, &user.Email, &user.PasswordHash, &user.Name, &user.Role,
		&user.CreatedAt, &user.UpdatedAt,
	)
	
	if err == sql.ErrNoRows {
		return echo.NewHTTPError(http.StatusUnauthorized, "invalid credentials")
	}
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "database error")
	}

	// Check password
	if !auth.CheckPasswordHash(req.Password, user.PasswordHash) {
		return echo.NewHTTPError(http.StatusUnauthorized, "invalid credentials")
	}

	// Generate JWT token
	token, err := auth.GenerateToken(&user, h.JWTSecret)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to generate token")
	}

	// Clear password hash from response
	user.PasswordHash = ""

	return c.JSON(http.StatusOK, LoginResponse{
		Token: token,
		User:  &user,
	})
}

func (h *AuthHandler) Logout(c echo.Context) error {
	// In a JWT-based auth system, logout is typically handled client-side
	// by removing the token from storage
	return c.JSON(http.StatusOK, map[string]string{
		"message": "logged out successfully",
	})
}

func (h *AuthHandler) Me(c echo.Context) error {
	user, ok := c.Get("user").(*auth.Claims)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "user not found in context")
	}

	// Get full user data from database
	var fullUser models.User
	query := `SELECT id, email, name, role, created_at, updated_at 
			  FROM users 
			  WHERE id = $1 AND deleted_at IS NULL`
	
	err := h.DB.QueryRow(query, user.UserID).Scan(
		&fullUser.ID, &fullUser.Email, &fullUser.Name, &fullUser.Role,
		&fullUser.CreatedAt, &fullUser.UpdatedAt,
	)
	
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "database error")
	}

	return c.JSON(http.StatusOK, fullUser)
}