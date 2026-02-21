package main

import (
	"cookeasy/apperror"
	"cookeasy/auth"
	"cookeasy/domain"
	"cookeasy/handlers"
	"database/sql"
	"log/slog"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

func main() {
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"}, // frontend URL
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type"},
		AllowCredentials: true, // VERY IMPORTANT for cookies
		MaxAge:           12 * time.Hour,
	}))

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		apperror.RegisterTags(v)
	}

	db, err := sql.Open("postgres", "postgres://postgres:secret@localhost:5432/cookeasy?sslmode=disable")
	if err != nil {
		slog.Error("Failed to connect to database", "error", err)
	}

	domainService := domain.NewDBService(db)

	handler := handlers.NewHandler(domainService)

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	r.POST("/signup", handler.CreateUser)

	r.POST("/login", handler.Login)
	r.POST("/logout", handler.Logout)

	r.GET("/profile", auth.AuthMiddleware, handler.GetUser)

	r.Run(":3001")
}
