package mygopherpostgres

import (
	"context"
	"fmt"
	"log"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// ConnectToPostgreSQLGormDB establishes a connection to the PostgreSQL database using GORM,
// with automatic retry logic and context-based timeout handling. It ensures that the
// `uuid-ossp` extension is enabled in the database upon successful connection.
//
// Parameters:
//   - ctx: A context to control the connection's cancellation and timeout.
//   - dsn: The Data Source Name (DSN), typically the database connection string.
//   - timeout: The total duration allowed for the connection attempts before timing out.
//   - maxRetries: The maximum number of connection attempts in case of failure.
//
// Returns:
//   - *gorm.DB: A pointer to the GORM DB instance if the connection is successful.
//   - error: An error describing the failure if the connection cannot be established
//     within the given number of retries.
//
// Example usage:
//
//	ctx := context.Background()
//	dsn := "postgres://username:password@localhost:5432/dbname?sslmode=disable"
//	timeout := 30 * time.Second
//	maxRetries := 3
//
//	db, err := ConnectToPostgreSQLGormDB(ctx, dsn, timeout, maxRetries)
//	if err != nil {
//	    log.Fatalf("Error connecting to the database: %v", err)
//	}
func ConnectToPostgreSQLGormDB(ctx context.Context, dsn string, timeout time.Duration, maxRetries int) (*gorm.DB, error) {
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
			db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
			if err == nil {
				// Successfully connected, enable the uuid-ossp extension if necessary
				log.Println("Connected to PostgreSQL using GORM successfully")
				return db, nil // Return the connected DB instance
			}

			// Log the failure and retry after a delay
			log.Printf("Connection attempt %d failed: %v", i+1, err)
			time.Sleep(retryDelay) // Wait before the next retry
		}
	}

	// Return error if all retries fail
	return nil, fmt.Errorf("failed to connect to PostgreSQL after %d retries: %w", maxRetries, err)
}
