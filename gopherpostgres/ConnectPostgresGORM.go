package gopherpostgres

import (
	"context"
	"fmt"
	"log"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// PostgresConfig encapsulates the configuration parameters for connecting to a PostgreSQL database.
//
// Fields:
//
// - DSN: PostgreSQL connection string (Data Source Name).
// - Timeout: Timeout duration for the connection attempt.
// - MaxRetries: Maximum number of retries before giving up.
//
// Example:
//
//	config := PostgresConfig{
//	    DSN:        "postgres://user:password@localhost:5432/mydb",
//	    Timeout:    10 * time.Second,
//	    MaxRetries: 3,
//	}
type PostgresConfig struct {
	DSN        string        // PostgreSQL connection string (Data Source Name)
	Timeout    time.Duration // Timeout duration for the connection attempt
	MaxRetries int           // Maximum number of retries before giving up
}

// ConnectToPostgresGORM connects to a PostgreSQL database using GORM with retries.
//
// This function attempts to connect to a PostgreSQL database using the GORM ORM library.
// It applies the provided configuration and tries to connect up to 'MaxRetries' times.
// The function uses a context with a timeout to control how long the connection attempt lasts.
// If successful, a *gorm.DB instance is returned, which allows performing ORM-based operations.
//
// Params:
//
//	ctx - The context for managing connection timeout and cancellation.
//	config - The configuration struct containing DSN, Timeout, and MaxRetries.
//
// Returns:
//
//	*gorm.DB - The connected GORM PostgreSQL database instance on success.
//	error - An error message if the connection fails after the retries.
//
// Example:
//
//	config := gopherpostgres.PostgresConfig{
//	    DSN:        "postgres://user:password@localhost:5432/mydb",
//	    Timeout:    10 * time.Second,
//	    MaxRetries: 3,
//	}
//
//	ctx := context.Background()
//	db, err := ConnectToPostgresGORM(ctx, config)
//	if err != nil {
//	    log.Fatalf("Failed to connect to PostgreSQL using GORM: %v", err)
//	}
func ConnectToPostgresGORM(ctx context.Context, config PostgresConfig) (*gorm.DB, error) {
	// Validate the DSN (database URL) input
	if config.DSN == "" {
		return nil, fmt.Errorf("missing required database URL (DSN)")
	}

	// Set a timeout for the connection operation using the context
	ctx, cancel := context.WithTimeout(ctx, config.Timeout)
	defer cancel()

	var db *gorm.DB
	var err error
	retryDelay := 5 * time.Second // Time to wait between retries

	// Attempt to connect with retries
	for i := 0; i < config.MaxRetries; i++ {
		select {
		case <-ctx.Done():
			// If context times out or is canceled, exit with an error
			return nil, fmt.Errorf("context timed out while trying to connect to database: %w", ctx.Err())
		default:
			// Try to open the connection using GORM
			log.Printf("Attempting to connect to PostgreSQL using GORM... (Attempt %d of %d)", i+1, config.MaxRetries)
			db, err = gorm.Open(postgres.Open(config.DSN), &gorm.Config{})
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
	log.Fatalf("Failed to connect to PostgreSQL using GORM after %d attempts: %v", config.MaxRetries, err)
	return nil, fmt.Errorf("failed to connect to PostgreSQL after %d retries: %w", config.MaxRetries, err)
}
