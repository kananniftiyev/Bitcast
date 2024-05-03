// Package db provides functionalities to interact with Firestore database.
package db

import (
	"cloud.google.com/go/firestore"
	"context"
	"fileguard/utils"
	firebase "firebase.google.com/go"
	"google.golang.org/api/option"
	"log"
	"regexp"
)

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

type Database struct {
	Client *firestore.Client
	ctx    context.Context
}

func NewDatabase() (*Database, error) {
	ctx := context.Background()
	sa := option.WithCredentialsFile(utils.FirebaseCredentialsFile)
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
