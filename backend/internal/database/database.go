package database

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

type DB struct {
	*sql.DB
}

func NewConnection(databaseURL string) (*DB, error) {
	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &DB{db}, nil
}

func (db *DB) CreateTables() error {
	queries := []string{
		`CREATE EXTENSION IF NOT EXISTS "uuid-ossp";`,
		
		`CREATE TYPE user_role AS ENUM ('admin', 'general');`,
		
		`CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			email VARCHAR(255) UNIQUE NOT NULL,
			password_hash VARCHAR(255) NOT NULL,
			name VARCHAR(255) NOT NULL,
			role user_role NOT NULL DEFAULT 'general',
			created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
			updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
			deleted_at TIMESTAMP WITH TIME ZONE NULL
		);`,
		
		`CREATE TABLE IF NOT EXISTS attendances (
			id SERIAL PRIMARY KEY,
			user_id INTEGER NOT NULL REFERENCES users(id),
			date DATE NOT NULL,
			clock_in_time TIMESTAMP WITH TIME ZONE NULL,
			clock_out_time TIMESTAMP WITH TIME ZONE NULL,
			work_hours DECIMAL(5,2) NULL,
			created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
			updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
			deleted_at TIMESTAMP WITH TIME ZONE NULL,
			UNIQUE(user_id, date)
		);`,
		
		`CREATE TYPE modification_request_type AS ENUM ('clock_in', 'clock_out');`,
		`CREATE TYPE modification_request_status AS ENUM ('pending', 'approved', 'rejected');`,
		
		`CREATE TABLE IF NOT EXISTS modification_requests (
			id SERIAL PRIMARY KEY,
			user_id INTEGER NOT NULL REFERENCES users(id),
			attendance_id INTEGER NOT NULL REFERENCES attendances(id),
			type modification_request_type NOT NULL,
			original_time TIMESTAMP WITH TIME ZONE NULL,
			requested_time TIMESTAMP WITH TIME ZONE NOT NULL,
			reason TEXT NOT NULL,
			status modification_request_status NOT NULL DEFAULT 'pending',
			approved_by INTEGER NULL REFERENCES users(id),
			approved_at TIMESTAMP WITH TIME ZONE NULL,
			rejection_reason TEXT NULL,
			created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
			updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
		);`,
		
		`CREATE TABLE IF NOT EXISTS audit_logs (
			id SERIAL PRIMARY KEY,
			user_id INTEGER NOT NULL REFERENCES users(id),
			action VARCHAR(255) NOT NULL,
			table_name VARCHAR(255) NOT NULL,
			record_id INTEGER NOT NULL,
			old_values JSONB NULL,
			new_values JSONB NULL,
			created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
		);`,
		
		`CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);`,
		`CREATE INDEX IF NOT EXISTS idx_attendances_user_id ON attendances(user_id);`,
		`CREATE INDEX IF NOT EXISTS idx_attendances_date ON attendances(date);`,
		`CREATE INDEX IF NOT EXISTS idx_modification_requests_user_id ON modification_requests(user_id);`,
		`CREATE INDEX IF NOT EXISTS idx_modification_requests_status ON modification_requests(status);`,
		`CREATE INDEX IF NOT EXISTS idx_audit_logs_user_id ON audit_logs(user_id);`,
	}

	for _, query := range queries {
		if _, err := db.Exec(query); err != nil {
			log.Printf("Error executing query: %s\n%v", query, err)
			return err
		}
	}

	return nil
}