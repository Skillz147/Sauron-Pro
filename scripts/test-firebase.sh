#!/bin/bash

set -e

echo "🔥 Testing Firebase Connection..."

# Check if we're in the right directory
if [ ! -f "firebaseAdmin.json" ]; then
    echo "❌ firebaseAdmin.json not found. Make sure you're in the Sauron directory."
    exit 1
fi

# Create the Go test file
cat > test_firebase_connection.go << 'EOF'
package main

import (
	"context"
	"fmt"
	"os"
	"time"

	firebase "firebase.google.com/go/v4"
	"google.golang.org/api/option"
)

func main() {
	fmt.Println("🔥 Testing Firebase connection...")
	
	// Check if credentials file exists
	if _, err := os.Stat("firebaseAdmin.json"); os.IsNotExist(err) {
		fmt.Println("❌ firebaseAdmin.json not found")
		os.Exit(1)
	}
	
	ctx := context.Background()
	opt := option.WithCredentialsFile("firebaseAdmin.json")
	
	// Initialize Firebase app
	fmt.Println("📡 Initializing Firebase app...")
	app, err := firebase.NewApp(ctx, nil, opt)
	if err != nil {
		fmt.Printf("❌ Firebase app initialization failed: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("✅ Firebase app initialized")
	
	// Initialize Firestore client
	fmt.Println("🗄️  Initializing Firestore client...")
	client, err := app.Firestore(ctx)
	if err != nil {
		fmt.Printf("❌ Firestore client initialization failed: %v\n", err)
		os.Exit(1)
	}
	defer client.Close()
	fmt.Println("✅ Firestore client initialized")
	
	// Test write operation
	fmt.Println("✍️  Testing write operation...")
	testDoc := map[string]interface{}{
		"timestamp": time.Now().Unix(),
		"test":      "connection_test",
		"status":    "testing",
		"server":    "production",
	}
	
	_, err = client.Collection("connection_tests").Doc("test_" + fmt.Sprintf("%d", time.Now().Unix())).Set(ctx, testDoc)
	if err != nil {
		fmt.Printf("❌ Firestore write test failed: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("✅ Write operation successful")
	
	// Test read operation
	fmt.Println("📖 Testing read operation...")
	docs, err := client.Collection("connection_tests").Limit(1).Documents(ctx).GetAll()
	if err != nil {
		fmt.Printf("❌ Firestore read test failed: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("✅ Read operation successful - found %d documents\n", len(docs))
	
	fmt.Println("\n🎉 Firebase connection test PASSED!")
	fmt.Println("✅ Firebase app: Working")
	fmt.Println("✅ Firestore client: Working") 
	fmt.Println("✅ Write operations: Working")
	fmt.Println("✅ Read operations: Working")
}
EOF

# Check if go.mod exists (for source directory) or use temporary module
if [ -f "go.mod" ]; then
    echo "📦 Using existing go.mod..."
    go run test_firebase_connection.go
else
    echo "📦 Creating temporary module..."
    # Create temporary module for standalone testing
    cat > go.mod << 'EOF'
module firebase-test

go 1.22

require (
	cloud.google.com/go/firestore v1.14.0
	firebase.google.com/go/v4 v4.12.0
	google.golang.org/api v0.149.0
)
EOF
    
    echo "⬇️  Downloading dependencies..."
    go mod download
    go run test_firebase_connection.go
    
    # Cleanup temporary files
    rm -f go.mod go.sum
fi

# Cleanup test file
rm -f test_firebase_connection.go

echo "🧹 Cleanup complete"
