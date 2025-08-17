#!/bin/bash

set -e

echo "🔥 Testing Firebase Connection..."

# Check if we're in the right directory
if [ ! -f "firebaseAdmin.json" ]; then
    echo "❌ firebaseAdmin.json not found. Make sure you're in the Sauron directory."
    exit 1
fi

# Create temporary directory for Firebase testing
TEMP_DIR=$(mktemp -d)
echo "📦 Creating temporary workspace: $TEMP_DIR"

# Copy required files to temp directory
cp firebaseAdmin.json "$TEMP_DIR/"
cd "$TEMP_DIR"

# Create the Go test file in temp directory
cat > test_firebase_connection.go << 'EOF'
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	firebase "firebase.google.com/go/v4"
	"google.golang.org/api/option"
)

type Credentials struct {
	ProjectID string `json:"project_id"`
}

func main() {
	fmt.Println("🔥 Testing Firebase connection...")
	
	// Check if credentials file exists
	if _, err := os.Stat("firebaseAdmin.json"); os.IsNotExist(err) {
		fmt.Println("❌ firebaseAdmin.json not found")
		os.Exit(1)
	}
	
	// Read and parse credentials file to get project ID
	credData, err := ioutil.ReadFile("firebaseAdmin.json")
	if err != nil {
		fmt.Printf("❌ Failed to read credentials file: %v\n", err)
		os.Exit(1)
	}

	var creds Credentials
	if err := json.Unmarshal(credData, &creds); err != nil {
		fmt.Printf("❌ Failed to parse credentials file: %v\n", err)
		os.Exit(1)
	}

	if creds.ProjectID == "" {
		fmt.Println("❌ No project_id found in credentials file")
		os.Exit(1)
	}

	fmt.Printf("📋 Project ID: %s\n", creds.ProjectID)
	
	ctx := context.Background()
	opt := option.WithCredentialsFile("firebaseAdmin.json")
	
	// Initialize Firebase app with project ID
	fmt.Println("📡 Initializing Firebase app...")
	config := &firebase.Config{ProjectID: creds.ProjectID}
	app, err := firebase.NewApp(ctx, config, opt)
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

# Create temporary module for standalone testing
echo "⬇️  Creating temporary Go module..."
cat > go.mod << 'EOF'
module firebase-test

go 1.22
EOF

echo "⬇️  Downloading Firebase dependencies..."
go mod init firebase-test 2>/dev/null || true
go get firebase.google.com/go/v4@latest
go get google.golang.org/api@latest
go mod tidy

echo "🚀 Running Firebase connection test..."
go run test_firebase_connection.go

# Return to original directory and cleanup
cd - > /dev/null
echo "🧹 Cleaning up temporary files..."
rm -rf "$TEMP_DIR"

# Cleanup test file
rm -f test_firebase_connection.go

echo "🧹 Cleanup complete"
