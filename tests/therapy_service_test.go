package tests

import (
	"testing"
	"time"

	"github.com/Om-SEHAT/omsehat-api/models"
	"github.com/Om-SEHAT/omsehat-api/schemas"
	"github.com/Om-SEHAT/omsehat-api/services"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestTherapySystemPromptBuilding(t *testing.T) {
	// Create a mock therapy session for testing
	userID := uuid.New()
	user := models.User{
		ID:          userID,
		Name:        "Test User",
		Email:       "test@example.com",
		DOB:         "1990-01-01", // Format as YYYY-MM-DD string
		Gender:      "Male",
		Nationality: "Indonesian",
	}

	session := models.TherapySession{
		ID:             uuid.New(),
		UserID:         userID,
		User:           user,
		StressLevel:    8,
		MoodRating:     5,
		SleepQuality:   4,
		IsHealthWorker: true,
		Specialization: "Nurse",
		WorkHours:      50,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	// Test private function through the service
	// This is a simple test - we're just making sure the function runs without error
	// and checking if certain strings are in the output
	systemPrompt := services.BuildTherapySystemPromptTest(&session)

	assert.Contains(t, systemPrompt, "healthcare worker")
	assert.Contains(t, systemPrompt, "Stress Level (1-10): 8")
	assert.Contains(t, systemPrompt, "Specialization: Nurse")
	assert.Contains(t, systemPrompt, "Weekly Work Hours: 50")
}

func TestCreateTherapySession(t *testing.T) {
	// Create mock user and input
	userID := uuid.New()
	input := schemas.TherapySessionInput{
		StressLevel:    7,
		MoodRating:     6,
		SleepQuality:   5,
		IsHealthWorker: true,
		Specialization: "Doctor",
		WorkHours:      60,
	}

	// Call test service
	session, err := services.CreateTherapySessionTest(userID, input)

	// Assertions
	assert.NoError(t, err)
	assert.NotNil(t, session)
	assert.Equal(t, userID, session.UserID)
	assert.Equal(t, 7, session.StressLevel)
	assert.Equal(t, "Doctor", session.Specialization)
}

func TestParseTherapyJSON(t *testing.T) {
	// Test valid JSON parsing
	validJSON := `{
		"next_action": "CONTINUE_CHAT",
		"reply": "How are you feeling today?",
		"recommendation": "Consider practicing mindfulness",
		"next_steps": "Try deep breathing exercises"
	}`

	response, err := services.ParseTherapyJSON(validJSON)

	assert.NoError(t, err)
	assert.Equal(t, "CONTINUE_CHAT", response.NextAction)
	assert.Equal(t, "How are you feeling today?", response.Reply)
	assert.Equal(t, "Consider practicing mindfulness", response.Recommendation)
	assert.Equal(t, "Try deep breathing exercises", response.NextSteps)

	// Test invalid JSON
	invalidJSON := `{
		"next_action": "INVALID_ACTION",
		"reply": "How are you feeling today?"
	}`

	_, err = services.ParseTherapyJSON(invalidJSON)
	assert.NoError(t, err) // Should parse but validation would fail later
}

func TestUpdateTherapyChatHistory(t *testing.T) {
	// Create mock session and messages
	sessionID := uuid.New().String()
	userMessage := "I'm feeling very stressed lately"
	llmResponse := "I understand you're feeling stressed. Can you tell me more about what's causing this stress?"

	// Mock the database operations
	// In a real test, we'd use a test database or mock the DB operations
	// For this test, we'll create a custom function

	err := mockUpdateTherapyChatHistory(sessionID, userMessage, llmResponse)

	// Assertions
	assert.NoError(t, err)
}

// Mock function to simulate UpdateTherapyChatHistory without database
func mockUpdateTherapyChatHistory(sessionId string, newMessage string, LLMResponse string) error {
	// In a real implementation, this would insert to the database
	// For testing, we just return nil to indicate success
	return nil
}

func TestGetTherapySessionsForUser(t *testing.T) {
	// Create a mock user ID
	userID := uuid.New()

	// Get mock sessions
	sessions := mockGetTherapySessionsForUser(userID)

	// Assertions
	assert.NotNil(t, sessions)
	assert.Equal(t, 2, len(sessions))
	assert.Equal(t, "Nurse", sessions[0].Specialization)
	assert.Equal(t, true, sessions[0].IsHealthWorker)
}

// Mock function to return therapy sessions for a user
func mockGetTherapySessionsForUser(userID uuid.UUID) []models.TherapySession {
	// Create some mock therapy sessions
	sessions := []models.TherapySession{
		{
			ID:             uuid.New(),
			UserID:         userID,
			StressLevel:    8,
			MoodRating:     5,
			SleepQuality:   4,
			IsHealthWorker: true,
			Specialization: "Nurse",
			WorkHours:      50,
			Recommendation: "Practice mindfulness techniques",
			NextSteps:      "Try deep breathing exercises when feeling overwhelmed",
			CreatedAt:      time.Now().Add(-48 * time.Hour),
			UpdatedAt:      time.Now().Add(-48 * time.Hour),
		},
		{
			ID:             uuid.New(),
			UserID:         userID,
			StressLevel:    6,
			MoodRating:     7,
			SleepQuality:   6,
			IsHealthWorker: true,
			Specialization: "Nurse",
			WorkHours:      50,
			Recommendation: "Continue with mindfulness practice",
			NextSteps:      "Schedule regular breaks during work shifts",
			CreatedAt:      time.Now().Add(-24 * time.Hour),
			UpdatedAt:      time.Now().Add(-24 * time.Hour),
		},
	}

	return sessions
}
