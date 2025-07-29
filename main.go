package main

import (
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/monoMonu/travel-itinerary-pdf/api"
)

func main() {
	app := gin.Default()

	app.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173", "https://vigovia-assessment.netlify.app"},
		AllowMethods:     []string{"GET", "POST", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	app.Static("/pdfs", "./pdfs")

	app.GET("/", func(reqCtx *gin.Context) {
		reqCtx.JSON(http.StatusOK, "Hello World")
	})

	app.POST("/generate-itinerary", api.GeneratePDF)

	app.Run(":3002")
}
