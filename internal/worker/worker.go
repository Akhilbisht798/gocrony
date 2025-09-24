package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/akhilbisht798/gocrony/internal/cache"
	"github.com/akhilbisht798/gocrony/internal/db"
	"github.com/akhilbisht798/gocrony/internal/models"
	"github.com/akhilbisht798/gocrony/internal/scheduler"
	"github.com/google/uuid"
)

const MAX_RETRY = 3

type Worker struct {
	ID     string
	client *http.Client
}

func NewWorker(id string) *Worker {
	return &Worker{
		ID:     id,
		client: &http.Client{},
	}
}

func (w *Worker) Start(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			log.Println("Worker stop due to context cancellation", w.ID)
			return
		default:
		}
		//TODO: change in future to check redis is working or nnot
		// otherwise this will be blocked forever.
		vals, err := cache.Rbd.BRPop(ctx, 0, scheduler.QUEUE).Result()
		if err != nil {
			log.Println("Error poping from queue ", err.Error())
			time.Sleep(1 * time.Second)
			continue
		}

		if len(vals) < 2 {
			continue
		}

		jobId := vals[1]
		go w.executeJobWithTimeout(jobId)
	}
}

func (w *Worker) executeJobWithTimeout(jobId string) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	done := make(chan error, 1)

	go func() {
		done <- w.executeJob(ctx, jobId)
	}()

	select {
	case err := <-done:
		if err != nil {
			log.Printf("Worker %s: job execution failed for %s: %v", w.ID, jobId, err)
		} else {
			log.Printf("Worker %s: job execution successfull for %s", w.ID, jobId)
		}
	case <-ctx.Done():
		log.Printf("Worker %s: Job %s timed out", w.ID, jobId)
		w.logJobExecution(jobId, string(models.StatusFailed), 0, "Job failed timeout", 0)
		w.updateJob(nil, jobId, models.StatusFailed)
	}
}

// Instead of returing error save the logs.
func (w *Worker) executeJob(ctx context.Context, jobId string) error {
	var job models.Job
	txn := db.DB.Where("id = ?", jobId).First(&job)
	if txn.Error != nil {
		return fmt.Errorf("Error: Job not found %s", txn.Error.Error())
	}
	if job.Retry > MAX_RETRY {
		w.updateJob(&job, job.ID.String(), models.StatusAborted)
		return fmt.Errorf("Error: Job aborted %s due to max retry", jobId)
	}
	switch job.Type {
	case models.JobTypeHTTP:
		w.executeHttpJob(ctx, &job)
	default:
		log.Println("Doesn't support this type right now.")
	}
	return nil
}

type HTTPRequestPayload struct {
	URL     string            `json:"url"`
	Method  string            `json:"method"`
	Headers map[string]string `json:"headers,omitempty"`
	Body    string            `json:"body,omitempty"`
}

func (w *Worker) executeHttpJob(ctx context.Context, job *models.Job) error {
	start := time.Now()

	var payload HTTPRequestPayload
	if err := json.Unmarshal(job.Payload, &payload); err != nil {
		w.logJobExecution(job.ID.String(), string(models.StatusFailed), 0, err.Error(), int64(time.Since(start).Milliseconds()))
		w.updateJob(job, job.ID.String(), models.StatusFailed)
		return err
	}
	if payload.URL == "" {
		return fmt.Errorf("Url is required")
	}
	if payload.Method == "" {
		payload.Method = "GET"
	}

	req, err := http.NewRequestWithContext(ctx, payload.Method, payload.URL, strings.NewReader(payload.Body))
	if err != nil {
		w.logJobExecution(job.ID.String(), string(models.StatusFailed), 0, err.Error(), int64(time.Since(start).Milliseconds()))
		w.updateJob(job, job.ID.String(), models.StatusFailed)
		return err
	}

	for k, v := range payload.Headers {
		req.Header.Set(k, v)
	}
	resp, err := w.client.Do(req)
	duration := time.Since(start).Milliseconds()

	if err != nil {
		w.logJobExecution(job.ID.String(), string(models.StatusFailed), 0, err.Error(), duration)
		w.updateJob(job, job.ID.String(), models.StatusFailed)
		return err
	}
	defer resp.Body.Close()
	bodyBytes, _ := io.ReadAll(io.LimitReader(resp.Body, 10*1024))
	w.logJobExecution(job.ID.String(), resp.Status, resp.StatusCode, string(bodyBytes), duration)
	w.updateJob(job, job.ID.String(), models.StatusPending) // Pending means ready to run again.
	return nil
}

func (w *Worker) logJobExecution(jobId string, status string, statusCode int, Response string, duration int64) {
	jobUuid, err := uuid.Parse(jobId)
	if err != nil {
		log.Println("Invalid jobId")
		return
	}
	logEntry := models.Logs{
		Status:     status,
		StatusCode: statusCode,
		Response:   Response,
		JobID:      jobUuid,
		Duration:   duration,
	}
	if err := db.DB.Create(&logEntry).Error; err != nil {
		log.Println("Error: creating log for jobId", jobId)
		return
	}
}

func (w *Worker) updateJob(job *models.Job, jobId string, status models.StatusType) {
	var updatedJob models.Job
	if job == nil {
		if err := db.DB.Where("id = ?", jobId).First(&updatedJob).Error; err != nil {
			log.Printf("Error: job not found %s", jobId)
			return
		}
	} else {
		updatedJob = *job
	}

	now := time.Now().UTC()
	updatedJob.Status = status
	updatedJob.LastRun = &now

	if status == models.StatusPending {
		updatedJob.NextRun, _ = scheduler.GetNextRun(job.Schedule, job.Timezone)
		updatedJob.Retry = 0
	}
	if status == models.StatusFailed {
		updatedJob.Retry = updatedJob.Retry + 1
		if updatedJob.Retry < MAX_RETRY {
			var retryAt time.Time
			if updatedJob.Retry == 1 {
				retryAt = time.Now().UTC().Add(1 * time.Minute)
			} else {
				delay := time.Duration(1<<updatedJob.Retry) * time.Minute
				retryAt = time.Now().UTC().Add(delay)
			}
			updatedJob.NextRun = &retryAt
		} else {
			updatedJob.Status = models.StatusAborted
		}
	}

	if err := db.DB.Save(updatedJob).Error; err != nil {
		log.Printf("Error Updating the job %s: %v", job.ID, err)
		return
	}
}
