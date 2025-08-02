package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type Jobs struct {
	ID        uuid.UUID      `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	URL       string         `json:"url"`
	Method    string         `json:"method"`
	Headers   datatypes.JSON `json:"headers" gorm:"type:jsonb"`
	Body      datatypes.JSON `json:"body" gorm:"type:jsonb"`
	Schedule  string         `json:"schedule"`
	LastRun   *time.Time     `json:"last_run,omitempty"`
	NextRun   *time.Time     `json:"next_run,omitempty"`
	CreatedAt time.Time      `json:"created_at"`
}

type Logs struct {
	ID       string    `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	JobId    string    `json:"job_id"`
	Status   string    `json:"status"`
	Response string    `json:"response"`
	RunAt    time.Time `json:"run_at"`
}
