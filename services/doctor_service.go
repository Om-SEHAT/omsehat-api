package services

import (
	"github.com/Om-SEHAT/omsehat-api/config"
	"github.com/Om-SEHAT/omsehat-api/models"
	"github.com/google/uuid"
)

func GetAllDoctors() []models.Doctor {
	var doctors []models.Doctor
	err := config.DB.Find(&doctors).Error
	if err != nil {
		return nil
	}
	return doctors
}

func GetDoctorByID(doctorID string) *models.Doctor {
	var doctor models.Doctor
	id, err := uuid.Parse(doctorID)
	if err != nil {
		return nil
	}
	err = config.DB.First(&doctor, "id = ?", id).Error
	if err != nil {
		return nil
	}
	return &doctor
}
