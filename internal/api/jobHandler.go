package api

import (
	"encoding/json"

	"github.com/akhilbisht798/gocrony/internal/db"
	"github.com/akhilbisht798/gocrony/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"gorm.io/datatypes"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
}

// TODO: add auth middleware and authentication later.
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

	headerJson, _ := json.Marshal(req.Headers)
	bodyJson, _ := json.Marshal(req.Body)

	job := models.Jobs{
		URL:      req.URL,
		Method:   req.Method,
		Headers:  datatypes.JSON(headerJson),
		Body:     datatypes.JSON(bodyJson),
		Schedule: req.Schedule,
	}
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

	updates := make(map[string]interface{})
	if req.URL != "" {
		updates["url"] = req.URL
	}
	if req.Method != "" {
		updates["method"] = req.Method
	}
	if req.Headers != nil {
		headerJson, _ := json.Marshal(req.Headers)
		updates["headers"] = datatypes.JSON(headerJson)
	}
	if req.Body != nil {
		bodyJson, _ := json.Marshal(req.Body)
		updates["body"] = datatypes.JSON(bodyJson)
	}
	if req.Schedule != "" {
		updates["schedule"] = req.Schedule
	}
	if err := db.DB.Model(&models.Jobs{}).Where("id = ?", req.ID).Updates(updates).Error; err != nil {
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
	var job models.Jobs
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
	var jobs []models.Jobs
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
	var job models.Jobs
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
