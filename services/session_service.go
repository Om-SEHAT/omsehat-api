package services

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/Om-SEHAT/omsehat-api/config"
	"github.com/Om-SEHAT/omsehat-api/models"
	"github.com/Om-SEHAT/omsehat-api/schemas"
	"github.com/Om-SEHAT/omsehat-api/utils"
	"github.com/google/uuid"
	"google.golang.org/genai"
	"gorm.io/gorm"
)

func convertMessageToGenaiContent(message models.Message) *genai.Content {
	// Determine the role of the message
	var role genai.Role
	if message.Role == "user" {
		role = genai.RoleUser
	} else if message.Role == "omsapa" {
		role = genai.RoleModel
	} else {
		return nil // Invalid role, return nil or handle error as needed
	}

	// Create a new genai.Content object from the message content and role
	content := genai.NewContentFromText(message.Content, role)
	return content
}

func GetLLMResponse(newMessage string, session *models.Session) (string, error) {
	// initialize the Gemini client with your API key and backend
	ctx := context.Background()
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  os.Getenv("GEMINI_API_KEY"),
		Backend: genai.BackendGeminiAPI,
	})

	if err != nil {
		log.Printf("Error creating Gemini client: %v\n", err)
		return "", fmt.Errorf("error creating LLM client: %w", err)
	}

	// the chat history is stored as a one-to-many relationship in the database
	var storedHistory []models.Message = session.Messages

	// build the genai history
	var genaiHistory []*genai.Content

	// build the system prompt using the session data
	systemPromptText := buildSystemPrompt(session)
	log.Printf("System Prompt: %s\n", systemPromptText)

	var temperature float32 = 0.8
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
				"next_action":  {Type: genai.TypeString, Enum: []string{"CONTINUE_CHAT", "APPOINTMENT"}},
				"reply":        {Type: genai.TypeString},
				"doctor_id":    {Type: genai.TypeString},
				"prediagnosis": {Type: genai.TypeString},
			},
			Required: []string{"next_action", "reply", "doctor_id", "prediagnosis"},
		},
	}

	// add the stored messages to the history
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

	// get the response from the LLM
	if res != nil && len(res.Candidates) > 0 && res.Candidates[0].Content != nil &&
		len(res.Candidates[0].Content.Parts) > 0 {
		text := res.Candidates[0].Content.Parts[0].Text
		return text, nil
	}

	return "", fmt.Errorf("no response from LLM")
}

func UpdateChatHistory(sessionId string, newMessage string, LLMResponse string) error {

	// get the session from the database
	var session models.Session
	err := config.DB.Where("id = ?", sessionId).First(&session).Error
	if err != nil {
		return fmt.Errorf("error fetching session: %v", err)
	}

	// append the new message and LLM response to the chat history
	now := time.Now()
	newUserMessage := models.Message{Role: "user", Content: newMessage, SessionID: session.ID, CreatedAt: now, UpdatedAt: now}
	newLLMResponse := models.Message{Role: "omsapa", Content: LLMResponse, SessionID: session.ID, CreatedAt: now.Add(time.Millisecond), UpdatedAt: now.Add(time.Millisecond)}

	// save the new messages to the database
	if err := config.DB.Create(&newUserMessage).Error; err != nil {
		return fmt.Errorf("error saving user message: %v", err)
	}

	if err := config.DB.Create(&newLLMResponse).Error; err != nil {
		return fmt.Errorf("error saving LLM response: %v", err)
	}

	return nil
}

