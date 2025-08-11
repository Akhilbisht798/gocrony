package models

type CreateJobRequest struct {
	URL      string                 `json:"url" validate:"required"`
	Method   string                 `json:"method" validate:"required"`
	Headers  map[string]string      `json:"headers" validate:"required"`
	Body     map[string]interface{} `json:"body"`
	Schedule string                 `json:"schedule" validate:"required"`
}

type UpdateJobRequest struct {
	ID       string                 `json:"id" validate:"required"`
	URL      string                 `json:"url"`
	Method   string                 `json:"method"`
	Headers  map[string]string      `json:"headers"`
	Body     map[string]interface{} `json:"body"`
	Schedule string                 `json:"schedule"`
}
