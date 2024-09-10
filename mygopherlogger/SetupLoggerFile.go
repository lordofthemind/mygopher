package mygopherlogger

import (
	"fmt"
	"io"
	"log"
	"os"
	"time"
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
	// Ensure the logs directory exists
	if err := os.MkdirAll("logs", 0755); err != nil {
		return nil, fmt.Errorf("failed to create logs directory: %w", err)
	}

	// Get the current date and time for the log file prefix
	currentTime := time.Now().Format("20060102_150405")
	logFileName = fmt.Sprintf("%s_%s", currentTime, logFileName)

	// Open the log file
	logFilePath := "logs/" + logFileName
	logFile, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		// If there's an error opening the log file, fallback to stdout only logging
		log.Printf("Error opening log file %s, falling back to stdout only: %v", logFilePath, err)
		logFile = nil // Continue without a log file
	} else {
		// Successfully opened the log file
		log.Printf("Logging initialized. Log file: %s", logFilePath)
	}

	// Set up multi-writer to write to stdout and file if possible
	var multiWriter io.Writer
	if logFile != nil {
		multiWriter = io.MultiWriter(os.Stdout, logFile)
	} else {
		multiWriter = os.Stdout
	}

	// Configure log format to include timestamp, log level, and file location
	log.SetOutput(multiWriter)
	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds | log.Lshortfile)

	return logFile, nil
}
