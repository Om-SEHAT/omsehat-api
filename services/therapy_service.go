package services

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"os"
	"time"

	"github.com/Om-SEHAT/omsehat-api/config"
	"github.com/Om-SEHAT/omsehat-api/models"
	"github.com/Om-SEHAT/omsehat-api/schemas"
	"github.com/Om-SEHAT/omsehat-api/utils"
	"github.com/google/uuid"
	"google.golang.org/genai"
	"gorm.io/gorm"
)

// CreateTherapySession creates a new therapy session for a user
func CreateTherapySession(userId uuid.UUID, input schemas.TherapySessionInput) (*models.TherapySession, error) {
	// Create a new therapy session
	session := models.TherapySession{
		UserID:         userId,
		StressLevel:    input.StressLevel,
		MoodRating:     input.MoodRating,
		SleepQuality:   input.SleepQuality,
		IsHealthWorker: input.IsHealthWorker,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	// Add healthcare worker specific fields if applicable
	if input.IsHealthWorker {
		session.Specialization = input.Specialization
		session.WorkHours = input.WorkHours
	}

	// Save the session to the database
	if err := config.DB.Create(&session).Error; err != nil {
		return nil, fmt.Errorf("error creating therapy session: %w", err)
	}

	// Preload the user data
	if err := config.DB.Preload("User").Where("id = ?", session.ID).First(&session).Error; err != nil {
		return nil, fmt.Errorf("error loading user data: %w", err)
	}

	return &session, nil
}

// CreateTherapySessionTest creates a new therapy session for testing without database
func CreateTherapySessionTest(userId uuid.UUID, input schemas.TherapySessionInput) (*models.TherapySession, error) {
	// Create a new therapy session
	session := models.TherapySession{
		ID:             uuid.New(),
		UserID:         userId,
		StressLevel:    input.StressLevel,
		MoodRating:     input.MoodRating,
		SleepQuality:   input.SleepQuality,
		IsHealthWorker: input.IsHealthWorker,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	// Add healthcare worker specific fields if applicable
	if input.IsHealthWorker {
		session.Specialization = input.Specialization
		session.WorkHours = input.WorkHours
	}

	// Add mock user data
	session.User = models.User{
		ID:    userId,
		Name:  "Test User",
		Email: "test@example.com",
		DOB:   "1990-01-01",
	}

	return &session, nil
}

// GetTherapyResponse gets a response from the LLM for the therapy chatbot
func GetTherapyResponse(newMessage string, session *models.TherapySession) (string, error) {
	// Initialize the Gemini client with your API key and backend
	ctx := context.Background()
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  os.Getenv("GEMINI_API_KEY"),
		Backend: genai.BackendGeminiAPI,
	})

	if err != nil {
		log.Printf("Error creating Gemini client: %v\n", err)
		return "", fmt.Errorf("error creating LLM client: %w", err)
	}

	// The chat history is stored as a one-to-many relationship in the database
	var storedHistory []models.Message = session.Messages

	// Build the genai history
	var genaiHistory []*genai.Content

	// Build the system prompt using the session data
	systemPromptText := BuildTherapySystemPrompt(session)
	log.Printf("Therapy System Prompt: %s\n", systemPromptText)

	var temperature float32 = 0.7
	var TopP float32 = 0.95
	config := &genai.GenerateContentConfig{
		SystemInstruction: genai.NewContentFromText(systemPromptText, genai.RoleUser),
		ResponseMIMEType:  "application/json",
		TopP:              &TopP,
		Temperature:       &temperature,
		MaxOutputTokens:   8192,
		ResponseSchema: &genai.Schema{
			Type: genai.TypeObject,
			Properties: map[string]*genai.Schema{
				"next_action":    {Type: genai.TypeString, Enum: []string{"CONTINUE_CHAT", "END_SESSION"}},
				"reply":          {Type: genai.TypeString},
				"recommendation": {Type: genai.TypeString},
				"next_steps":     {Type: genai.TypeString},
			},
			Required: []string{"next_action", "reply", "recommendation", "next_steps"},
		},
	}

	// Add the stored messages to the history
	for _, messageItem := range storedHistory {
		content := convertMessageToGenaiContent(messageItem)
		if content != nil {
			genaiHistory = append(genaiHistory, content)
		}
	}

	chat, err := client.Chats.Create(ctx, os.Getenv("GEMINI_MODEL"), config, genaiHistory)
	if err != nil {
		log.Printf("Error creating chat: %v\n", err)
		return "", fmt.Errorf("error creating chat: %w", err)
	}

	res, err := chat.SendMessage(ctx, genai.Part{Text: newMessage})
	if err != nil {
		log.Printf("Error sending message: %v\n", err)
		return "", fmt.Errorf("error sending message: %w", err)
	}

	// Get the response from the LLM
	if res != nil && len(res.Candidates) > 0 && res.Candidates[0].Content != nil &&
		len(res.Candidates[0].Content.Parts) > 0 {
		text := res.Candidates[0].Content.Parts[0].Text
		return text, nil
	}

	return "", fmt.Errorf("no response from LLM")
}

