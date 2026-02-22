package handlers

import (
	"cookeasy/apperror"
	"cookeasy/auth"
	"cookeasy/domain"
	"cookeasy/models"
	"errors"
	"net/http"

	"golang.org/x/crypto/bcrypt"

	"github.com/gin-gonic/gin"
)

func (h *Handler) Login(c *gin.Context) {

	ctx := c.Request.Context()

	var req models.Login
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": apperror.CustomValidationError(err),
		})
		return
	}

	user, err := h.service.GetUserByEmail(ctx, req.Email)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "invalid email or password",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "internal server error",
		})
		return
	}

	err = bcrypt.CompareHashAndPassword(
		[]byte(user.Password),
		[]byte(req.Password),
	)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "invalid email or password",
		})
		return
	}

	token, err := auth.GenerateToken(user.ID.String())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to generate token",
		})
		return
	}

	c.SetSameSite(http.SameSiteNoneMode)

	c.SetCookie(
		"auth_token",
		token,
		3600*24, // 1 day
		"/",
		"",   // change in local to (localhost) and in prod to your domain
		true, // secure (true in prod https)
		true, // httpOnly
	)

	c.JSON(http.StatusOK, gin.H{
		"message": "login successful",
	})
}
