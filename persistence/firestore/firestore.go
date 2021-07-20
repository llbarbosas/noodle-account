package firestore

import (
	"context"
	"log"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"google.golang.org/api/option"
)

type Repository struct {
	client *firestore.Client
}

func NewRepository(t interface{}) Repository {
	ctx := context.Background()
	sa := option.WithCredentialsFile("path/to/serviceAccount.json")

	app, err := firebase.NewApp(ctx, nil, sa)

	if err != nil {
		log.Fatalln(err)
	}

	client, err := app.Firestore(ctx)

	if err != nil {
		log.Fatalln(err)
	}

	return Repository{
		client: client,
	}
}
