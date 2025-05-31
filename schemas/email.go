package schemas

type Email struct {
	To      string   `json:"to" binding:"required,email"`
	From    string   `json:"from_email,omitempty" binding:"omitempty,email"`
	Subject string   `json:"subject" binding:"required"`
	Body    string   `json:"body" binding:"required"`
	HTML    string   `json:"html,omitempty"`
	CC      []string `json:"cc,omitempty"`
	BCC     []string `json:"bcc,omitempty"`
}
