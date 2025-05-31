package controllers

import (
	"sort"
	"time"

	"github.com/Om-SEHAT/omsehat-api/schemas"
	"github.com/Om-SEHAT/omsehat-api/services"
	"github.com/Om-SEHAT/omsehat-api/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

var validate = validator.New()

func RegisterUser(c *gin.Context) {
	var input schemas.RegisterUserInput

	// bind and validate the request body to the input struct
	if valid, _ := utils.BindAndValidate(c, &input); !valid {
		return // the response has already been sent in the utility function
	}

	// Call the service to register the user
	user, err := services.RegisterUser(input)

	if err != nil {
		c.JSON(500, gin.H{"message": err.Error()})
		return
	}

	// Send user id, name, and email in the response to prompt OTP verification
	c.JSON(200, gin.H{
		"message": "User registered successfully",
		"user": gin.H{
			"id":    user.ID,
			"name":  user.Name,
			"email": user.Email,
		},
	})
}

func VerifyOTP(c *gin.Context) {
	var input schemas.OTPInput

	// bind and validate the request body to the input struct
	if valid, _ := utils.BindAndValidate(c, &input); !valid {
		return // The response has already been sent in the utility function
	}

	// Call the service to verify the OTP
	session, err := services.ValidateOTP(input)

	if err != nil {
		c.JSON(500, gin.H{"message": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "OTP verified successfully", "session": session})
}

func GetUserDetails(c *gin.Context) {
	userID := c.Param("id")

	// Parse userID to UUID
	id, err := uuid.Parse(userID)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid user ID"})
		return
	}

	// Fetch user details
	user := services.GetUserByID(id)
	if user == nil {
		c.JSON(404, gin.H{"error": "User not found"})
		return
	}

	// Fetch sessions for the user
	sessions := services.GetSessionsByUserID(id)
	if sessions == nil {
		c.JSON(404, gin.H{"error": "No sessions found for the user"})
		return
	}

	// Separate current session and history sessions
	var currentSession gin.H
	var historySessions []gin.H
	latestTimestamp := time.Time{} // Initialize to zero value

	for _, session := range sessions {
		if session.CreatedAt.After(latestTimestamp) {
			if currentSession != nil {
				historySessions = append(historySessions, gin.H{
					"session_id":       currentSession["session_id"],
					"bodytemp":         currentSession["bodytemp"],
					"doctor_diagnosis": currentSession["doctor_diagnosis"],
					"heartrate":        currentSession["heartrate"],
					"height":           currentSession["height"],
					"prediagnosis":     currentSession["prediagnosis"],
					"weight":           currentSession["weight"],
					"created_at":       currentSession["created_at"],
				})
			}
			latestTimestamp = session.CreatedAt

			// Fetch queue for the current session
			queue := services.GetQueueBySessionID(session.ID)

			currentSession = gin.H{
				"queue":            queue,
				"session_id":       session.ID,
				"bodytemp":         session.Bodytemp,
				"doctor_diagnosis": session.DoctorDiagnosis,
				"heartrate":        session.Heartrate,
				"height":           session.Height,
				"prediagnosis":     session.Prediagnosis,
				"weight":           session.Weight,
				"created_at":       session.CreatedAt, // Ensure created_at is set here
			}
		} else {
			historySessions = append(historySessions, gin.H{
				"session_id":       session.ID,
				"bodytemp":         session.Bodytemp,
				"doctor_diagnosis": session.DoctorDiagnosis,
				"heartrate":        session.Heartrate,
				"height":           session.Height,
				"prediagnosis":     session.Prediagnosis,
				"weight":           session.Weight,
				"created_at":       session.CreatedAt, // Use session.CreatedAt directly
			})
		}
	}

	// Sort history sessions by created_at in descending order
	sort.Slice(historySessions, func(i, j int) bool {
		return historySessions[i]["created_at"].(time.Time).After(historySessions[j]["created_at"].(time.Time))
	})

	// Respond with the new structure
	c.JSON(200, gin.H{
		"user": gin.H{
			"id":          user.ID,
			"name":        user.Name,
			"email":       user.Email,
			"gender":      user.Gender,
			"nationality": user.Nationality,
			"age":         utils.DateToAgeString(user.DOB),
		},
		"current_session":  currentSession,
		"history_sessions": historySessions,
	})
}
