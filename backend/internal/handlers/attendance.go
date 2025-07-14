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

type AttendanceHandler struct {
	DB *database.DB
}

func (h *AttendanceHandler) ClockIn(c echo.Context) error {
	user, ok := c.Get("user").(*auth.Claims)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "user not found in context")
	}

	now := time.Now()
	today := now.Truncate(24 * time.Hour)

	// Check if already clocked in today
	var existingID int
	checkQuery := `SELECT id FROM attendances 
				   WHERE user_id = $1 AND date = $2 AND clock_in_time IS NOT NULL AND deleted_at IS NULL`
	err := h.DB.QueryRow(checkQuery, user.UserID, today).Scan(&existingID)
	
	if err == nil {
		return echo.NewHTTPError(http.StatusBadRequest, "already clocked in today")
	}
	if err != sql.ErrNoRows {
		return echo.NewHTTPError(http.StatusInternalServerError, "database error")
	}

	// Insert or update attendance record
	query := `INSERT INTO attendances (user_id, date, clock_in_time, created_at, updated_at)
			  VALUES ($1, $2, $3, $4, $5)
			  ON CONFLICT (user_id, date) 
			  DO UPDATE SET clock_in_time = $3, updated_at = $5
			  RETURNING id`
	
	var attendanceID int
	err = h.DB.QueryRow(query, user.UserID, today, now, now, now).Scan(&attendanceID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to record clock in")
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message":        "clocked in successfully",
		"attendance_id":  attendanceID,
		"clock_in_time":  now,
	})
}

func (h *AttendanceHandler) ClockOut(c echo.Context) error {
	user, ok := c.Get("user").(*auth.Claims)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "user not found in context")
	}

	now := time.Now()
	today := now.Truncate(24 * time.Hour)

	// Check if clocked in today and not already clocked out
	var attendanceID int
	var clockInTime time.Time
	checkQuery := `SELECT id, clock_in_time FROM attendances 
				   WHERE user_id = $1 AND date = $2 AND clock_in_time IS NOT NULL 
				   AND clock_out_time IS NULL AND deleted_at IS NULL`
	
	err := h.DB.QueryRow(checkQuery, user.UserID, today).Scan(&attendanceID, &clockInTime)
	if err == sql.ErrNoRows {
		return echo.NewHTTPError(http.StatusBadRequest, "not clocked in or already clocked out")
	}
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "database error")
	}

	// Calculate work hours
	duration := now.Sub(clockInTime)
	workHours := duration.Hours()

	// Update attendance record
	query := `UPDATE attendances 
			  SET clock_out_time = $1, work_hours = $2, updated_at = $3
			  WHERE id = $4`
	
	_, err = h.DB.Exec(query, now, workHours, now, attendanceID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to record clock out")
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message":         "clocked out successfully",
		"attendance_id":   attendanceID,
		"clock_out_time":  now,
		"work_hours":      workHours,
	})
}

func (h *AttendanceHandler) GetToday(c echo.Context) error {
	user, ok := c.Get("user").(*auth.Claims)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "user not found in context")
	}

	today := time.Now().Truncate(24 * time.Hour)

	var attendance models.Attendance
	query := `SELECT id, user_id, date, clock_in_time, clock_out_time, work_hours, created_at, updated_at
			  FROM attendances 
			  WHERE user_id = $1 AND date = $2 AND deleted_at IS NULL`
	
	err := h.DB.QueryRow(query, user.UserID, today).Scan(
		&attendance.ID, &attendance.UserID, &attendance.Date,
		&attendance.ClockInTime, &attendance.ClockOutTime, &attendance.WorkHours,
		&attendance.CreatedAt, &attendance.UpdatedAt,
	)
	
	if err == sql.ErrNoRows {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"date":           today,
			"clock_in_time":  nil,
			"clock_out_time": nil,
			"work_hours":     nil,
		})
	}
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "database error")
	}

	return c.JSON(http.StatusOK, attendance)
}

func (h *AttendanceHandler) GetHistory(c echo.Context) error {
	user, ok := c.Get("user").(*auth.Claims)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "user not found in context")
	}

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

	// Check if user is admin and requesting another user's data
	targetUserID := user.UserID
	if user.Role == models.RoleAdmin && c.QueryParam("user_id") != "" {
		if uid, err := strconv.Atoi(c.QueryParam("user_id")); err == nil {
			targetUserID = uid
		}
	}

	query := `SELECT id, user_id, date, clock_in_time, clock_out_time, work_hours, created_at, updated_at
			  FROM attendances 
			  WHERE user_id = $1 AND deleted_at IS NULL
			  ORDER BY date DESC
			  LIMIT $2 OFFSET $3`

	rows, err := h.DB.Query(query, targetUserID, limit, offset)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "database error")
	}
	defer rows.Close()

	var attendances []models.Attendance
	for rows.Next() {
		var attendance models.Attendance
		err := rows.Scan(
			&attendance.ID, &attendance.UserID, &attendance.Date,
			&attendance.ClockInTime, &attendance.ClockOutTime, &attendance.WorkHours,
			&attendance.CreatedAt, &attendance.UpdatedAt,
		)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "scan error")
		}
		attendances = append(attendances, attendance)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"attendances": attendances,
		"page":        page,
		"limit":       limit,
	})
}