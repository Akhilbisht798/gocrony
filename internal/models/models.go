package models

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type JobType string

const (
	JobTypeHTTP JobType = "http"
	// JobTypeShell JobType = "shell"
	JobTypeSQL   JobType = "sql"
	JobTypeQueue JobType = "queue"
)

type Job struct {
	ID        uuid.UUID       `gorm:"type:uuid;primaryKey" json:"id"`
	LastRun   *time.Time      `json:"last_run,omitempty"`
	NextRun   *time.Time      `json:"next_run,omitempty"`
	CreatedAt time.Time       `json:"created_at"`
	Schedule  string          `json:"schedule" gorm:"default:'* * * * *'"`
	Name      string          `json:"name"`
	Payload   json.RawMessage `json:"payload"` // one-time or recurring
	Type      JobType         `json:"type"`
}

type Logs struct {
	ID       uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	JobId    uuid.UUID `json:"job_id" gorm:"type:uuid"`
	Status   int       `json:"status"`
	Response string    `json:"response"`
	RunAt    time.Time `json:"run_at"`
}

type User struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey"`
	Email     string    `gorm:"uniqueIndex"`
	Name      string
	AvatarUrl string
	CreatedAt time.Time

	Identities []UserIdentity
}

type UserIdentity struct {
	ID           uuid.UUID `gorm:"type:uuid;primaryKey"`
	UserID       uuid.UUID `gorm:"type:uuid;index"`
	User         User      `gorm:"foreignKey:UserID;references:ID"` // <-- add this
	Provider     string    // e.g., "google", "email"
	ProviderID   string    // e.g., Google sub ID, or email
	PasswordHash string    // only used if provider == "email"

}

func (job *Job) BeforeCreate(tx *gorm.DB) (err error) {
	job.ID = uuid.New()
	return
}

func (log *Logs) BeforeCreate(tx *gorm.DB) (err error) {
	log.ID = uuid.New()
	return
}

func (user *User) BeforeCreate(tx *gorm.DB) (err error) {
	user.ID = uuid.New()
	return
}

func (userIdentity *UserIdentity) BeforeCreate(tx *gorm.DB) (err error) {
	userIdentity.ID = uuid.New()
	return
}
