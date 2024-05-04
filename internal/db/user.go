// Package provides way to CRUD operations with User Collection
package db

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"google.golang.org/api/iterator"
	"log"
	"time"
)

func (db *Database) CreateNewUser(username, email string) error {
	query := db.Client.Collection("Users").Where("email", "==", email).Limit(1)
	iter := query.Documents(db.ctx)
	defer iter.Stop()

	_, err := iter.Next()
	if err != iterator.Done {
		return nil
	}

	if !emailRegex.MatchString(email) {
		return errors.New("Invalid email address")
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
		return err
	}

	return nil
}

func (db *Database) GetUserByID(userID string) (map[string]interface{}, error) {
	// Get user document reference
	query := db.Client.Collection("Users").Where("user_id", "==", userID).Limit(1)

	// Execute the query
	iter := query.Documents(db.ctx)
	defer iter.Stop()

	// Retrieve the first document from the query result
	doc, err := iter.Next()
	if err == iterator.Done {
		return nil, fmt.Errorf("user with ID %s not found", userID)
	}
	if err != nil {
		return nil, err
	}

	// Extract user data from the document
	userData := doc.Data()

	return userData, nil
}

func (db *Database) DeleteUserByID(userID string) {
	query := db.Client.Collection("Users").Where("user_id", "==", userID).Limit(1)

	iter := query.Documents(db.ctx)
	defer iter.Stop()

	doc, err := iter.Next()

	if err == iterator.Done {
		log.Fatal("user with ID %s not found", userID)
		return
	}
	if err != nil {
		log.Fatalf("Error getting user: %v", err)
		return
	}

	// Delete the user document
	_, err = doc.Ref.Delete(db.ctx)
	if err != nil {
		log.Fatalf("Error deleting user: %v", err)
		return
	}

	log.Println("User Deleted")
}
