package handlers

import "github.com/gin-gonic/gin"

func (h *Handler) Logout(c *gin.Context) {
	c.SetCookie(
		"auth_token",
		"",
		-1,
		"/",
		"localhost",
		false,
		true,
	)

	c.JSON(200, gin.H{"message": "logged out"})
}
