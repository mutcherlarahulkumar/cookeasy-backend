package handlers

import (
	"cookeasy/apperror"
	"cookeasy/domain"
	"cookeasy/models"
	"errors"
	"log/slog"
	"net/http"

	"golang.org/x/crypto/bcrypt"

	"github.com/gin-gonic/gin"
)

func (h *Handler) CreateUser(c *gin.Context) {

	ctx := c.Request.Context()

	var newUser models.SignUp

	if err := c.ShouldBindJSON(&newUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": apperror.CustomValidationError(err),
		})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(newUser.Password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to process password",
		})
		return
	}

	newUser.Password = string(hashedPassword)

	err = h.service.CreateUser(ctx, newUser)
	if err != nil {

		if errors.Is(err, domain.ErrEmailAlreadyExists) {
			c.JSON(http.StatusConflict, gin.H{
				"error": "email already exists",
			})
			return
		}
		slog.Error("Failed to create user", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to create user",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Signup Successful",
	})
}
