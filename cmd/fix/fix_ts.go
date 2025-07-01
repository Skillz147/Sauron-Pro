package main

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go/v4"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

func main() {
	ctx := context.Background()
	app, err := firebase.NewApp(ctx, nil, option.WithCredentialsFile("../../firebaseAdmin.json"))
	if err != nil {
		log.Fatalf("❌ Firebase init failed: %v", err)
	}

	client, err := app.Firestore(ctx)
	if err != nil {
		log.Fatalf("❌ Firestore client init failed: %v", err)
	}
	defer client.Close()

	// List all user result sets
	users, err := client.Collection("results").Documents(ctx).GetAll()
	if err != nil {
		log.Fatalf("❌ Failed to get user docs: %v", err)
	}

	for _, userDoc := range users {
		userID := userDoc.Ref.ID
		fmt.Println("🔍 Checking user:", userID)

		entries := client.Collection("results").Doc(userID).Collection("entries")
		iter := entries.Documents(ctx)

		for {
			doc, err := iter.Next()
			if err == iterator.Done {
				break
			}
			if err != nil {
				log.Printf("⚠️ Failed to read doc: %v", err)
				continue
			}

			data := doc.Data()
			fmt.Printf("🔎 Doc %s fields:\n", doc.Ref.ID)
			for k, v := range data {
				fmt.Printf("   • %s (%T): %v\n", k, v, v)
			}

			var raw interface{}
			for k, v := range data {
				if strings.EqualFold(k, "ts") {
					raw = v
					break
				}
			}

			if raw == nil {
				fmt.Println("   ❌ No TS/ts field found, skipping")
				continue
			}

			var ts int64
			switch v := raw.(type) {
			case int64:
				ts = normalizeTS(v)
			case int:
				ts = normalizeTS(int64(v))
			case float64:
				ts = normalizeTS(int64(v))
			case string:
				parsed, err := strconv.ParseInt(v, 10, 64)
				if err != nil {
					fmt.Println("   ⚠️ Cannot parse string TS:", err)
					continue
				}
				ts = normalizeTS(parsed)
			default:
				fmt.Printf("   ⚠️ Unknown TS type %T\n", v)
				continue
			}

			// Prepare update ops
			updates := []firestore.Update{
				{Path: "ts", Value: ts},
			}

			// Delete all mis-cased variants
			for _, k := range []string{"TS", "Ts", "tS"} {
				if _, ok := data[k]; ok {
					updates = append(updates, firestore.Update{Path: k, Value: firestore.Delete})
				}
			}

			_, err = doc.Ref.Update(ctx, updates)
			if err != nil {
				log.Printf("❌ Failed to patch doc %s: %v", doc.Ref.ID, err)
			} else {
				fmt.Printf("✅ Fixed doc %s → ts: %d\n", doc.Ref.ID, ts)
			}
		}
	}

	fmt.Println("🎯 Migration done.")
}

func normalizeTS(input int64) int64 {
	if input < 1e12 {
		return input * 1000 // convert seconds → ms
	}
	return input // already ms
}