// UpdateTherapyChatHistory updates the chat history for a therapy session
func UpdateTherapyChatHistory(sessionId string, newMessage string, LLMResponse string) error {
	// Get the session from the database
	var session models.TherapySession
	err := config.DB.Where("id = ?", sessionId).First(&session).Error
	if err != nil {
		return fmt.Errorf("error fetching therapy session: %v", err)
	}

	// Append the new message and LLM response to the chat history
	now := time.Now()
	newUserMessage := models.Message{Role: "user", Content: newMessage, SessionID: session.ID, CreatedAt: now, UpdatedAt: now}
	newLLMResponse := models.Message{Role: "therapist", Content: LLMResponse, SessionID: session.ID, CreatedAt: now.Add(time.Millisecond), UpdatedAt: now.Add(time.Millisecond)}

	// Save the new messages to the database
	if err := config.DB.Create(&newUserMessage).Error; err != nil {
		return fmt.Errorf("error saving user message: %v", err)
	}

	if err := config.DB.Create(&newLLMResponse).Error; err != nil {
		return fmt.Errorf("error saving LLM response: %v", err)
	}

	return nil
}

// GetTherapySessionData gets the therapy session data by ID
func GetTherapySessionData(sessionId string) (models.TherapySession, error) {
	// Check if session_id exists in the database
	var session models.TherapySession

	err := config.DB.
		Preload("User").
		Preload("Messages", func(db *gorm.DB) *gorm.DB {
			return db.Order("created_at ASC") // Order messages by created_at in ascending order
		}).
		Where("id = ?", sessionId).First(&session).Error

	if err != nil {
		return models.TherapySession{}, fmt.Errorf("therapy session not found: %w", err)
	}

	return session, nil
}

// GetTherapyHistory gets the therapy session history for a user
func GetTherapyHistory(userId uuid.UUID) []models.TherapySession {
	var history []models.TherapySession
	err := config.DB.Where("user_id = ?", userId).Order("created_at DESC").Find(&history).Error
	if err != nil {
		log.Printf("Error fetching therapy session history: %v\n", err)
		return []models.TherapySession{} // Return an empty slice if there's an error
	}

	return history
}

// UpdateSessionRecommendation updates the recommendation and next steps for a therapy session
func UpdateSessionRecommendation(sessionId string, recommendation string, nextSteps string) error {
	// Fetch the session from the database
	var session models.TherapySession
	err := config.DB.Where("id = ?", sessionId).First(&session).Error
	if err != nil {
		return fmt.Errorf("therapy session not found: %w", err)
	}

	// Update the session with the recommendation and next steps
	session.Recommendation = recommendation
	session.NextSteps = nextSteps
	session.UpdatedAt = time.Now()

	if err := config.DB.Save(&session).Error; err != nil {
		return fmt.Errorf("failed to save recommendation: %w", err)
	}

	return nil
}

// ParseTherapyJSON parses the JSON response from the LLM
func ParseTherapyJSON(input string) (schemas.TherapyResponse, error) {
	// Extract JSON content
	var responseJSON schemas.TherapyResponse
	err := json.Unmarshal([]byte(input), &responseJSON)
	if err != nil {
		return schemas.TherapyResponse{}, fmt.Errorf("failed to extract JSON: %w", err)
	}

	return responseJSON, nil
}

