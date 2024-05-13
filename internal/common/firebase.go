package common

import (
	"context"
	"fileguard/utils"
	firebase "firebase.google.com/go"
	"google.golang.org/api/option"
	"log"
	"sync"
)

var (
	firebaseApp *firebase.App
	firebaseCtx context.Context
)

var lock = &sync.Mutex{}

func GetFirebaseApp() (*firebase.App, context.Context, error) {
	lock.Lock()
	defer lock.Unlock()

	if firebaseApp != nil {
		log.Println("Single instance already created.")
		return firebaseApp, firebaseCtx, nil
	}

	log.Println("Creating instance now.")
	config := &firebase.Config{
		StorageBucket: "fileguard-cf4d3.appspot.com",
	}

	firebaseCtx = context.Background()

	opt := option.WithCredentialsFile(utils.FirebaseCredentialsFile)
	app, err := firebase.NewApp(firebaseCtx, config, opt)
	if err != nil {
		return nil, nil, err
	}

	firebaseApp = app

	return firebaseApp, firebaseCtx, nil
}
