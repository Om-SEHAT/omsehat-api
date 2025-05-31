package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID          uuid.UUID `json:"id" gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	Name        string    `json:"name" gorm:"type:varchar(100);not null"`
	Email       string    `json:"email" gorm:"type:varchar(100);unique;not null"`
	Nationality string    `json:"nationality" gorm:"type:varchar(100);not null"`
	DOB         string    `json:"dob" gorm:"type:date;not null"`
	Gender      string    `json:"gender" gorm:"type:varchar(10);not null"`
	OTP         string    `json:"-" gorm:"type:varchar(6)"`
	CreatedAt   time.Time `json:"created_at" gorm:"type:timestamp;not null"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"type:timestamp;not null"`
	Sessions    []Session `json:"sessions" gorm:"foreignKey:UserID"`
}
