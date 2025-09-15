package models

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type JobType string

const (
	JobTypeHTTP JobType = "http"
	// JobTypeShell JobType = "shell"
	JobTypeSQL   JobType = "sql"
	JobTypeQueue JobType = "queue"
)

type Job struct {
	ID        uuid.UUID       `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	LastRun   *time.Time      `json:"last_run,omitempty"`
	NextRun   *time.Time      `json:"next_run,omitempty"`
	CreatedAt time.Time       `json:"created_at"`
	Schedule  string          `json:"schedule" gorm:"default:'* * * * *'"`
	Name      string          `json:"name"`
	Payload   json.RawMessage `json:"payload"` // one-time or recurring
	Type      JobType         `json:"type"`
}

type Logs struct {
	ID       uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	JobId    uuid.UUID `json:"job_id" gorm:"type:uuid"`
	Status   int       `json:"status"`
	Response string    `json:"response"`
	RunAt    time.Time `json:"run_at"`
}
