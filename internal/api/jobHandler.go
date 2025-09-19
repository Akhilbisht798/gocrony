package api

import (
	"github.com/akhilbisht798/gocrony/internal/db"
	"github.com/akhilbisht798/gocrony/internal/middleware"
	"github.com/akhilbisht798/gocrony/internal/models"
	"github.com/akhilbisht798/gocrony/internal/scheduler"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
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
		Name:     req.Name,
		Schedule: req.Schedule,
		Type:     req.Type,
		Payload:  req.Payload,
		UserID:   userId,
	}

	job.NextRun, _ = scheduler.GetNextRun(req.Schedule)

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
			"error": "error parsing userId: " + err.Error(),
		})
		return
	}
	var req models.UpdateJobRequest
	id := c.Param("id")
	if id == "" {
		c.JSON(400, gin.H{
			"error": "invalid paramerter",
		})
		return
	}
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
	tx := db.DB.Model(&models.Job{}).
		Where("id = ? AND user_id = ?", id, userId).
		Updates(updates)
	if tx.Error != nil {
		c.JSON(500, gin.H{
			"error": "failed to update job: " + tx.Error.Error(),
		})
		return
	}
	if tx.RowsAffected == 0 {
		c.JSON(404, gin.H{
			"error": "job not found or you don't have permission to update this job",
		})
		return
	}

	c.JSON(200, gin.H{
		"message": "job updated successfully",
		"data":    updates,
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
		c.JSON(404, gin.H{ "error": "job not found or not owned by user." })
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
