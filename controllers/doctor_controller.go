package controllers

import (
	"github.com/Om-SEHAT/omsehat-api/services"
	"github.com/Om-SEHAT/omsehat-api/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func DoctorDiagnose(c *gin.Context) {
	sessionId := c.Param("id")

	// Parse the diagnosis from the request body
	var input struct {
		Diagnosis string `json:"diagnosis" validate:"required"`
	}

	if valid, _ := utils.BindAndValidate(c, &input); !valid {
		return // The response has already been sent in the utility function
	}

	// Call the service to save the diagnosis
	if err := services.DoctorDiagnose(sessionId, input.Diagnosis); err != nil {
		c.JSON(500, gin.H{"message": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "Diagnosis saved successfully"})
}

func GetAllDoctors(c *gin.Context) {
	// Fetch all doctors from the database
	doctors := services.GetAllDoctors()

	if len(doctors) == 0 {
		c.JSON(404, gin.H{"message": "No doctors found"})
		return
	}

	// Respond with the list of doctors
	c.JSON(200, gin.H{"doctors": doctors})
}

func GetDoctorDetails(c *gin.Context) {
	doctorID := c.Param("id")

	// Parse doctorID to UUID
	id, err := uuid.Parse(doctorID)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid doctor ID"})
		return
	}

	// Fetch doctor details from the database
	doctor := services.GetDoctorByID(doctorID)
	if doctor == nil {
		c.JSON(404, gin.H{"error": "Doctor not found"})
		return
	}

	// Fetch appointment counts and current queue
	totalAppointments := services.GetTotalAppointments(id)
	dailyAppointments := services.GetDailyAppointments(id)

	// If empty queue, just let currentQueue be nil
	currentQueue, err := services.GetCurrentQueue(id)

	// Respond with aggregated data
	c.JSON(200, gin.H{
		"doctor":                     doctor,
		"appointment_count_all_time": totalAppointments,
		"appointment_count_daily":    dailyAppointments,
		"current_queue":              currentQueue, // Current queue ID
	})
}
