package models

import (
	"time"

	"github.com/google/uuid"
)

type TherapySession struct {
	ID             uuid.UUID `json:"id" gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	UserID         uuid.UUID `json:"user_id" gorm:"type:uuid;not null"`
	User           User      `json:"user" gorm:"foreignKey:UserID"`
	StressLevel    int       `json:"stress_level" gorm:"type:int;not null"`  // Scale 1-10
	MoodRating     int       `json:"mood_rating" gorm:"type:int;not null"`   // Scale 1-10
	SleepQuality   int       `json:"sleep_quality" gorm:"type:int;not null"` // Scale 1-10
	IsHealthWorker bool      `json:"is_health_worker" gorm:"type:boolean;not null"`
	Specialization string    `json:"specialization" gorm:"type:varchar(100);"` // For healthcare workers
	WorkHours      int       `json:"work_hours" gorm:"type:int;"`              // Weekly work hours for healthcare workers
	Messages       []Message `json:"messages" gorm:"foreignKey:SessionID"`
	Recommendation string    `json:"recommendation" gorm:"type:text;"`
	NextSteps      string    `json:"next_steps" gorm:"type:text;"`
	CreatedAt      time.Time `json:"created_at" gorm:"type:timestamp;not null"`
	UpdatedAt      time.Time `json:"updated_at" gorm:"type:timestamp;not null"`
}
