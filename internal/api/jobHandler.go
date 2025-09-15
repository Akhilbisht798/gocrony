package api

import (
	"github.com/akhilbisht798/gocrony/internal/db"
	"github.com/akhilbisht798/gocrony/internal/models"
	"github.com/akhilbisht798/gocrony/internal/scheduler"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
}

// TODO: add auth middleware and authentication later.
// TODO: check the method for url before putting it on db.
// TODO: check the cron schedule before putting it on db.
func CreateNewJob(c *gin.Context) {
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
		Name:     req.Name,
		Schedule: req.Schedule,
		Type:     req.Type,
		Payload:  req.Payload,
	}

	job.NextRun, _ = scheduler.GetNextRun(req.Schedule)

	if err := db.DB.Create(&job).Error; err != nil {
		c.JSON(500, gin.H{
			"error": "failed to create job: " + err.Error(),
		})
	}

	c.JSON(200, gin.H{
		"message": "Job recived",
		"id":      job.ID,
	})
}

func UpdateJob(c *gin.Context) {
	var req models.UpdateJobRequest
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

	updates := make(map[string]any)
	if req.Schedule != "" {
		updates["Schedule"] = req.Schedule
		updates["NextRun"], _ = scheduler.GetNextRun(req.Schedule)
	}
	if req.Name != "" {
		updates["Name"] = req.Name
	}
	if req.Type != "" {
		updates["Type"] = req.Type
	}
	if req.Payload != nil {
		updates["Payload"] = req.Payload
	}
	if err := db.DB.Model(&models.Job{}).Where("id = ?", req.ID).Updates(updates).Error; err != nil {
		c.JSON(500, gin.H{
			"error": "failed to update job: " + err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"message": "job updated successfully",
		"data":    updates,
	})
}

func GetJob(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(400, gin.H{
			"error": "invalid paramerter",
		})
		return
	}
	var job models.Job
	if err := db.DB.First(&job, "id=?", id).Error; err != nil {
		c.JSON(400, gin.H{
			"error": "job not found",
		})
	}
	c.JSON(200, gin.H{
		"job": job,
	})
}

// TODO: Only that jobs that are created by him. so add user auth.
func GetAllJobs(c *gin.Context) {
	var jobs []models.Job
	if err := db.DB.Find(&jobs).Error; err != nil {
		c.JSON(500, gin.H{
			"error": "failed to fetch jobs: " + err.Error(),
		})
	}
	c.JSON(200, gin.H{
		"jobs": jobs,
	})
}

func DeleteJob(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(400, gin.H{
			"error": "invalid paramerter",
		})
		return
	}
	var job models.Job
	if err := db.DB.Delete(&job, "id = ?", id).Error; err != nil {
		c.JSON(500, gin.H{
			"error": "failed to delete job: " + err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"message": "successfully deleted",
		"id":      id,
	})
}

func RunJob(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(400, gin.H{
			"error": "invalid paramerter",
		})
		return
	}
	// Function to run the job would go here.
	c.JSON(200, gin.H{
		"message": "successfully deleted",
	})
}
