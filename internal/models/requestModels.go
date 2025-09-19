package models

import (
	"encoding/json"
)

type CreateJobRequest struct {
	Name     string          `json:"name" validate:"required"`
	Payload  json.RawMessage `json:"payload" validate:"required"`
	Schedule string          `json:"schedule" validate:"required"`
	Type     JobType         `json:"type" validate:"required,oneof=http sql queue"`
}

type UpdateJobRequest struct {
	Name     string          `json:"name,omitempty"`
	Payload  json.RawMessage `json:"payload,omitempty"`
	Schedule string          `json:"schedule,omitempty"`
	Type     JobType         `json:"type,omitempty" validate:"omitempty,oneof=http sql queue"`
}

type UserSignUpRequest struct {
	Email     string `json:"email" validate:"required,email"`
	Name      string `json:"name" validate:"required"`
	Password  string `json:"password" validate:"required,min=6"`
	AvatarUrl string `json:"avatar_url,omitempty"`
}

type UserLoginRequest struct {
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required,min=6"`
}
