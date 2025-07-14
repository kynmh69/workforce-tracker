package handlers

import (
	"database/sql"
	"net/http"
	"strconv"
	"time"
	"workforce-tracker-backend/internal/database"
	"workforce-tracker-backend/internal/models"
	"workforce-tracker-backend/pkg/auth"

	"github.com/labstack/echo/v4"
)

type UserHandler struct {
	DB *database.DB
}

type CreateUserRequest struct {
	Email    string      `json:"email" validate:"required,email"`
	Password string      `json:"password" validate:"required,min=8"`
	Name     string      `json:"name" validate:"required"`
	Role     models.Role `json:"role" validate:"required"`
}

type UpdateUserRequest struct {
	Email    string      `json:"email" validate:"email"`
	Password string      `json:"password,omitempty"`
	Name     string      `json:"name"`
	Role     models.Role `json:"role"`
}

func (h *UserHandler) CreateUser(c echo.Context) error {
	var req CreateUserRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request body")
	}

	// Hash password
	hashedPassword, err := auth.HashPassword(req.Password)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to hash password")
	}

	// Insert user
	query := `INSERT INTO users (email, password_hash, name, role, created_at, updated_at)
			  VALUES ($1, $2, $3, $4, $5, $6)
			  RETURNING id`
	
	now := time.Now()
	var userID int
	err = h.DB.QueryRow(query, req.Email, hashedPassword, req.Name, req.Role, now, now).Scan(&userID)
	if err != nil {
		if err.Error() == `pq: duplicate key value violates unique constraint "users_email_key"` {
			return echo.NewHTTPError(http.StatusConflict, "email already exists")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to create user")
	}

	// Return created user (without password)
	user := models.User{
		ID:        userID,
		Email:     req.Email,
		Name:      req.Name,
		Role:      req.Role,
		CreatedAt: now,
		UpdatedAt: now,
	}

	return c.JSON(http.StatusCreated, user)
}

func (h *UserHandler) GetUsers(c echo.Context) error {
	// Parse query parameters
	page, _ := strconv.Atoi(c.QueryParam("page"))
	if page < 1 {
		page = 1
	}
	
	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	if limit < 1 || limit > 100 {
		limit = 20
	}
	
	offset := (page - 1) * limit
	
	search := c.QueryParam("search")
	role := c.QueryParam("role")

	// Build query
	baseQuery := `SELECT id, email, name, role, created_at, updated_at
				  FROM users 
				  WHERE deleted_at IS NULL`
	countQuery := `SELECT COUNT(*) FROM users WHERE deleted_at IS NULL`
	args := []interface{}{}
	argIndex := 1

	if search != "" {
		baseQuery += ` AND (name ILIKE $` + strconv.Itoa(argIndex) + ` OR email ILIKE $` + strconv.Itoa(argIndex) + `)`
		countQuery += ` AND (name ILIKE $` + strconv.Itoa(argIndex) + ` OR email ILIKE $` + strconv.Itoa(argIndex) + `)`
		args = append(args, "%"+search+"%")
		argIndex++
	}

	if role != "" {
		baseQuery += ` AND role = $` + strconv.Itoa(argIndex)
		countQuery += ` AND role = $` + strconv.Itoa(argIndex)
		args = append(args, role)
		argIndex++
	}

	// Get total count
	var total int
	err := h.DB.QueryRow(countQuery, args...).Scan(&total)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "database error")
	}

	// Get users
	baseQuery += ` ORDER BY created_at DESC LIMIT $` + strconv.Itoa(argIndex) + ` OFFSET $` + strconv.Itoa(argIndex+1)
	args = append(args, limit, offset)

	rows, err := h.DB.Query(baseQuery, args...)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "database error")
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		err := rows.Scan(&user.ID, &user.Email, &user.Name, &user.Role, &user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "scan error")
		}
		users = append(users, user)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"users": users,
		"total": total,
		"page":  page,
		"limit": limit,
	})
}

func (h *UserHandler) GetUser(c echo.Context) error {
	userID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid user ID")
	}

	var user models.User
	query := `SELECT id, email, name, role, created_at, updated_at
			  FROM users 
			  WHERE id = $1 AND deleted_at IS NULL`
	
	err = h.DB.QueryRow(query, userID).Scan(
		&user.ID, &user.Email, &user.Name, &user.Role, &user.CreatedAt, &user.UpdatedAt,
	)
	
	if err == sql.ErrNoRows {
		return echo.NewHTTPError(http.StatusNotFound, "user not found")
	}
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "database error")
	}

	return c.JSON(http.StatusOK, user)
}

func (h *UserHandler) UpdateUser(c echo.Context) error {
	userID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid user ID")
	}

	var req UpdateUserRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request body")
	}

	// Build update query dynamically
	setParts := []string{}
	args := []interface{}{}
	argIndex := 1

	if req.Email != "" {
		setParts = append(setParts, "email = $"+strconv.Itoa(argIndex))
		args = append(args, req.Email)
		argIndex++
	}

	if req.Name != "" {
		setParts = append(setParts, "name = $"+strconv.Itoa(argIndex))
		args = append(args, req.Name)
		argIndex++
	}

	if req.Role != "" {
		setParts = append(setParts, "role = $"+strconv.Itoa(argIndex))
		args = append(args, req.Role)
		argIndex++
	}

	if req.Password != "" {
		hashedPassword, err := auth.HashPassword(req.Password)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "failed to hash password")
		}
		setParts = append(setParts, "password_hash = $"+strconv.Itoa(argIndex))
		args = append(args, hashedPassword)
		argIndex++
	}

	if len(setParts) == 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "no fields to update")
	}

	// Add updated_at
	setParts = append(setParts, "updated_at = $"+strconv.Itoa(argIndex))
	args = append(args, time.Now())
	argIndex++

	// Add WHERE clause
	args = append(args, userID)

	query := `UPDATE users SET ` + 
		setParts[0]
	for i := 1; i < len(setParts); i++ {
		query += ", " + setParts[i]
	}
	query += ` WHERE id = $` + strconv.Itoa(argIndex) + ` AND deleted_at IS NULL`

	result, err := h.DB.Exec(query, args...)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to update user")
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return echo.NewHTTPError(http.StatusNotFound, "user not found")
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "user updated successfully",
	})
}

func (h *UserHandler) DeleteUser(c echo.Context) error {
	userID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid user ID")
	}

	// Soft delete
	query := `UPDATE users SET deleted_at = $1 WHERE id = $2 AND deleted_at IS NULL`
	result, err := h.DB.Exec(query, time.Now(), userID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to delete user")
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return echo.NewHTTPError(http.StatusNotFound, "user not found")
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "user deleted successfully",
	})
}