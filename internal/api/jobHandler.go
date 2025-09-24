package api

import (
	"errors"

	"github.com/akhilbisht798/gocrony/internal/db"
	"github.com/akhilbisht798/gocrony/internal/middleware"
	"github.com/akhilbisht798/gocrony/internal/models"
	"github.com/akhilbisht798/gocrony/internal/scheduler"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
}

func CreateNewJob(c *gin.Context) {
	userId, err := middleware.ParseUserID(c)
	if err != nil {
		c.JSON(401, gin.H{
			"error": "error parsing userId: " + err.Error(),
		})
		return
	}

	var req models.CreateJobRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(400, gin.H{
			"error": "invalid JSON",
		})
		return
	}

	if err := validate.Struct(req); err != nil {
		c.JSON(400, gin.H{
			"error": "validation failed: " + err.Error(),
		})
		return
	}

	job := models.Job{
		Name:      req.Name,
		Schedule:  req.Schedule,
		Type:      req.Type,
		Payload:   req.Payload,
		UserID:    userId,
		Recurring: *req.Recurring,
		Enabled:   *req.Enabled,
		Timezone:  req.Timezone,
	}

	nextRun, err := scheduler.GetNextRun(req.Schedule, req.Timezone)
	if err != nil {
		c.JSON(500, gin.H{
			"error": "failed to create job: " + err.Error(),
		})
		return
	}
	job.NextRun = nextRun

	if err := db.DB.Create(&job).Error; err != nil {
		c.JSON(500, gin.H{
			"error": "failed to create job: " + err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"message": "Job recived",
		"id":      job.ID,
	})
}

func UpdateJob(c *gin.Context) {
	userId, err := middleware.ParseUserID(c)
	if err != nil {
		c.JSON(401, gin.H{
			"error": "unauthorized: " + err.Error(),
		})
		return
	}

	// Validate job ID parameter
	jobID := c.Param("id")
	if jobID == "" {
		c.JSON(400, gin.H{
			"error": "job ID is required",
		})
		return
	}

	// Parse and validate request
	var req models.UpdateJobRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(400, gin.H{
			"error": "invalid JSON format",
		})
		return
	}

	if err := validate.Struct(req); err != nil {
		c.JSON(400, gin.H{
			"error": "validation failed: " + err.Error(),
		})
		return
	}

	// Check if job exists and user has permission
	var existingJob models.Job
	if err := db.DB.Where("id = ? AND user_id = ?", jobID, userId).First(&existingJob).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(404, gin.H{
				"error": "job not found or access denied",
			})
		} else {
			c.JSON(500, gin.H{
				"error": "failed to fetch job: " + err.Error(),
			})
		}
		return
	}

	// Build updates map
	updates := make(map[string]any)
	shouldRecalculateNextRun := false

	if req.Name != "" {
		updates["name"] = req.Name
	}

	if req.Type != "" {
		updates["type"] = req.Type
	}

	if req.Payload != nil {
		updates["payload"] = req.Payload
	}

	if req.Recurring != nil {
		updates["recurring"] = *req.Recurring
	}

	if req.Enabled != nil {
		updates["enabled"] = *req.Enabled
	}

	// Handle schedule update
	if req.Schedule != "" {
		// Get timezone (use existing or new)
		timezone := existingJob.Timezone
		if req.Timezone != "" {
			timezone = req.Timezone
		}

		// Validate schedule format with timezone
		nextRun, err := scheduler.GetNextRun(req.Schedule, timezone)
		if err != nil {
			c.JSON(400, gin.H{
				"error": "invalid schedule format: " + err.Error(),
			})
			return
		}

		updates["schedule"] = req.Schedule
		updates["next_run"] = nextRun
		shouldRecalculateNextRun = true

		// Reset job status when schedule changes
		updates["status"] = models.StatusPending
		updates["retry"] = 0
	}

	// Handle timezone update - also triggers next_run recalculation
	if req.Timezone != "" {
		updates["timezone"] = req.Timezone
		shouldRecalculateNextRun = true
	}

	// Recalculate next_run if timezone changed but schedule didn't
	if shouldRecalculateNextRun && req.Schedule == "" {
		// Use existing schedule with new timezone
		nextRun, err := scheduler.GetNextRun(existingJob.Schedule, req.Timezone)
		if err != nil {
			c.JSON(400, gin.H{
				"error": "failed to recalculate next run: " + err.Error(),
			})
			return
		}
		updates["next_run"] = nextRun

		// Reset job status when timezone changes
		updates["status"] = models.StatusPending
		updates["retry"] = 0
	}

	// Perform update
	if err := db.DB.Model(&existingJob).Updates(updates).Error; err != nil {
		c.JSON(500, gin.H{
			"error": "failed to update job: " + err.Error(),
		})
		return
	}

	// Fetch updated job for response
	var updatedJob models.Job
	if err := db.DB.Where("id = ?", jobID).First(&updatedJob).Error; err != nil {
		c.JSON(500, gin.H{
			"error": "failed to fetch updated job: " + err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"message": "job updated successfully",
		"job":     updatedJob,
	})
}

func GetJob(c *gin.Context) {
	userId, err := middleware.ParseUserID(c)
	if err != nil {
		c.JSON(401, gin.H{
			"error": "error parsing userId: " + err.Error(),
		})
		return
	}
	id := c.Param("id")
	if id == "" {
		c.JSON(400, gin.H{
			"error": "invalid paramerter",
		})
		return
	}
	var job models.Job
	if err := db.DB.First(&job, "id = ? AND user_id = ?", id, userId).Error; err != nil {
		c.JSON(404, gin.H{
			"error": "job not found",
		})
		return
	}
	c.JSON(200, gin.H{
		"job": job,
	})
}

// TODO: Only that jobs that are created by him. so add user auth.
func GetAllJobs(c *gin.Context) {
	userId, err := middleware.ParseUserID(c)
	if err != nil {
		c.JSON(401, gin.H{
			"error": "error parsing userId: " + err.Error(),
		})
		return
	}
	var jobs []models.Job
	tx := db.DB.Where("user_id = ?", userId).Find(&jobs)
	if tx.Error != nil {
		c.JSON(500, gin.H{
			"error": "failed to fetch jobs: " + tx.Error.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"jobs": jobs,
	})
}

func DeleteJob(c *gin.Context) {
	userId, err := middleware.ParseUserID(c)
	if err != nil {
		c.JSON(401, gin.H{
			"error": "error parsing userId: " + err.Error(),
		})
		return
	}
	id := c.Param("id")
	if id == "" {
		c.JSON(400, gin.H{
			"error": "invalid paramerter",
		})
		return
	}
	tx := db.DB.Delete(&models.Job{}, "id = ? AND user_id = ?", id, userId)
	if tx.Error != nil {
		c.JSON(500, gin.H{
			"error": "failed to delete job: " + tx.Error.Error(),
		})
		return
	}
	if tx.RowsAffected == 0 {
		c.JSON(404, gin.H{"error": "job not found or not owned by user."})
		return
	}
	c.JSON(200, gin.H{
		"message": "successfully deleted",
		"id":      id,
	})
}

func RunJob(c *gin.Context) {
	_, err := middleware.ParseUserID(c)
	if err != nil {
		c.JSON(401, gin.H{
			"error": "error parsing userId: " + err.Error(),
		})
		return
	}
	id := c.Param("id")
	if id == "" {
		c.JSON(400, gin.H{
			"error": "invalid paramerter",
		})
		return
	}
	// Function to run the job would go here.
	c.JSON(200, gin.H{
		"message": "successfully runned the job",
	})
}
