package gopherpostgres

import (
	"context"
	"fmt"
	"log"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// ConnectToPostgresGORM connects to a PostgreSQL database using GORM with retries.
//
// This function attempts to connect to a PostgreSQL database using the GORM ORM library.
// It applies the provided DSN and tries to connect up to 'maxRetries' times. The function
// also uses a context with a timeout to control how long the connection attempt lasts.
// If successful, a *gorm.DB instance is returned, which allows performing ORM-based
// operations.
//
// Params:
//
//	ctx - The context for managing connection timeout and cancellation.
//	dsn - The PostgreSQL connection string (Data Source Name).
//	timeout - The timeout duration for the connection attempt.
//	maxRetries - The maximum number of retries before giving up.
//
// Returns:
//
//	*gorm.DB - The connected GORM PostgreSQL database instance on success.
//	error - An error message if the connection fails after the retries.
//
// Example usage:
//
//	ctx := context.Background()
//	db, err := ConnectToPostgresGORM(ctx, "postgres://user:password@localhost:5432/mydb", 10*time.Second, 3)
//	if err != nil {
//	    log.Fatalf("Failed to connect to PostgreSQL using GORM: %v", err)
//	}
//
// Once connected, you can use GORM's ORM features for database operations like querying,
// inserting, updating, and deleting records.
func ConnectToPostgresGORM(ctx context.Context, dsn string, timeout time.Duration, maxRetries int) (*gorm.DB, error) {
	// Set a timeout for the connection operation using the context
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	// Validate the DSN (database URL) input
	if dsn == "" {
		return nil, fmt.Errorf("missing required database URL (DSN)")
	}

	var db *gorm.DB
	var err error
	retryDelay := 5 * time.Second // Time to wait between retries

	// Attempt to connect with retries
	for i := 0; i < maxRetries; i++ {
		select {
		case <-ctx.Done():
			// If context times out or is canceled, exit with an error
			return nil, fmt.Errorf("context timed out while trying to connect to database: %w", ctx.Err())
		default:
			// Try to open the connection using GORM
			log.Printf("Attempting to connect to PostgreSQL using GORM... (Attempt %d of %d)", i+1, maxRetries)
			db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
			if err == nil {
				// Successfully connected
				log.Println("Connected to PostgreSQL using GORM successfully")
				return db, nil // Return the connected DB instance
			}

			// Log the failure and retry after a delay
			log.Printf("Connection attempt %d failed: %v", i+1, err)
			log.Printf("Retrying connection in %v seconds...", retryDelay.Seconds())
			time.Sleep(retryDelay) // Wait before the next retry
		}
	}

	// Log final failure before exiting
	log.Fatalf("Failed to connect to PostgreSQL using GORM after %d attempts: %v", maxRetries, err)
	return nil, fmt.Errorf("failed to connect to PostgreSQL after %d retries: %w", maxRetries, err)
}

// package main

// import (
// 	"context"
// 	"log"
// 	"time"

// 	"github.com/lordofthemind/mygopher/gopherpostgres"
// )

// func main() {
// 	ctx := context.Background()
// 	db, err := gopherpostgres.ConnectToPostgresGORM(ctx, "postgres://user:password@localhost:5432/mydb", 10*time.Second, 3)
// 	if err != nil {
// 		// This log will not be hit because ConnectToPostgresGORM exits the application on failure.
// 		log.Fatalf("Unable to continue: %v", err)
// 	}
// 	defer db.Close()

// 	// Continue with your application logic...
// }
