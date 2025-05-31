package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/Om-SEHAT/omsehat-api/config"
	"github.com/Om-SEHAT/omsehat-api/models"
	"github.com/Om-SEHAT/omsehat-api/schemas"
)

func generateOTP() string {
	// Generate a random 6-digit OTP
	otp := fmt.Sprintf("%06d", rand.Intn(1000000))
	return otp
}

func ValidateOTP(input schemas.OTPInput) (*models.Session, error) {
	// get the user from the input
	var user models.User
	err := config.DB.Where("email = ?", input.Email).First(&user).Error
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	// check if OTP sent to the user is valid
	if user.OTP != input.OTP {
		return nil, fmt.Errorf("invalid OTP")
	}

	// create a new session for the user with the data from the input
	newSession := models.Session{
		UserID:    user.ID,
		Weight:    input.Weight,
		Height:    input.Height,
		Heartrate: input.Heartrate,
		Bodytemp:  input.Bodytemp,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// save the session to the database
	err = config.DB.Create(&newSession).Error
	if err != nil {
		return nil, fmt.Errorf("failed to create session: %w", err)
	}

	// update user's OTP to nil after successful validation
	user.OTP = ""
	err = config.DB.Save(&user).Error
	if err != nil {
		return nil, fmt.Errorf("failed to update user OTP: %w", err)
	}

	return &newSession, nil
}

func sendOTPEmail(to string, otp string, token string) (map[string]interface{}, error) {

	email := schemas.Email{
		To:      to,
		Subject: "Your OTP Code",
		Body:    "",
		From:    "omsehat@sportsnow.app",
		HTML:    injectOtpIntoHtml(otp),
	}

	// URL of the email service
	url := "http://52.230.88.220:16250/send-email"

	jsonData, err := json.Marshal(email)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal email request: %w", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send HTTP request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("error from email service: %s", resp.Status)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return result, nil
}

func injectOtpIntoHtml(otpCode string) string {
	// Path to the HTML file
	htmlFilePath := "emails/otp_mail.html"

	// Read the HTML file
	htmlBytes, err := os.ReadFile(htmlFilePath)
	if err != nil {
		log.Fatalf("Failed to read HTML file: %v", err)
	}

	// Convert the file content to a string
	htmlString := string(htmlBytes)

	// Replace the placeholder with the OTP code
	return strings.ReplaceAll(htmlString, "{{otp_code}}", otpCode)
}
