// Package provides way to CRUD operations with User Collection
package db

import (
	"errors"
	"fileguard/utils"
	"github.com/google/uuid"
	"google.golang.org/api/iterator"
	"log"
	"time"
)

func (db *Database) CreateNewUser(username, email string) (string, error) {
	query := db.Client.Collection("Users").Where("email", "==", email).Limit(1)
	iter := query.Documents(db.ctx)
	defer iter.Stop()

	_, err := iter.Next()
	if err != iterator.Done {
		return "", utils.ErrUserAlreadySigned
	}

	if !emailRegex.MatchString(email) {
		return "", errors.New("Invalid email address")
	}

	currentTime := time.Now()
	userID := uuid.New().String()

	_, _, err = db.Client.Collection("Users").Add(db.ctx, map[string]interface{}{
		"user_id":       userID,
		"username":      username,
		"email":         email,
		"creation_date": currentTime,
	})

	if err != nil {
		return "", err
	}

	return userID, nil
}

func (db *Database) GetUserByID(userID string) (map[string]interface{}, error) {
	// Get user document reference
	doc, err := db.getRecord("Users", "user_id", userID)
	if err != nil {
		return nil, err
	}

	return doc.Data(), nil
}

func (db *Database) GetUserByEmail(email string) (map[string]interface{}, error) {
	doc, err := db.getRecord("Users", "email", email)
	if err != nil {
		return nil, err
	}
	return doc.Data(), nil
}

func (db *Database) DeleteUserByID(userID string) error {
	doc, err := db.getRecord("Users", "user_id", userID)
	if err != nil {
		log.Fatalf("Error getting user: %v", err)
		return err
	}

	if _, err = doc.Ref.Delete(db.ctx); err != nil {
		log.Fatalf("Error deleting user: %v", err)
		return err
	}

	return nil

}
