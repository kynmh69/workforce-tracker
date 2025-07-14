package models

import (
	"database/sql/driver"
	"fmt"
	"time"
)

type Role string

const (
	RoleAdmin   Role = "admin"
	RoleGeneral Role = "general"
)

func (r Role) String() string {
	return string(r)
}

func (r *Role) Scan(value interface{}) error {
	switch s := value.(type) {
	case string:
		*r = Role(s)
	case []byte:
		*r = Role(s)
	case nil:
		*r = ""
	default:
		return fmt.Errorf("cannot scan %T into Role", value)
	}
	return nil
}

func (r Role) Value() (driver.Value, error) {
	return string(r), nil
}

type User struct {
	ID           int       `json:"id" db:"id"`
	Email        string    `json:"email" db:"email"`
	PasswordHash string    `json:"-" db:"password_hash"`
	Name         string    `json:"name" db:"name"`
	Role         Role      `json:"role" db:"role"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
	DeletedAt    *time.Time `json:"deleted_at,omitempty" db:"deleted_at"`
}

type Attendance struct {
	ID           int        `json:"id" db:"id"`
	UserID       int        `json:"user_id" db:"user_id"`
	Date         time.Time  `json:"date" db:"date"`
	ClockInTime  *time.Time `json:"clock_in_time" db:"clock_in_time"`
	ClockOutTime *time.Time `json:"clock_out_time" db:"clock_out_time"`
	WorkHours    *float64   `json:"work_hours" db:"work_hours"`
	CreatedAt    time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at" db:"updated_at"`
	DeletedAt    *time.Time `json:"deleted_at,omitempty" db:"deleted_at"`
}

type ModificationRequestType string

const (
	ModificationTypeClockIn  ModificationRequestType = "clock_in"
	ModificationTypeClockOut ModificationRequestType = "clock_out"
)

type ModificationRequestStatus string

const (
	ModificationStatusPending  ModificationRequestStatus = "pending"
	ModificationStatusApproved ModificationRequestStatus = "approved"
	ModificationStatusRejected ModificationRequestStatus = "rejected"
)

type ModificationRequest struct {
	ID              int                       `json:"id" db:"id"`
	UserID          int                       `json:"user_id" db:"user_id"`
	AttendanceID    int                       `json:"attendance_id" db:"attendance_id"`
	Type            ModificationRequestType   `json:"type" db:"type"`
	OriginalTime    *time.Time                `json:"original_time" db:"original_time"`
	RequestedTime   time.Time                 `json:"requested_time" db:"requested_time"`
	Reason          string                    `json:"reason" db:"reason"`
	Status          ModificationRequestStatus `json:"status" db:"status"`
	ApprovedBy      *int                      `json:"approved_by" db:"approved_by"`
	ApprovedAt      *time.Time                `json:"approved_at" db:"approved_at"`
	RejectionReason *string                   `json:"rejection_reason" db:"rejection_reason"`
	CreatedAt       time.Time                 `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time                 `json:"updated_at" db:"updated_at"`
}

type AuditLog struct {
	ID        int       `json:"id" db:"id"`
	UserID    int       `json:"user_id" db:"user_id"`
	Action    string    `json:"action" db:"action"`
	TableName string    `json:"table_name" db:"table_name"`
	RecordID  int       `json:"record_id" db:"record_id"`
	OldValues string    `json:"old_values" db:"old_values"` // JSON string
	NewValues string    `json:"new_values" db:"new_values"` // JSON string
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}