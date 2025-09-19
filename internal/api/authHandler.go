package api

import (
	"context"
	"errors"
	"log"
	"net/http"

	"github.com/akhilbisht798/gocrony/internal/auth"
	"github.com/akhilbisht798/gocrony/internal/db"
	"github.com/akhilbisht798/gocrony/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/markbates/goth/gothic"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func GetAuthCallbackFunction(c *gin.Context) {
	c.Request = c.Request.WithContext(
		context.WithValue(c.Request.Context(), "provider", c.Param("provider")),
	)

	gothicUser, err := gothic.CompleteUserAuth(c.Writer, c.Request)
	if err != nil {
		log.Println("Error completing user auth:", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication failed"})
		return
	}

	var user models.User
	var userIdentity models.UserIdentity
	res := db.DB.Preload("User").Where("provider = ? AND provider_id = ?", gothicUser.Provider, gothicUser.UserID).First(&userIdentity)
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			if err := db.DB.Where("email = ?", gothicUser.Email).First(&user).Error; err == nil {
				newIdentity := models.UserIdentity{
					UserID:     user.ID,
					Provider:   gothicUser.Provider,
					ProviderID: gothicUser.UserID,
				}
				db.DB.Create(&newIdentity)
			} else {
				user = models.User{
					Email:     gothicUser.Email,
					Name:      gothicUser.Name,
					AvatarUrl: gothicUser.AvatarURL,
				}
				db.DB.Create(&user)
				newIdentity := models.UserIdentity{
					UserID:     user.ID,
					Provider:   gothicUser.Provider,
					ProviderID: gothicUser.UserID,
				}
				db.DB.Create(&newIdentity)
			}
		} else {
			log.Println("Internal server error: ", res.Error)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error" + res.Error.Error()})
			return
		}
	}
	if res.Error == nil {
		user = userIdentity.User
	}
	//Send back jwt.
	token, err := auth.GenrateJWT(user.ID.String(), user.Email, gothicUser.Provider)
	if err != nil {
		log.Println("Error generating JWT:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Authentication successful",
		"token":   token,
	})
}

func Logout(c *gin.Context) {
	gothic.Logout(c.Writer, c.Request)
	c.JSON(http.StatusOK, gin.H{"message": "Logged out"})
}

func GetAuthProvider(c *gin.Context) {
	c.Request = c.Request.WithContext(
		context.WithValue(c.Request.Context(), "provider", c.Param("provider")),
	)
	gothic.BeginAuthHandler(c.Writer, c.Request)
}

func EmailPasswordAuthSignUp(c *gin.Context) {
	var req models.UserSignUpRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(400, gin.H{
			"error": "invalid request body: " + err.Error(),
		})
		return
	}

	if err := validate.Struct(req); err != nil {
		c.JSON(400, gin.H{
			"error": "validation failed: " + err.Error(),
		})
		return
	}
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(500, gin.H{
			"error": "failed to hash password: " + err.Error(),
		})
		return
	}
	var user models.User
	err = db.DB.Preload("Identities").First(&user, "email = ?", req.Email).Error
	if err == nil {
		for _, identity := range user.Identities {
			if identity.Provider == "email" {
				c.JSON(400, gin.H{
					"error": "user already exists",
				})
				return
			}
		}

		newIdentity := models.UserIdentity{
			UserID:       user.ID,
			Provider:     string(auth.Email),
			ProviderID:   req.Email,
			PasswordHash: string(hashPassword),
		}
		if err := db.DB.Create(&newIdentity).Error; err != nil {
			c.JSON(500, gin.H{"error": "failed to create identity: " + err.Error()})
			return
		}

		c.JSON(200, gin.H{"message": "Email/password identity added to existing user"})
		return
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		newUser := models.User{
			Email:     req.Email,
			Name:      req.Name,
			AvatarUrl: req.AvatarUrl,
		}
		if err := db.DB.Create(&newUser).Error; err != nil {
			c.JSON(500, gin.H{"error": "failed to create user: " + err.Error()})
		}
		newIdentity := models.UserIdentity{
			UserID:       newUser.ID,
			Provider:     string(auth.Email),
			ProviderID:   req.Email,
			PasswordHash: string(hashPassword),
		}
		if err := db.DB.Create(&newIdentity).Error; err != nil {
			c.JSON(500, gin.H{"error": "failed to create identity: " + err.Error()})
			return
		}

		c.JSON(200, gin.H{"message": "Email/password identity added to existing user"})
		return
	}

	c.JSON(500, gin.H{"error": "database error: " + err.Error()})
}

func EmailPasswordAuthSignIn(c *gin.Context) {
	var req models.UserLoginRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(400, gin.H{
			"error": "invalid request body: " + err.Error(),
		})
		return
	}

	if err := validate.Struct(req); err != nil {
		c.JSON(400, gin.H{
			"error": "validation failed: " + err.Error(),
		})
		return
	}
	var userIdentity models.UserIdentity
	res := db.DB.Preload("User").Where("provider_id = ? AND provider = ?", req.Email, "email").First(&userIdentity)
	if res.Error != nil {
		c.JSON(400, gin.H{
			"error": "user not found, Sign up first",
		})
		return
	}

	err := bcrypt.CompareHashAndPassword([]byte(userIdentity.PasswordHash), []byte(req.Password))
	if err != nil {
		c.JSON(400, gin.H{
			"error": "wrong password.",
		})
		return
	}

	token, err := auth.GenrateJWT(userIdentity.User.ID.String(), userIdentity.User.Email, string(auth.Email))
	if err != nil {
		c.JSON(400, gin.H{
			"error": "error creating token",
		})
		return
	}

	c.JSON(200, gin.H{
		"token":   token,
		"message": "succesfully logged in",
	})
}
