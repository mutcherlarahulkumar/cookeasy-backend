package main

import (
	"cookeasy/apperror"
	"cookeasy/auth"
	"cookeasy/domain"
	"cookeasy/handlers"
	"database/sql"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/gin-contrib/cors"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

func main() {
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{
			"http://localhost:3000",
			"https://cookeasy-frontend-22sf1xw9x-rahulkumars-projects-54f77606.vercel.app",
		},
		AllowMethods: []string{
			"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS",
		},
		AllowHeaders: []string{
			"Origin",
			"Content-Type",
			"Authorization",
			"Accept",
		},
		ExposeHeaders: []string{
			"Content-Length",
		},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		apperror.RegisterTags(v)
	}

	dbURL := os.Getenv("DATABASE_URL")

	db, err := sql.Open("postgres", dbURL)
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

	port := os.Getenv("PORT")
	if port == "" {
		port = "3001"
	}

	r.Run(":" + port)
}
