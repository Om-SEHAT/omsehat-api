package models

type Doctor struct {
	ID        string `json:"id" gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	Name      string `json:"name" gorm:"type:varchar(100);not null"`
	Email     string `json:"email" gorm:"type:varchar(100);unique;not null"`
	Specialty string `json:"specialty" gorm:"type:varchar(100);not null"`
	Roomno 	 string `json:"roomno" gorm:"type:varchar(10);not null"`
}
