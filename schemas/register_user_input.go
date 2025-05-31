package schemas

type RegisterUserInput struct {
	Name        string  `json:"name" validate:"required"`
	Email       string  `json:"email" validate:"required,email"`
	Nationality string  `json:"nationality" validate:"required"`
	DOB         string  `json:"dob" validate:"required"`
	Weight      float32 `json:"weight" validate:"required"`
	Height      float32 `json:"height" validate:"required"`
	Heartrate   float32 `json:"heartrate" validate:"required"`
	Bodytemp    float32 `json:"bodytemp" validate:"required"`
	Gender      string  `json:"gender" validate:"required"`
}
