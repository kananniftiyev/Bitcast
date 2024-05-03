package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/google/uuid" // Import package for generating UUIDs
)

var ErrUserAlreadyExists = errors.New("User Already exists")

// Token represents an authentication token
type Token struct {
	UserID      string
	AccessToken string
	Expiration  time.Time
}

// GenerateToken generates a token for the given user
func GenerateToken(userID string) *Token {
	accessToken := generateRandomToken()
	expiration := time.Now().Add(24 * time.Hour) // Example: Token expires in 24 hours

	return &Token{
		UserID:      userID,
		AccessToken: accessToken,
		Expiration:  expiration,
	}
}

// SaveToken saves the token to a file
func SaveToken(token *Token) error {
	dir, err := os.Getwd()
	if err != nil {
		return err
	}

	// Define the filename for the token file
	filename := filepath.Join(dir, "save_data.json")

	// Check if the file already exists
	if _, err := os.Stat(filename); err == nil {
		return fmt.Errorf("token file already exists")
	}

	data, err := json.Marshal(token)
	if err != nil {
		return err
	}

	// Todo: change this to docs
	return ioutil.WriteFile(filename, data, 0644)
}

// LoadToken loads the token from a JSON file in the same folder as the project
func LoadToken() (*Token, error) {
	dir, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	// Todo: change this to docs
	filename := filepath.Join(dir, "save_data.json")
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	var token Token
	err = json.Unmarshal(data, &token)
	if err != nil {
		return nil, err
	}
	return &token, nil
}

// generateRandomToken generates a random access token
func generateRandomToken() string {
	uuidWithHyphen := uuid.New()
	uuid := uuidWithHyphen.String()
	// Remove hyphens from the generated UUID to get a random token
	randomToken := removeHyphens(uuid)
	return randomToken
}

// removeHyphens removes hyphens from a string
func removeHyphens(s string) string {
	var result string
	for _, char := range s {
		if char != '-' {
			result += string(char)
		}
	}
	return result
}
