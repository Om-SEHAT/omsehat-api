package utils

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// BindAndValidate is a utility function to bind JSON input and validate the struct.
func BindAndValidate(c *gin.Context, input interface{}) (bool, map[string]string) {
	// bind the request body to the input struct
	if err := c.ShouldBindJSON(input); err != nil {
		// log the error for debugging
		log.Printf("Error binding JSON: %v", err)
		c.JSON(400, gin.H{"message": "Invalid JSON format"})
		return false, nil
	}

	// validate the input struct
	validate := validator.New()
	if err := validate.Struct(input); err != nil {
		// collect validation errors in a map
		errors := map[string]string{}
		for _, err := range err.(validator.ValidationErrors) {
			errors[err.Field()] = err.Tag()
		}
		c.JSON(400, gin.H{
			"message": "Validation failed",
			"details": errors,
		})
		return false, errors
	}

	// everything is good
	return true, nil
}
