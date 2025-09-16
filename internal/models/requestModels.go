package models

import "encoding/json"

type CreateJobRequest struct {
	Name     string          `json:"name" validate:"required"`
	Payload  json.RawMessage `json:"payload" validate:"required"`
	Schedule string          `json:"schedule" validate:"required"`
	Type     JobType         `json:"type" validate:"required,oneof=http sql queue"`
}

type UpdateJobRequest struct {
	ID       string          `json:"id" validate:"required"`
	Name     string          `json:"name" validate:"required"`
	Payload  json.RawMessage `json:"payload" validate:"required"`
	Schedule string          `json:"schedule" validate:"required"`
	Type     JobType         `json:"type" validate:"required,oneof=http sql queue"`
}
