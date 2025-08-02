package api

import (
	"encoding/json"
	"log"

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
	log.Println(req.URL)
	c.JSON(200, gin.H{
		"message": "Job recived",
		"data":    req,
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
	c.JSON(200, gin.H{
		"id": id,
	})
}

// TODO: Only that jobs that are created by him. so add user auth.
func GetAllJobs(c *gin.Context) {
	c.JSON(200, gin.H{
		"jobs": "o",
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
	c.JSON(200, gin.H{
		"message": "successfully deleted",
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
	c.JSON(200, gin.H{
		"message": "successfully deleted",
	})
}