func buildSystemPrompt(session *models.Session) string {
	// Build the system prompt using the user's data
	userDataText := fmt.Sprintf(
		"\nHere's the user's data: \n\nName:%s\nAge:%s\nGender:%s\nNationality:%s\nWeight: %d\nHeight: %d\nHeartrate: %d\nBodytemp: %f\n",
		session.User.Name,
		utils.DateToAgeString(session.User.DOB),
		session.User.Gender,
		session.User.Nationality,
		session.Weight,
		session.Height,
		session.Heartrate,
		session.Bodytemp,
	)

	doctors := GetAllDoctors()

	// Convert the doctors to a string representation
	var doctorList []string
	for _, doctor := range doctors {
		doctorList = append(doctorList, fmt.Sprintf("- [%s] %s (%s)\n", doctor.ID, doctor.Name, doctor.Specialty))
	}
	doctorListText := fmt.Sprintf("\nHere are the doctors available [ID] Name (Specialty):\n%s", strings.Join(doctorList, ""))

	// Get history of sessions
	var history []models.Session = GetHistory(session)

	// Convert history of sessions to a string representation
	var historyList []string
	if len(history) > 0 {
		for _, sessionItem := range history {
			historyList = append(historyList, fmt.Sprintf(
				"[%s]\nWeight: %f\nHeight: %f\nHeartrate: %f\nBodytemp: %f\nPrediagnosis: %s\n",
				sessionItem.CreatedAt.Format("2006-01-02 15:04:05"),
				sessionItem.Weight,
				sessionItem.Height,
				sessionItem.Heartrate,
				sessionItem.Bodytemp,
				sessionItem.Prediagnosis,
			))
		}
	} else {
		historyList = append(historyList, "No previous sessions found.\n")
	}
	historyListText := fmt.Sprintf("\nHere are the previous sessions:\n%s", strings.Join(historyList, ""))

	// hardcode initial system prompt
	systemPrompt := `You are a health expert fluent in Indonesian, passionate about helping patients understand their symptoms and connect them with the right doctor. Your goal is to guide patients step-by-step, choose a doctor based on their symptoms and the available doctor list, and ensure they feel comfortable and informed.
Follow the conversation flow and output format strictly as described below. It is VERY IMPORTANT that you adhere to the JSON format:

0.  JSON FORMAT FOR EVERY RESPONSE
-  Output the response in valid JSON format ONLY. Do not include any surrounding text, explanations, or formatting outside the JSON structure.
-  JSON Output Format:
	{
	\"next_action\": \"CONTINUE_CHAT\" or \"APPOINTMENT\",
	\"reply\": \"Your text reply here\",
	\"doctor_id\": \"selected doctor_id\" (only if next_action is APPOINTMENT),
	\"prediagnosis\": \"Your pre-diagnosis based on the conversation\" (only if next_action is APPOINTMENT)
	}
-  Detailed explanation of each field:
	- next_action: A string indicating the next step in the conversation. Must be either \"CONTINUE_CHAT\" or \"APPOINTMENT\".
	- reply: A string containing your response to the patient.
		-  If next_action is \"CONTINUE_CHAT\", this should be the next question(s) or statement to keep the conversation flowing.
		-  If next_action is \"APPOINTMENT\", this should be a confirmation message to the patient, informing them of the doctor they are assigned to and that their queue number has been sent to their email.  Be friendly and reassuring. For example: \"Based on your symptoms, I recommend you see Dr. Udin (General Practitioner). Your queue number has been sent to your email address.\"
		doctor_id: A string containing the ID of the selected doctor. This *MUST be included if and only if next_action is \"APPOINTMENT\".  You MUST choose a doctor from the provided list of doctors. If no doctor seems appropriate based on the conversation, choose a General Practitioner.
		prediagnosis: A string containing your pre-diagnosis based on the conversation. This *MUST be included if and only if next_action is \"APPOINTMENT\".  Be brief and provide a likely possible diagnosis.
-  Example JSON Response (for CONTINUE_CHAT):
	{
	\"next_action\": \"CONTINUE_CHAT\",
	\"reply\": \"Can you describe the location of the pain more specifically?  Is it sharp, dull, or throbbing?\",
	\"doctor_id\": null,
	\"prediagnosis\": null
	}
-  Example JSON Response (for APPOINTMENT):
	{
	\"next_action\": \"APPOINTMENT\",
	\"reply\": \"Based on your symptoms, I recommend you see Dr. Jane Doe (Cardiologist). Your queue number has been sent to your email address.\",
	\"doctor_id\": \"edd248b7-75d3-4af2-a954-183970124e9d\",
	\"prediagnosis\": \"Possible arrhythmia\"
	}

1. Conversation Flow:
	1. First Response:
		- Greet the patient warmly and ask them to choose their preferred language:
			- a. Bahasa Indonesia
			- b. English
		- Add: \"Choose the language that makes you feel most comfortable.\"
		- Default language: Bahasa Indonesia.
	2. Second Response (AFTER language selection):
		- Start with a friendly greeting.
		- Ask 3 simple, easy-to-understand questions to begin. Provide 3 quick-answer examples in parentheses for each question.
	3. Follow-Up Questions:
		- Based on the patient's answers, ask progressively specific follow-up questions (e.g., symptom type, duration, severity, associated symptoms, medication use).
		- Try to understand the patient's condition, don't try to just pass the problem to the doctor.
		- Limit to 3 questions per follow-up. Provide 3 quick-answer examples in parentheses for each question.
	4. Decision Points:
		- If the conversation is sufficient for a preliminary diagnosis OR the user requests an appointment:
			- Generate a pre-diagnosis.
			- Select a suitable doctor_id from the provided doctor list. If no suitable doctor is available based on the conversation, assign the patient to a General Practitioner.
			- Make the next_action \"APPOINTMENT\".
			- Don't assume their sickness based on the symptoms, ask their symptoms first.
2. Important Notes:
	- Always adhere strictly to the JSON format and the defined conversation flow.
	- Prioritize patient comfort and understanding throughout the interaction.
	- Ensure the JSON output is valid and contains no additional text or formatting outside the JSON structure.
	- When next_action is \"APPOINTMENT\",  ALWAYS populate the doctor_id and prediagnosis fields using the information you have gathered.  If you are uncertain about the prediagnosis, give the most likely possibility.
	- If you are unable to determine the doctor_id from the symptoms the patient is providing, default to a General Practitioner from the list.  Do not return an empty doctor_id.`

	// Build the system prompt text
	systemPromptText := fmt.Sprintf("%s %s %s %s\nCurrent Time: %s", systemPrompt, userDataText, doctorListText, historyListText, time.Now().Format("2006-01-02 15:04:05"))

	return systemPromptText
}

