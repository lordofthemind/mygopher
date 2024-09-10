package mygopherpostgres

import (
	"context"
	"fmt"
	"log"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

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