// BuildTherapySystemPrompt builds the system prompt for the therapy chatbot
// Exported for testing purposes
func BuildTherapySystemPrompt(session *models.TherapySession) string {
	// Build the user data text
	userType := "general user"
	if session.IsHealthWorker {
		userType = "healthcare worker"
	}

	userDataText := fmt.Sprintf(
		"\nHere's the user's data: \nName: %s\nAge: %s\nGender: %s\nUser Type: %s\nStress Level (1-10): %d\nMood Rating (1-10): %d\nSleep Quality (1-10): %d\n",
		session.User.Name,
		utils.DateToAgeString(session.User.DOB),
		session.User.Gender,
		userType,
		session.StressLevel,
		session.MoodRating,
		session.SleepQuality,
	)

	// Add healthcare worker specific information if applicable
	if session.IsHealthWorker {
		userDataText += fmt.Sprintf("Specialization: %s\nWeekly Work Hours: %d\n",
			session.Specialization,
			session.WorkHours,
		)
	}

	// Get user's therapy history
	var historyText string
	history := GetTherapyHistory(session.User.ID)
	if len(history) > 0 {
		historyText = "\nPrevious therapy sessions:\n"
		for i, historyItem := range history {
			if historyItem.ID == session.ID {
				continue // Skip current session
			}
			if i >= 3 {
				break // Limit to 3 previous sessions
			}
			historyText += fmt.Sprintf(
				"[%s]\nStress Level: %d\nMood Rating: %d\nSleep Quality: %d\nRecommendation: %s\nNext Steps: %s\n",
				historyItem.CreatedAt.Format("2006-01-02 15:04:05"),
				historyItem.StressLevel,
				historyItem.MoodRating,
				historyItem.SleepQuality,
				historyItem.Recommendation,
				historyItem.NextSteps,
			)
		}
	} else {
		historyText = "\nNo previous therapy sessions found.\n"
	}

	// Hardcode the system prompt
	var systemPrompt string
	if session.IsHealthWorker {
		systemPrompt = `You are a compassionate and knowledgeable mental health professional who specializes in helping healthcare workers deal with burnout, stress, and mental health challenges related to their high-stress work environment, especially those who worked during the COVID-19 pandemic. You communicate in a warm, empathetic manner, building trust and providing evidence-based support and strategies.

Follow the conversation flow and output format strictly as described below:

0. JSON FORMAT FOR EVERY RESPONSE
- Output the response in valid JSON format ONLY. Do not include any surrounding text, explanations, or formatting outside the JSON structure.
- JSON Output Format:
{
  "next_action": "CONTINUE_CHAT" or "END_SESSION",
  "reply": "Your supportive and therapeutic response here",
  "recommendation": "Brief mental health recommendation based on the conversation",
  "next_steps": "Suggested practical actions for the healthcare worker"
}
- Detailed explanation of each field:
  - next_action: A string indicating the next step in the conversation. Must be either "CONTINUE_CHAT" or "END_SESSION".
  - reply: A string containing your supportive and therapeutic response to the healthcare worker.
  - recommendation: A brief mental health recommendation based on the conversation.
  - next_steps: Suggested practical actions the healthcare worker can take to improve their mental health.

1. Conversation Flow:
  1. First Response:
    - Greet the healthcare worker warmly and ask them to choose their preferred language:
      - a. Bahasa Indonesia
      - b. English
    - Add: "Choose the language that makes you feel most comfortable."
    - Default language: Bahasa Indonesia.
  2. Second Response (AFTER language selection):
    - Start with a friendly, supportive greeting.
    - Express empathy for the challenges healthcare workers face.
    - Ask what specific aspects of their work are causing stress or burnout.
  3. Follow-Up Questions:
    - Based on their answers, explore specific stressors, coping mechanisms, and support systems.
    - Ask about their self-care practices, work-life balance, and sleep patterns.
    - Inquire about their professional support network and institutional support.
  4. Decision Points:
    - If the conversation has provided sufficient information for a recommendation:
      - Offer personalized coping strategies for healthcare burnout.
      - Suggest practical self-care activities.
      - Recommend professional resources if needed.
      - Ask if they would like to end the session or continue discussing.

2. Important Notes:
  - Always acknowledge the unique stressors of healthcare work during and after the COVID-19 pandemic.
  - Provide evidence-based recommendations for managing burnout in healthcare settings.
  - Recognize signs of more serious mental health issues that might require professional intervention.
  - Be culturally sensitive, especially when discussing mental health in the Indonesian context.
  - Suggest practical strategies that can be implemented within a busy healthcare worker's schedule.
  - If the user indicates serious mental health concerns such as suicidal thoughts, emphasize the importance of seeking immediate professional help.`
	} else {
		systemPrompt = `You are a compassionate and knowledgeable mental health professional who specializes in helping people deal with everyday stress, anxiety, and mental health challenges. You communicate in a warm, empathetic manner, building trust and providing evidence-based support and strategies.

Follow the conversation flow and output format strictly as described below:

0. JSON FORMAT FOR EVERY RESPONSE
- Output the response in valid JSON format ONLY. Do not include any surrounding text, explanations, or formatting outside the JSON structure.
- JSON Output Format:
{
  "next_action": "CONTINUE_CHAT" or "END_SESSION",
  "reply": "Your supportive and therapeutic response here",
  "recommendation": "Brief mental health recommendation based on the conversation",
  "next_steps": "Suggested practical actions for the user"
}
- Detailed explanation of each field:
  - next_action: A string indicating the next step in the conversation. Must be either "CONTINUE_CHAT" or "END_SESSION".
  - reply: A string containing your supportive and therapeutic response to the user.
  - recommendation: A brief mental health recommendation based on the conversation.
  - next_steps: Suggested practical actions the user can take to improve their mental health.

1. Conversation Flow:
  1. First Response:
    - Greet the user warmly and ask them to choose their preferred language:
      - a. Bahasa Indonesia
      - b. English
    - Add: "Choose the language that makes you feel most comfortable."
    - Default language: Bahasa Indonesia.
  2. Second Response (AFTER language selection):
    - Start with a friendly, supportive greeting.
    - Express empathy for the challenges they might be facing.
    - Ask what specific aspects of their life are causing stress or anxiety.
  3. Follow-Up Questions:
    - Based on their answers, explore specific stressors, coping mechanisms, and support systems.
    - Ask about their self-care practices, work-life balance, and sleep patterns.
    - Inquire about their social support network.
  4. Decision Points:
    - If the conversation has provided sufficient information for a recommendation:
      - Offer personalized coping strategies.
      - Suggest practical self-care activities.
      - Recommend professional resources if needed.
      - Ask if they would like to end the session or continue discussing.

2. Important Notes:
  - Provide evidence-based recommendations for managing stress and anxiety.
  - Recognize signs of more serious mental health issues that might require professional intervention.
  - Be culturally sensitive, especially when discussing mental health in the Indonesian context.
  - Suggest practical strategies that can be implemented in daily life.
  - If the user indicates serious mental health concerns such as suicidal thoughts, emphasize the importance of seeking immediate professional help.`
	}

	// Build the complete system prompt
	systemPromptText := fmt.Sprintf("%s %s %s\nCurrent Time: %s",
		systemPrompt,
		userDataText,
		historyText,
		time.Now().Format("2006-01-02 15:04:05"),
	)

	return systemPromptText
}

