package utils

import (
	"fmt"
	"log"
	"time"
)

// calculate age (year, month, day) from date of birth
func DateToAgeString(dateOfBirth string) string {
	// Strip the time and timezone portion if present (assumes format is RFC3339)
	dateOfBirth = dateOfBirth[:10] // Keeps only the YYYY-MM-DD part

	// Parse the date string into a time.Time object
	dob, err := time.Parse("2006-01-02", dateOfBirth)
	if err != nil {
		// log error
		log.Printf("Error parsing date of birth: %v\n", err)
		return ""
	}

	// Get the current date
	currentDate := time.Now()

	// Calculate the difference between the current date and the date of birth
	years := currentDate.Year() - dob.Year()

	months := currentDate.Month() - dob.Month()
	if months < 0 {
		years--
		months += 12
	}
	days := currentDate.Day() - dob.Day()
	if days < 0 {
		months--
		days += 30 // Approximation, can be improved for exact days in month
		if months < 0 {
			years--
			months += 12
		}
	}

	// Format the age string
	ageString := fmt.Sprintf("%d years, %d months, %d days", years, months, days)

	return ageString
}
