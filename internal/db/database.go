package db

import (
	"cloud.google.com/go/firestore"
	"context"
	firebase "firebase.google.com/go"
	"google.golang.org/api/option"
	"log"
)

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
func (db *Database) CreateNewUser() {
	_, _, err := db.Client.Collection("Users").Add(db.ctx, map[string]interface{}{
		"first": "Ada",
		"last":  "Lovelace",
		"born":  1815,
	})

	if err != nil {
		log.Fatal("Fail")
	}
}
