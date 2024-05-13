// Package db provides functionalities to interact with Firestore database.
package db

import (
	"cloud.google.com/go/firestore"
	"context"
	"fileguard/internal/common"
	"fmt"
	"google.golang.org/api/iterator"
	"regexp"
)

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

type Database struct {
	Client *firestore.Client
	ctx    context.Context
}

func NewDatabase() (*Database, error) {
	app, ctx, err := common.GetFirebaseApp()
	if err != nil {
		panic(err)
	}

	client, err := app.Firestore(ctx)
	if err != nil {
		return nil, err
	}

	return &Database{Client: client, ctx: ctx}, nil
}

func (db *Database) getRecord(collectionName, fieldName, value string) (*firestore.DocumentSnapshot, error) {
	query := db.Client.Collection(collectionName).Where(fieldName, "==", value).Limit(1)
	iter := query.Documents(db.ctx)
	defer iter.Stop()
	doc, err := iter.Next()
	if err == iterator.Done {
		return nil, fmt.Errorf("%s record not found", collectionName)
	}
	if err != nil {
		return nil, err
	}
	return doc, nil
}
