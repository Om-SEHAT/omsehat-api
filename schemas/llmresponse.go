package schemas

type LLMResponse struct {
	NextAction string `json:"next_action"`
	Reply      string `json:"reply"`
	DoctorID  string `json:"doctor_id"`
	PreDiagnosis string `json:"prediagnosis"`
}