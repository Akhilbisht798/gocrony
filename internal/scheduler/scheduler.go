package scheduler

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/akhilbisht798/gocrony/internal/cache"
	"github.com/akhilbisht798/gocrony/internal/db"
	"github.com/akhilbisht798/gocrony/internal/models"
	"github.com/robfig/cron/v3"
)

const QUEUE = "jobs"

// TODO: save errors and response as logs.
func Scheduler() {
	log.Println("job scheduler started")
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		log.Println("RUNNING AGAIN")
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		if err := getJobsAndSchedule(ctx); err != nil {
			log.Printf("Error processing schedule jobs: %v", err)
		}
		cancel()
	}
}

func getJobsAndSchedule(ctx context.Context) error {
	var jobs []models.Job
	now := time.Now().UTC()

	err := db.DB.Where("next_run <= ? AND enabled = ? AND (status IS NULL OR status = '' OR status = ? OR status = ?)",
		now, true, models.StatusPending, models.StatusFailed).Find(&jobs).Error
	if err != nil {
		return fmt.Errorf("Error: failed to fetch schedule jobs %w", err)
	}
	log.Printf("Found %d jobs to schedule", len(jobs))
	if len(jobs) == 0 {
		return nil
	}

	for _, job := range jobs {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			if err := processJobs(ctx, &job); err != nil {
				log.Printf("Error Processing job %s: %v", job.ID, err.Error())
				continue
			}
		}
	}
	return nil
}

func processJobs(ctx context.Context, job *models.Job) error {
	// Queue it.
	err := enqueueJob(ctx, job.ID.String())
	if err != nil {
		log.Println("Error pushing to queue: ", err.Error())
		return err
	}
	// Status to be queued.
	updates := map[string]any{
		"status": models.StatusQueued,
	}
	log.Println("Enqueued job and updating it.")
	err = db.DB.Model(job).Updates(updates).Error
	if err != nil {
		return fmt.Errorf("Error: unable to update status of the job %w", err)
	}
	log.Println("updated success.")
	return nil
}

func enqueueJob(ctx context.Context, jobId string) error {
	if cache.Rbd == nil {
		return errors.New("redis client not initialized.")
	}
	if jobId == "" {
		return errors.New("Job Id cannot be empty")
	}
	return cache.Rbd.LPush(ctx, QUEUE, jobId).Err()
}

func GetNextRun(schedule string, tzone string) (*time.Time, error) {
	loc, err := time.LoadLocation(tzone)
	if err != nil {
		return &time.Time{}, err
	}
	scheduler, err := cron.ParseStandard(schedule)
	now := time.Now().In(loc)
	if err != nil {
		return &now, err
	}
	nextRun := scheduler.Next(now).UTC()
	return &nextRun, nil
}
