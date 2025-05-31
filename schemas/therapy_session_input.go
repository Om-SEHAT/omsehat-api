package schemas

type TherapySessionInput struct {
	StressLevel    int    `json:"stress_level" validate:"required,min=1,max=10"`
	MoodRating     int    `json:"mood_rating" validate:"required,min=1,max=10"`
	SleepQuality   int    `json:"sleep_quality" validate:"required,min=1,max=10"`
	IsHealthWorker bool   `json:"is_health_worker" validate:"required"`
	Specialization string `json:"specialization" validate:"omitempty,min=2,max=50"` // Required only if IsHealthWorker is true
	WorkHours      int    `json:"work_hours" validate:"omitempty,min=1,max=168"`    // Required only if IsHealthWorker is true, max hours in a week
}