func GetSessionData(sessionId string) (models.Session, error) {
	// check if session_id exists in the database
	var session models.Session

	err := config.DB.
		Preload("User").
		Preload("Messages", func(db *gorm.DB) *gorm.DB {
			return db.Order("created_at ASC") // Order messages by created_at in descending order (for latest messages first)
		}).
		Where("id = ?", sessionId).First(&session).Error

	if err != nil {
		return models.Session{}, fmt.Errorf("session not found: %w", err)
	}

	return session, nil
}

func GetHistory(session *models.Session) []models.Session {
	var history []models.Session
	err := config.DB.Where("user_id = ?", session.User.ID).Where("id != ?", session.ID).Order("created_at DESC").Find(&history).Error
	if err != nil {
		log.Printf("Error fetching session history: %v\n", err)
		return []models.Session{} // Return an empty slice if there's an error
	}

	return history
}

func DoctorDiagnose(sessionId string, diagnosis string) error {
	// Fetch the session from the database
	var session models.Session
	err := config.DB.Where("id = ?", sessionId).First(&session).Error
	if err != nil {
		return fmt.Errorf("session not found: %w", err)
	}

	// Update the session with the doctor's diagnosis
	session.DoctorDiagnosis = diagnosis
	session.UpdatedAt = time.Now()

	if err := config.DB.Save(&session).Error; err != nil {
		return fmt.Errorf("failed to save diagnosis: %w", err)
	}

	return nil
}

// RemoveMarkdownAndExtractJSON removes Markdown syntax and extracts the JSON content
func ParseJSON(input string) (schemas.LLMResponse, error) {

	// Extract JSON content
	var responseJSON schemas.LLMResponse
	err := json.Unmarshal([]byte(input), &responseJSON)
	if err != nil {
		return schemas.LLMResponse{}, fmt.Errorf("failed to extract JSON: %w", err)
	}

	return responseJSON, nil
}

func GetSessionsByUserID(userID uuid.UUID) []models.Session {
	var sessions []models.Session
	err := config.DB.Where("user_id = ?", userID).Find(&sessions).Error
	if err != nil {
		return nil
	}
	return sessions
}
