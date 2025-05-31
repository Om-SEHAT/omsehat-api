package schemas

type SessionChatInput struct {
	NewMessage string `json:"new_message" validate:"required"`
}
