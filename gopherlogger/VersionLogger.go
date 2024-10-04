package gopherlogger

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"
)

// SetUpLoggerFile sets up logging to both a file and stdout.
//
// This function creates a log file in the "logs" directory with the provided log file name,
// prefixed by the current date and time. If the log file cannot be opened, the function falls
// back to logging to stdout. It configures the log format to include timestamps, log levels,
// and the file location of the log entry. The logs are written to both the log file and stdout.
//
// Params:
//
//	logFileName - The base name for the log file (e.g., "app.log").
//
// Returns:
//
//	*os.File - A pointer to the log file if successfully opened, or nil if it falls back to stdout.
//	error    - An error message if the log file could not be opened.
//
// Example usage:
//
//	logFile, err := SetUpLoggerFile("app.log")
//	if err != nil {
//	    log.Fatalf("Failed to initialize logger: %v", err)
//	}
//	defer logFile.Close()
//
// In this example, the logger writes to both stdout and a file named with the current timestamp.
func SetUpLoggerFile(logFileName string) (*os.File, error) {
	// Ensure the logs directory exists
	if err := os.MkdirAll("logs", 0755); err != nil {
		return nil, fmt.Errorf("failed to create logs directory: %w", err)
	}

	// Get the current date and time for the log file prefix
	currentTime := time.Now().Format("20060102_150405")
	logFilePath := filepath.Join("logs", fmt.Sprintf("%s_%s", currentTime, logFileName))

	// Open the log file
	logFile, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Printf("Error opening log file %s, using stdout: %v", logFilePath, err)
		log.SetOutput(os.Stdout)
		return nil, fmt.Errorf("failed to open log file: %w", err)
	}

	// Set up multi-writer to write to both stdout and log file
	multiWriter := io.MultiWriter(os.Stdout, logFile)

	// Configure log format to include timestamp, log level, and file location
	log.SetOutput(multiWriter)
	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds | log.Lshortfile)

	log.Printf("Logging initialized. Log file: %s", logFilePath)
	return logFile, nil
}
