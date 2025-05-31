package controllers

import (
	"log"

	"github.com/Om-SEHAT/omsehat-api/config"
	"github.com/Om-SEHAT/omsehat-api/models"
	"github.com/Om-SEHAT/omsehat-api/schemas"
	"github.com/Om-SEHAT/omsehat-api/services"
	"github.com/Om-SEHAT/omsehat-api/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// CreateTherapySession creates a new therapy session
func CreateTherapySession(c *gin.Context) {
	// Get the user ID from the context
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(401, gin.H{"message": "Unauthorized"})
		return
	}

	// Parse user ID as UUID
	uid, err := uuid.Parse(userID.(string))
	if err != nil {
		c.JSON(400, gin.H{"message": "Invalid user ID"})
		return
	}

	// Get the therapy session input
	var input schemas.TherapySessionInput

	// Bind and validate the input
	if valid, _ := utils.BindAndValidate(c, &input); !valid {
		return // The response has already been sent in the utility function
	}

	// Additional validation for healthcare workers
	if input.IsHealthWorker {
		if input.Specialization == "" {
			c.JSON(400, gin.H{"message": "Specialization is required for healthcare workers"})
			return
		}
		if input.WorkHours <= 0 {
			c.JSON(400, gin.H{"message": "Work hours must be greater than 0 for healthcare workers"})
			return
		}
	}

	// Create the therapy session
	session, err := services.CreateTherapySession(uid, input)
	if err != nil {
		c.JSON(500, gin.H{"message": err.Error()})
		return
	}

	// Return the new session
	c.JSON(201, session)
}

// GenerateTherapyResponse generates a response for a therapy session
func GenerateTherapyResponse(c *gin.Context) {
	session_id := c.Param("id")

	// Check if session_id exists in the database
	var existingSession models.TherapySession

	err := config.DB.Preload("User").
		Preload("Messages", func(db *gorm.DB) *gorm.DB {
			return db.Order("created_at ASC") // Order messages by created_at in ascending order
		}).
		Where("id = ?", session_id).First(&existingSession).Error

	if err != nil {
		c.JSON(404, gin.H{"message": "Therapy session not found"})
		return
	}

	// Get the new message from user
	var input schemas.TherapyChatInput

	// Bind and validate the input
	if valid, _ := utils.BindAndValidate(c, &input); !valid {
		return // The response has already been sent in the utility function
	}

	// Get the message reply from LLM
	reply, err := services.GetTherapyResponse(input.NewMessage, &existingSession)
	if err != nil {
		c.JSON(500, gin.H{"message": err.Error()})
		return
	}

	// Parse the LLM response
	var LLMResponse schemas.TherapyResponse
	LLMResponse, err = services.ParseTherapyJSON(reply)
	if err != nil {
		c.JSON(500, gin.H{"message": "Error parsing LLM response: " + err.Error()})
		return
	}

	// From the LLM response determine the next action
	log.Println("Therapy LLM Response Next Action:", LLMResponse.NextAction)
	log.Println("-----------------------------------")
	log.Println("Therapy LLM Response:", LLMResponse.Reply)

	if next_action := LLMResponse.NextAction; next_action == "CONTINUE_CHAT" {
		// Just continue the chat
	} else if next_action == "END_SESSION" {
		// Update the session with recommendation and next steps
		err = services.UpdateSessionRecommendation(session_id, LLMResponse.Recommendation, LLMResponse.NextSteps)
		if err != nil {
			c.JSON(500, gin.H{"message": err.Error()})
			return
		}
	} else {
		c.JSON(500, gin.H{"message": "Invalid next action"})
		return
	}

	// Update the chat history with the new message and LLM response
	err = services.UpdateTherapyChatHistory(session_id, input.NewMessage, LLMResponse.Reply)
	if err != nil {
		c.JSON(500, gin.H{"message": err.Error()})
		return
	}

	// Send the response back to the client
	c.JSON(200, gin.H{
		"message":        "Chat history updated successfully",
		"next_action":    LLMResponse.NextAction,
		"reply":          LLMResponse.Reply,
		"recommendation": LLMResponse.Recommendation,
		"next_steps":     LLMResponse.NextSteps,
		"session_id":     session_id,
	})
}

// GetActiveTherapySession gets an active therapy session
func GetActiveTherapySession(c *gin.Context) {
	session_id := c.Param("id")

	var session models.TherapySession

	session, err := services.GetTherapySessionData(session_id)
	if err != nil {
		c.JSON(404, gin.H{"message": "Therapy session not found"})
		return
	}

	// Make sure session is active (recommendation is not set yet)
	if session.Recommendation != "" {
		c.JSON(400, gin.H{"message": "Therapy session is completed"})
		return
	}

	// Send the response back to the client
	c.JSON(200, session)
}

// GetTherapySessionHistory gets the therapy session history for a user
func GetTherapySessionHistory(c *gin.Context) {
	// Get the user ID from the context
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(401, gin.H{"message": "Unauthorized"})
		return
	}

	// Parse user ID as UUID
	uid, err := uuid.Parse(userID.(string))
	if err != nil {
		c.JSON(400, gin.H{"message": "Invalid user ID"})
		return
	}

	// Get the therapy session history
	sessions := services.GetTherapySessionsByUserID(uid)

	// Send the response back to the client
	c.JSON(200, gin.H{
		"sessions": sessions,
	})
}

// GetTherapySessionDetail gets the details of a therapy session
func GetTherapySessionDetail(c *gin.Context) {
	session_id := c.Param("id")

	var session models.TherapySession

	session, err := services.GetTherapySessionData(session_id)
	if err != nil {
		c.JSON(404, gin.H{"message": "Therapy session not found"})
		return
	}

	// Send the response back to the client
	c.JSON(200, session)
}

// GetUserMentalHealthSummary gets a summary of a user's mental health over time
func GetUserMentalHealthSummary(c *gin.Context) {
	// Get the user ID from the context
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(401, gin.H{"message": "Unauthorized"})
		return
	}

	// Parse user ID as UUID
	uid, err := uuid.Parse(userID.(string))
	if err != nil {
		c.JSON(400, gin.H{"message": "Invalid user ID"})
		return
	}

	// Get the mental health summary
	summary, err := services.GetMentalHealthSummary(uid)
	if err != nil {
		c.JSON(500, gin.H{"message": err.Error()})
		return
	}

	// Send the response back to the client
	c.JSON(200, summary)
}
