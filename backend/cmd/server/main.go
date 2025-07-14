package main

import (
	"log"
	"workforce-tracker-backend/internal/config"
	"workforce-tracker-backend/internal/database"
	"workforce-tracker-backend/internal/handlers"
	"workforce-tracker-backend/internal/middleware"
	"workforce-tracker-backend/pkg/auth"

	"github.com/labstack/echo/v4"
	echomiddleware "github.com/labstack/echo/v4/middleware"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Connect to database
	db, err := database.NewConnection(cfg.DatabaseURL)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	// Create tables
	if err := db.CreateTables(); err != nil {
		log.Fatal("Failed to create tables:", err)
	}

	// Create default admin user if not exists
	if err := createDefaultAdmin(db); err != nil {
		log.Fatal("Failed to create default admin:", err)
	}

	// Initialize Echo
	e := echo.New()

	// Middleware
	e.Use(echomiddleware.Logger())
	e.Use(echomiddleware.Recover())
	e.Use(echomiddleware.CORS())

	// Initialize handlers
	authHandler := &handlers.AuthHandler{
		DB:        db,
		JWTSecret: cfg.JWTSecret,
	}
	userHandler := &handlers.UserHandler{DB: db}
	attendanceHandler := &handlers.AttendanceHandler{DB: db}

	// Routes
	api := e.Group("/api")

	// Auth routes (public)
	auth := api.Group("/auth")
	auth.POST("/login", authHandler.Login)
	auth.POST("/logout", authHandler.Logout)
	
	// Protected routes
	protected := api.Group("", middleware.JWTMiddleware(cfg.JWTSecret))
	protected.GET("/auth/me", authHandler.Me)

	// Attendance routes (protected)
	attendance := protected.Group("/attendance")
	attendance.POST("/clock-in", attendanceHandler.ClockIn)
	attendance.POST("/clock-out", attendanceHandler.ClockOut)
	attendance.GET("/today", attendanceHandler.GetToday)
	attendance.GET("/history", attendanceHandler.GetHistory)

	// User management routes (admin only)
	admin := protected.Group("/users", middleware.AdminOnly())
	admin.POST("", userHandler.CreateUser)
	admin.GET("", userHandler.GetUsers)
	admin.GET("/:id", userHandler.GetUser)
	admin.PUT("/:id", userHandler.UpdateUser)
	admin.DELETE("/:id", userHandler.DeleteUser)

	// Start server
	log.Printf("Server starting on port %s", cfg.Port)
	log.Fatal(e.Start(":" + cfg.Port))
}

func createDefaultAdmin(db *database.DB) error {
	// Check if admin exists
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM users WHERE role = 'admin' AND deleted_at IS NULL").Scan(&count)
	if err != nil {
		return err
	}

	if count > 0 {
		return nil // Admin already exists
	}

	// Create default admin
	hashedPassword, err := auth.HashPassword("admin123")
	if err != nil {
		return err
	}

	query := `INSERT INTO users (email, password_hash, name, role)
			  VALUES ($1, $2, $3, $4)`
	
	_, err = db.Exec(query, "admin@example.com", hashedPassword, "Administrator", "admin")
	if err != nil {
		return err
	}

	log.Println("Default admin user created:")
	log.Println("  Email: admin@example.com")
	log.Println("  Password: admin123")
	log.Println("  Please change the password after first login!")

	return nil
}