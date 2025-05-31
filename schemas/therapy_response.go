package schemas

type TherapyResponse struct {
	NextAction     string `json:"next_action"`    // CONTINUE_CHAT or END_SESSION
	Reply          string `json:"reply"`          // The chatbot's response
	Recommendation string `json:"recommendation"` // Mental health recommendation
	NextSteps      string `json:"next_steps"`     // Suggested actions for the user
}
