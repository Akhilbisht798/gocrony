package models

import "encoding/json"

type CreateJobRequest struct {
	Name     string          `json:"url" validate:"required"`
	Payload  json.RawMessage `json:"method" validate:"required"`
	Schedule string          `json:"schedule" validate:"required"`
	Type     JobType         `json:"type" validate:"required,oneof=http sql queue"`
}

type UpdateJobRequest struct {
	ID       string          `json:"id" validate:"required"`
	Name     string          `json:"url" validate:"required"`
	Payload  json.RawMessage `json:"method" validate:"required"`
	Schedule string          `json:"schedule" validate:"required"`
	Type     JobType         `json:"type" validate:"required,oneof=http sql queue"`
}
