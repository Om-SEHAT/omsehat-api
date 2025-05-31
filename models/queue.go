package models

import (
	"time"

	"github.com/google/uuid"
)

type Queue struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	DoctorID  uuid.UUID `json:"doctor_id" gorm:"type:uuid;not null"`
	Doctor    Doctor    `json:"doctor" gorm:"foreignKey:DoctorID"`
	SessionID uuid.UUID `json:"session_id" gorm:"type:uuid;not null"`
	Session   Session   `json:"session" gorm:"foreignKey:SessionID"`
	CreatedAt time.Time `json:"created_at" gorm:"type:timestamp;not null"`
	UpdatedAt time.Time `json:"updated_at" gorm:"type:timestamp;not null"`
	Number    int       `json:"number" gorm:"type:int;not null"`
}
