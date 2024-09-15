package gopherpostgres

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq"
)

// ConnectPostgresDB connects to a PostgreSQL database using the sql package with retries.
//
// This function attempts to connect to a PostgreSQL database using the provided Data Source Name (DSN),
// retrying the connection up to 'maxRetries' times. It uses a context with a timeout to ensure
// that the connection does not hang indefinitely. If the connection is successful, it returns
// a *sql.DB instance for database operations.
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
//	*sql.DB - The connected PostgreSQL database instance on success.
//	error - An error message if the connection fails after the retries.
//
// Example usage:
//
//	ctx := context.Background()
//	db, err := ConnectPostgresDB(ctx, "postgres://user:password@localhost:5432/mydb", 10*time.Second, 3)
//	if err != nil {
//	    log.Fatalf("Failed to connect to PostgreSQL: %v", err)
//	}
//	defer db.Close()
//
// Once connected, you can perform SQL operations like querying or executing statements.
func ConnectPostgresDB(ctx context.Context, dsn string, timeout time.Duration, maxRetries int) (*sql.DB, error) {
	// Set a timeout for the connection operation using the context
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	// Validate the DSN (database URL) input
	if dsn == "" {
		return nil, fmt.Errorf("missing required database URL (DSN)")
	}

	var db *sql.DB
	var err error
	retryDelay := 5 * time.Second // Time to wait between retries

	// Attempt to connect with retries
	for i := 0; i < maxRetries; i++ {
		select {
		case <-ctx.Done():
			// If context times out or is canceled, exit with an error
			return nil, fmt.Errorf("context timed out while trying to connect to database: %w", ctx.Err())
		default:
			// Try to open the connection using the standard library's sql package
			db, err = sql.Open("postgres", dsn)
			if err == nil {
				// Ping the database to ensure connection is established
				err = db.PingContext(ctx)
				if err == nil {
					log.Println("Connected to PostgreSQL successfully")
					return db, nil // Return the connected DB instance
				}
			}

			// Log the failure and retry after a delay
			log.Printf("Connection attempt %d failed: %v", i+1, err)
			time.Sleep(retryDelay) // Wait before the next retry
		}
	}

	// Return error if all retries fail
	return nil, fmt.Errorf("failed to connect to PostgreSQL after %d retries: %w", maxRetries, err)
}
