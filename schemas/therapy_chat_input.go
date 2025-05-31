package schemas

type TherapyChatInput struct {
	NewMessage string `json:"new_message" validate:"required"`
}
