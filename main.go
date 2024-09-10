package mygopher

import (
	"context"
	"database/sql"
	"os"
	"time"

	"github.com/lordofthemind/mygopher/mygopherlogger"
	"github.com/lordofthemind/mygopher/mygophermongodb"
	"github.com/lordofthemind/mygopher/mygopherpostgres"
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
)

// SetUpLoggerFile sets up logging to both a file and stdout.
//
// It creates a "logs" directory if it doesn't already exist. It generates a log file name
// based on the current date and time and opens this file for appending. The function configures
// logging to output to both the file and stdout. If there is an issue opening the log file,
// it falls back to logging only to stdout.
//
// Parameters:
//   - logFileName: The base name for the log file (e.g., "app.log"). The final log file name
//     will include a timestamp prefix to ensure uniqueness.
//
// Returns:
//   - *os.File: A pointer to the opened log file which should be closed by the caller when
//     logging is complete or when the application shuts down.
//   - error: An error is returned if the function fails to create the log directory or open
//     the log file. If an error occurs while opening the log file, it falls back to stdout
//     logging and logs the error.
//
// Example usage:
//
//	logFile, err := logger.SetUpLoggerFile("app.log")
//	if err != nil {
//	    log.Fatalf("Failed to set up logger: %v", err)
//	}
//	defer logFile.Close() // Ensure to close the log file when the application exits
func SetUpLoggerFile(logFileName string) (*os.File, error) {
	logFile, err := mygopherlogger.SetUpLoggerFile(logFileName)
	return logFile, err
}

// ConnectToMongoDB establishes a connection to the MongoDB database using the official MongoDB Go driver,
// with retry and context-based timeout handling. It returns a client and database instance
// that can be used for subsequent MongoDB operations.
//
// Parameters:
//   - ctx: A context to control the connection's cancellation and timeout.
//   - dsn: The MongoDB connection string (Data Source Name).
//   - timeout: The total duration allowed for the connection attempts before timing out.
//   - maxRetries: The maximum number of connection attempts in case of failure.
//   - dbName: The name of the MongoDB database to connect to.
//
// Returns:
//   - *mongo.Client: A pointer to the MongoDB client instance if the connection is successful.
//   - *mongo.Database: A pointer to the specific MongoDB database instance.
//   - error: An error describing the failure if the connection cannot be established
//     within the given number of retries.
//
// Example usage:
//
//	ctx := context.Background()
//	dsn := os.Getenv("MONGO_DOCKER_CONNECTION_URL")
//	timeout := 30 * time.Second
//	maxRetries := 5
//	dbName := "polyglot" // Replace with your actual database name
//
//	client, db, err := ConnectToMongoDB(ctx, dsn, timeout, maxRetries, dbName)
//	if err != nil {
//	    log.Fatalf("Error connecting to MongoDB: %v", err)
//	}
func ConnectToMongoDB(ctx context.Context, dsn string, timeout time.Duration, maxRetries int, dbName string) (*mongo.Client, *mongo.Database, error) {
	mongoClient, mongoDatabase, err := mygophermongodb.ConnectToMongoDB(ctx, dsn, timeout, maxRetries, dbName)
	return mongoClient, mongoDatabase, err
}

// ConnectPostgresDB establishes a connection to the PostgreSQL database using the `database/sql` package,
// with retry and context-based timeout handling. It provides a connection pool that can be used
// for database operations throughout the application's lifetime.
//
// Parameters:
//   - ctx: A context to control the connection's cancellation and timeout.
//   - dsn: The Data Source Name (DSN), typically the database connection string.
//   - timeout: The total duration allowed for the connection attempts before timing out.
//   - maxRetries: The maximum number of connection attempts in case of failure.
//
// Returns:
//   - *sql.DB: A pointer to the database connection pool if the connection is successful.
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
//	db, err := ConnectPostgresDB(ctx, dsn, timeout, maxRetries)
//	if err != nil {
//	    log.Fatalf("Error connecting to the database: %v", err)
//	}
//	defer db.Close() // Always ensure to close the database connection when done
func ConnectPostgresDB(ctx context.Context, dsn string, timeout time.Duration, maxRetries int) (*sql.DB, error) {
	SQLdb, err := mygopherpostgres.ConnectPostgresDB(ctx, dsn, timeout, maxRetries)
	return SQLdb, err
}

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
	GORMdb, err := mygopherpostgres.ConnectToPostgreSQLGormDB(ctx, dsn, timeout, maxRetries)
	return GORMdb, err
}