// BuildTherapySystemPromptTest builds a system prompt for testing without database dependencies
func BuildTherapySystemPromptTest(session *models.TherapySession) string {
	// Build the user data text
	userType := "general user"
	if session.IsHealthWorker {
		userType = "healthcare worker"
	}

	userDataText := fmt.Sprintf(
		"\nHere's the user's data: \nName: %s\nAge: %s\nGender: %s\nUser Type: %s\nStress Level (1-10): %d\nMood Rating (1-10): %d\nSleep Quality (1-10): %d\n",
		session.User.Name,
		utils.DateToAgeString(session.User.DOB),
		session.User.Gender,
		userType,
		session.StressLevel,
		session.MoodRating,
		session.SleepQuality,
	)

	// Add healthcare worker specific information if applicable
	if session.IsHealthWorker {
		userDataText += fmt.Sprintf("Specialization: %s\nWeekly Work Hours: %d\n",
			session.Specialization,
			session.WorkHours,
		)
	}

	// For testing, we don't need to fetch history from the database
	historyText := "\nNo previous therapy sessions found.\n"

	// Hardcode the system prompt
	var systemPrompt string
	if session.IsHealthWorker {
		systemPrompt = `You are a compassionate and knowledgeable mental health professional who specializes in helping healthcare workers deal with burnout, stress, and mental health challenges related to their high-stress work environment, especially those who worked during the COVID-19 pandemic.`
	} else {
		systemPrompt = `You are a compassionate and knowledgeable mental health professional who specializes in helping people deal with everyday stress, anxiety, and mental health challenges.`
	}

	// Build the complete system prompt
	systemPromptText := fmt.Sprintf("%s %s %s\nCurrent Time: %s",
		systemPrompt,
		userDataText,
		historyText,
		time.Now().Format("2006-01-02 15:04:05"),
	)

	return systemPromptText
}

