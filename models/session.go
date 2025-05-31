package models

import (
	"time"

	"github.com/google/uuid"
)

type Session struct {
	ID              uuid.UUID `json:"id" gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	UserID          uuid.UUID `json:"user_id" gorm:"type:uuid;not null"`
	User            User      `json:"user" gorm:"foreignKey:UserID"`
	Weight          float32   `json:"weight" gorm:"type:float;not null"`
	Height          float32   `json:"height" gorm:"type:float;not null"`
	Heartrate       float32   `json:"heartrate" gorm:"type:float;not null"`
	Bodytemp        float32   `json:"bodytemp" gorm:"type:float;not null"`
	Messages        []Message `json:"messages" gorm:"foreignKey:SessionID"`
	Prediagnosis    string    `json:"prediagnosis" gorm:"type:varchar(100);"`
	DoctorDiagnosis string    `json:"doctor_diagnosis" gorm:"type:varchar(100);"`
	CreatedAt       time.Time `json:"created_at" gorm:"type:timestamp;not null"`
	UpdatedAt       time.Time `json:"updated_at" gorm:"type:timestamp;not null"`
}
