package main

import (
	"log"
	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/cors"
	"crud/api"
	"crud/config"
)

func main() {
	// Connect to database
	db := config.ConnectDB()
	defer db.Close()

	// Initialize Gin router
	r := gin.Default()

	// Enable CORS
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"}, // Allow requests from your frontend
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// Register API routes
	api.RegisterRoutes(r, db)
	api.AuthRoutes(r, db)
	
	log.Println("🚀 Server running on port 9090")
	r.Run(":9090")
}


