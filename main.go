package main

import (
	"log"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"github.com/Om-SEHAT/omsehat-api/config"
	"github.com/Om-SEHAT/omsehat-api/controllers"
	"github.com/Om-SEHAT/omsehat-api/middleware"
)

func main() {
	// load .env file
	if _, err := os.Stat(".env"); err == nil {
		err = godotenv.Load()
		if err != nil {
			log.Fatal("Error loading .env file")
		}
	} else if os.IsNotExist(err) {
		log.Println(".env file not found, using environment variables from Docker")
	} else {
		log.Fatal("Error checking .env file:", err)
	}

	// Connect to the database
	config.ConnectDatabase()

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true

	r := gin.Default()
	r.Use(cors.New(corsConfig))

	// register routes
	r.POST("/register", controllers.RegisterUser)
	r.POST("/verify-otp", controllers.VerifyOTP)

	// session routes
	r.GET("/session/:id", controllers.GetActiveSession)
	r.POST("/session/:id", controllers.GenerateSessionResponse)
	r.POST("/session/:id/diagnose", controllers.DoctorDiagnose)

	// therapy session routes
	therapyRoutes := r.Group("/therapy")
	therapyRoutes.Use(middleware.AuthMiddleware())
	{
		therapyRoutes.POST("", controllers.CreateTherapySession)
		therapyRoutes.GET("/:id", controllers.GetActiveTherapySession)
		therapyRoutes.POST("/:id", controllers.GenerateTherapyResponse)
		therapyRoutes.GET("/history", controllers.GetTherapySessionHistory)
		therapyRoutes.GET("/detail/:id", controllers.GetTherapySessionDetail)
		therapyRoutes.GET("/mental-health-summary", controllers.GetUserMentalHealthSummary)
	}

	// queue routes
	r.GET("/queue/:doctor_id", controllers.GetCurrentQueue)

	// doctor routes
	r.GET("/doctors", controllers.GetAllDoctors)
	r.GET("/doctor/:id", controllers.GetDoctorDetails)

	// user routes
	r.GET("/user/:id", controllers.GetUserDetails)

	// test routes
	r.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong - om sehat API is running")
	})

	r.Run(":8080") // listen and serve on
}