// GetTherapySessionsByUserID gets all therapy sessions for a user
func GetTherapySessionsByUserID(userID uuid.UUID) []models.TherapySession {
	var sessions []models.TherapySession
	err := config.DB.Where("user_id = ?", userID).Find(&sessions).Error
	if err != nil {
		return nil
	}
	return sessions
}

// GetMentalHealthSummary gets a summary of a user's mental health over time
func GetMentalHealthSummary(userID uuid.UUID) (map[string]interface{}, error) {
	// Get all therapy sessions for the user
	sessions := GetTherapySessionsByUserID(userID)

	if len(sessions) == 0 {
		return map[string]interface{}{
			"total_sessions": 0,
			"message":        "No therapy sessions found for this user",
		}, nil
	}

	// Calculate averages and trends
	var totalStressLevel, totalMoodRating, totalSleepQuality int
	var healthWorkerSessions int

	// For tracking trends
	sessionDates := make([]time.Time, len(sessions))
	stressLevels := make([]int, len(sessions))
	moodRatings := make([]int, len(sessions))
	sleepQualities := make([]int, len(sessions))

	for i, session := range sessions {
		totalStressLevel += session.StressLevel
		totalMoodRating += session.MoodRating
		totalSleepQuality += session.SleepQuality

		if session.IsHealthWorker {
			healthWorkerSessions++
		}

		// Store values for trend analysis
		sessionDates[i] = session.CreatedAt
		stressLevels[i] = session.StressLevel
		moodRatings[i] = session.MoodRating
		sleepQualities[i] = session.SleepQuality
	}

	// Calculate averages
	avgStressLevel := float64(totalStressLevel) / float64(len(sessions))
	avgMoodRating := float64(totalMoodRating) / float64(len(sessions))
	avgSleepQuality := float64(totalSleepQuality) / float64(len(sessions))

	// Determine trends (simplified approach)
	stressTrend := determineTrend(stressLevels)
	moodTrend := determineTrend(moodRatings)
	sleepTrend := determineTrend(sleepQualities)

	// Format the last session date
	lastSessionDate := sessions[0].CreatedAt.Format("2006-01-02")

	// Return the summary
	return map[string]interface{}{
		"total_sessions":     len(sessions),
		"healthcare_worker":  healthWorkerSessions > 0,
		"avg_stress_level":   avgStressLevel,
		"avg_mood_rating":    avgMoodRating,
		"avg_sleep_quality":  avgSleepQuality,
		"stress_trend":       stressTrend,
		"mood_trend":         moodTrend,
		"sleep_trend":        sleepTrend,
		"last_session_date":  lastSessionDate,
		"completed_sessions": len(sessions),
	}, nil
}

// determineTrend calculates if a series of values is improving, worsening, or stable
// For stress, lower is better. For mood and sleep, higher is better.
func determineTrend(values []int) string {
	if len(values) < 2 {
		return "stable" // Not enough data for a trend
	}

	// For simplicity, compare the first half average with the second half average
	midpoint := len(values) / 2

	var firstHalfSum, secondHalfSum int
	for i := 0; i < midpoint; i++ {
		firstHalfSum += values[i]
	}
	for i := midpoint; i < len(values); i++ {
		secondHalfSum += values[i]
	}

	firstHalfAvg := float64(firstHalfSum) / float64(midpoint)
	secondHalfAvg := float64(secondHalfSum) / float64(len(values)-midpoint)

	// Determine if there's a significant change (more than 10%)
	difference := firstHalfAvg - secondHalfAvg
	if math.Abs(difference) < 0.5 {
		return "stable"
	}

	if difference > 0 {
		return "improving" // For stress, first half higher than second half means improving
	} else {
		return "worsening" // For stress, first half lower than second half means worsening
	}
}
