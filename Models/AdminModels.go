package models

import (
	"database/sql"
	"fmt"
)

// Users represents the user model
type Users struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	UserID   int    `json:"id"`
}

// RegistrationCheck checks user credentials and returns a result message
func RegistrationCheck(Email, Password string) (string, error) {
	var userID int

	// Query to check for a matching user
	LoginQuery := `SELECT user_id FROM admin_registration WHERE email = ? AND password = ?`

	// Execute the query
	err := DB.QueryRow(LoginQuery, Email, Password).Scan(&userID)
	if err != nil {
		if err == sql.ErrNoRows {
			// No matching user found
			return "Invalid user", nil
		}
		// Other errors
		fmt.Printf("Error querying user: %v\n", err)
		return "", err
	}

	// Login successful
	return "Success", nil
}
