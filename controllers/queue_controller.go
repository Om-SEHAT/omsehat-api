package controllers

import (
	"github.com/Om-SEHAT/omsehat-api/services"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func GetCurrentQueue(c *gin.Context) {
	// Get the current queue for the user
	doctorID := c.Param("doctor_id")

	// parse the doctorID to UUID
	doctorUUID, err := uuid.Parse(doctorID)

	if err != nil {
		c.JSON(400, gin.H{"message": "Invalid doctor ID"})
		return
	}

	// get the current queue
	queue, err := services.GetCurrentQueue(doctorUUID)
	if err != nil {
		c.JSON(500, gin.H{"message": err.Error()})
		return
	}

	// check if queue is nil
	if queue == nil {
		c.JSON(404, gin.H{"message": "No queue found"})
		return
	}

	// return the queue
	c.JSON(200, gin.H{
		"queue": queue,
	})
}
