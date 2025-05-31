package services

import (
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/Om-SEHAT/omsehat-api/config"
	"github.com/Om-SEHAT/omsehat-api/models"
	"github.com/Om-SEHAT/omsehat-api/schemas"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func RegisterUser(input schemas.RegisterUserInput) (*models.User, error) {
	// Generate OTP
	otp := generateOTP()

	// Check if the user already exists in the database
	var existingUser models.User
	err := config.DB.Preload("Sessions").Where("email = ?", input.Email).First(&existingUser).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		// User does not exist, create a new user
		newUser := models.User{
			Name:        input.Name,
			Email:       input.Email,
			Nationality: input.Nationality,
			DOB:         input.DOB,
			Gender:      input.Gender,
			OTP:         otp,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}

		if err := config.DB.Create(&newUser).Error; err != nil {
			return nil, fmt.Errorf("failed to create user: %w", err)
		}
		existingUser = newUser

	} else if err == nil {
		// update existing user with new OTP
		existingUser.OTP = otp
		existingUser.Name = input.Name
		existingUser.Nationality = input.Nationality
		existingUser.DOB = input.DOB
		existingUser.Gender = input.Gender

		existingUser.UpdatedAt = time.Now()

		if err := config.DB.Save(&existingUser).Error; err != nil {
			return nil, fmt.Errorf("failed to update user: %w", err)
		}
	} else {
		// Some other error occurred while checking for existing user
		return nil, fmt.Errorf("failed to check for existing user: %w", err)
	}

	// Send OTP email to the user
	_, err = sendOTPEmail(existingUser.Email, otp, os.Getenv("EMAIL_OTP_TOKEN"))
	if err != nil {
		log.Printf("Error sending OTP email: %v\n", err)
	}

	// registration success, return the user object
	return &existingUser, nil
}

func GetUserByID(userID uuid.UUID) *models.User {
	var user models.User
	err := config.DB.First(&user, "id = ?", userID).Error
	if err != nil {
		return nil
	}
	return &user
}
