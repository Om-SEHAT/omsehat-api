package controllers

import (
	"log"
	"os"

	"github.com/Om-SEHAT/omsehat-api/config"
	"github.com/Om-SEHAT/omsehat-api/models"
	"github.com/Om-SEHAT/omsehat-api/schemas"
	"github.com/Om-SEHAT/omsehat-api/services"
	"github.com/Om-SEHAT/omsehat-api/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GenerateSessionResponse(c *gin.Context) {
	session_id := c.Param("id")

	// check if session_id exists in the database
	var existingSession models.Session

	err := config.DB.Preload("User").
		Preload("Messages", func(db *gorm.DB) *gorm.DB {
			return db.Order("created_at ASC") // Order messages by created_at in ascending order (for earlier messages first)
		}).
		Where("id = ?", session_id).First(&existingSession).Error

	if err != nil {
		c.JSON(404, gin.H{"message": "Session not found"})
		return
	}

	// get the new message from user
	var input schemas.SessionChatInput

	// bind and validate the input
	if valid, _ := utils.BindAndValidate(c, &input); !valid {
		return // The response has already been sent in the utility function
	}

	// get the message reply from LLM
	reply, err := services.GetLLMResponse(input.NewMessage, &existingSession)
	if err != nil {
		c.JSON(500, gin.H{"message": err.Error()})
		return
	}

	// parse the LLM response
	var LLMResponse schemas.LLMResponse
	LLMResponse, err = services.ParseJSON(reply)

	// queue var
	var queue *models.Queue = nil
	var currentQueue *models.Queue

	// from the LLM response determine the next action
	log.Println("LLM Response Next Action:", LLMResponse.NextAction)
	log.Println("-----------------------------------")
	log.Println("LLM Response Doctor ID:", LLMResponse.DoctorID)
	log.Println("-----------------------------------")
	log.Println("LLM Response:", LLMResponse.Reply)
	if next_action := LLMResponse.NextAction; next_action == "CONTINUE_CHAT" {
		// just continue

	} else if next_action == "APPOINTMENT" {
		// create queue
		queue, err = services.GenerateQueue(session_id, LLMResponse.DoctorID)
		if err != nil {
			c.JSON(500, gin.H{"message": err.Error()})
			return
		}

		// preload queue's doctor
		err = config.DB.Preload("Doctor").Where("id = ?", queue.ID).First(&queue).Error
		if err != nil {
			c.JSON(500, gin.H{"message": err.Error()})
			return
		}

		// send email to the user
		currentQueue, err = services.GetCurrentQueue(queue.DoctorID)
		if err != nil {
			c.JSON(500, gin.H{"message": err.Error()})
			return
		}

		_, err = services.SendQueueEmail(existingSession.User.Email, queue.Number, currentQueue.Number, os.Getenv("EMAIL_TOKEN"), queue.Doctor)
		if err != nil {
			log.Println("Error sending email:", err)
		}

		// update the session's prediagnosis
		existingSession.Prediagnosis = LLMResponse.PreDiagnosis

		err = config.DB.Save(&existingSession).Error
		if err != nil {
			c.JSON(500, gin.H{"message": err.Error()})
			return
		}

	} else {
		c.JSON(500, gin.H{"message": "Invalid next action"})
		return
	}

	// update the chat history with the new message and LLM response
	err = services.UpdateChatHistory(session_id, input.NewMessage, LLMResponse.Reply)
	if err != nil {
		c.JSON(500, gin.H{"message": err.Error()})
		return
	}

	// send the response back to the client
	c.JSON(200, gin.H{
		"message":       "Chat history updated successfully",
		"next_action":   LLMResponse.NextAction,
		"reply":         LLMResponse.Reply,
		"session_id":    session_id,
		"queue":         queue,        // queue is nil if next_action is not APPOINTMENT
		"current_queue": currentQueue, // currentQueue is nil if next_action is not APPOINTMENT
	})
}

func GetActiveSession(c *gin.Context) {
	session_id := c.Param("id")

	var session models.Session

	session, err := services.GetSessionData(session_id)
	if err != nil {
		c.JSON(404, gin.H{"message": "Session not found"})
		return
	}

	// make sure sessoin is active (prediagnosis is not done yet)
	if session.Prediagnosis != "" {
		c.JSON(400, gin.H{"message": "Session is completed"})
		return
	}

	// send the response back to the client
	c.JSON(200, session)
}
