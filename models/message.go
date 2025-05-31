package models

import (
	"time"

	"github.com/google/uuid"
)

type Message struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	Role      string    `json:"role" gorm:"type:varchar(50);not null"`
	Content   string    `json:"content" gorm:"type:text;not null"`
	SessionID uuid.UUID `json:"session_id" gorm:"type:uuid;not null"`
	Session   Session   `json:"-" gorm:"foreignKey:SessionID"`
	CreatedAt time.Time `json:"created_at" gorm:"type:timestamp;default:now()"`
	UpdatedAt time.Time `json:"updated_at" gorm:"type:timestamp;default:now()"`
}
