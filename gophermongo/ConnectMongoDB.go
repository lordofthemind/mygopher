package gophermongo

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// ConnectToMongoDB establishes a connection to MongoDB with retries and a context timeout.
//
// This function attempts to connect to MongoDB using the provided connection string (DSN),
// retrying the connection up to 'maxRetries' times with a delay of 5 seconds between retries.
// It also applies a timeout to the entire connection attempt using the context.
//
// Params:
//
//	ctx - The context for connection management (with timeout support).
//	dsn - The MongoDB connection string (Data Source Name).
//	timeout - The timeout duration for the connection attempt.
//	maxRetries - The maximum number of retries before giving up.
//
// Returns:
//
//	*mongo.Client - The connected MongoDB client instance on success.
//	error - An error message if the connection fails.
//
// Example usage:
//
//	ctx := context.Background()
//	client, err := ConnectToMongoDB(ctx, "mongodb://localhost:27017", 10*time.Second, 3)
//	if err != nil {
//	    log.Fatalf("Failed to connect to MongoDB: %v", err)
//	}
//	defer client.Disconnect(ctx)
func ConnectToMongoDB(ctx context.Context, dsn string, timeout time.Duration, maxRetries int) (*mongo.Client, error) {
	// Set a timeout for the connection operation using the context
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	// Validate the DSN (connection string) input
	if dsn == "" {
		return nil, fmt.Errorf("missing required MongoDB connection string (DSN)")
	}

	var client *mongo.Client
	var err error
	retryDelay := 5 * time.Second // Time to wait between retries

	// Attempt to connect with retries
	for i := 0; i < maxRetries; i++ {
		select {
		case <-ctx.Done():
			// If context times out or is canceled, exit with an error
			return nil, fmt.Errorf("context timed out while trying to connect to MongoDB: %w", ctx.Err())
		default:
			// Try to establish a connection to MongoDB
			client, err = mongo.Connect(ctx, options.Client().ApplyURI(dsn))
			if err == nil {
				// Successfully connected, verify the connection
				if err = client.Ping(ctx, nil); err != nil {
					// If ping fails, return an error
					return nil, fmt.Errorf("failed to ping MongoDB: %w", err)
				}
				log.Println("Connected to MongoDB successfully")
				return client, nil
			}

			// Log the failure and retry after a delay
			log.Printf("Connection attempt %d failed: %v\n", i+1, err)
			time.Sleep(retryDelay) // Wait before the next retry
		}
	}

	// Return error if all retries fail
	return nil, fmt.Errorf("failed to connect to MongoDB after %d retries: %w", maxRetries, err)
}
