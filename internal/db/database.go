package db

import (
	"cloud.google.com/go/firestore"
	"context"
	firebase "firebase.google.com/go"
	"fmt"
	"github.com/google/uuid"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
	"log"
	"regexp"
	"time"
)

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

type Database struct {
	Client *firestore.Client
	ctx    context.Context
}

func NewDatabase(projectID string) (*Database, error) {
	ctx := context.Background()
	sa := option.WithCredentialsFile("C:/Users/kenan/Documents/GitHub/fileguard/internal/db/fileguard.json")
	conf := &firebase.Config{ProjectID: "fileguard-cf4d3"}
	app, err := firebase.NewApp(ctx, conf, sa)
	if err != nil {
		log.Fatalln(err)
	}

	client, err := app.Firestore(ctx)
	if err != nil {
		return nil, err
	}

	return &Database{Client: client, ctx: ctx}, nil

}

func (db *Database) CreateNewUser(username, email, hashed_password string) {
	if !emailRegex.MatchString(email) {
		log.Fatal("Invalid email address")
		return
	}

	currentTime := time.Now()
	userID := uuid.New().String()

	_, _, err := db.Client.Collection("Users").Add(db.ctx, map[string]interface{}{
		"user_id":       userID,
		"username":      username,
		"email":         email,
		"password":      hashed_password,
		"creation_date": currentTime,
	})

	if err != nil {
		log.Fatal("Fail")
	}
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
