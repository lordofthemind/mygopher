package mygopherlogger

import (
	"bytes"
	"log"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// TestSetUpLoggerFile tests the SetUpLoggerFile function for creating logs directory,
// writing logs to a file and stdout, and handling file opening errors gracefully.
func TestSetUpLoggerFile(t *testing.T) {
	// Set up
	logFileName := "test_app.log"
	logDir := "logs"

	// Cleanup before and after the test
	_ = os.RemoveAll(logDir)
	defer os.RemoveAll(logDir)

	// Create a temporary buffer to capture log output
	var logBuffer bytes.Buffer
	log.SetOutput(&logBuffer)

	// Call the function
	logFile, err := SetUpLoggerFile(logFileName)
	if err != nil {
		t.Fatalf("SetUpLoggerFile returned an unexpected error: %v", err)
	}

	// Check if the logs directory was created
	if _, err := os.Stat(logDir); os.IsNotExist(err) {
		t.Fatalf("Expected logs directory to be created, but it does not exist: %v", err)
	}

	// Check if the log file is created with the correct naming convention
	expectedLogFilePrefix := time.Now().Format("20060102_150405")
	expectedLogFileName := filepath.Join(logDir, expectedLogFilePrefix+"_"+logFileName)

	if logFile != nil {
		if _, err := os.Stat(expectedLogFileName); os.IsNotExist(err) {
			t.Fatalf("Expected log file %s to be created, but it does not exist: %v", expectedLogFileName, err)
		}
		defer logFile.Close()
	} else {
		t.Logf("Log file is nil, using stdout only")
	}

	// Write test log messages
	log.Println("Test log message")

	// Check that the log messages are written to stdout or log file
	if logFile != nil {
		// Verify content written to log file
		content, err := os.ReadFile(expectedLogFileName)
		if err != nil {
			t.Fatalf("Failed to read log file: %v", err)
		}
		if !strings.Contains(string(content), "Test log message") {
			t.Fatalf("Expected log message not found in log file")
		}
	} else {
		// Verify content written to stdout (captured in logBuffer)
		if !strings.Contains(logBuffer.String(), "Test log message") {
			t.Fatalf("Expected log message not found in stdout")
		}
	}

	// Simulate error in opening file
	// Temporarily set the function to try opening a directory as a file
	os.Mkdir("logs/fail.log", 0755) // Create a directory with the log file name to force an error
	defer os.Remove("logs/fail.log")

	_, err = SetUpLoggerFile("fail.log")
	if err == nil {
		t.Fatalf("Expected error when trying to log to a directory, but got none")
	}
	if !strings.Contains(logBuffer.String(), "falling back to stdout only") {
		t.Fatalf("Expected fallback to stdout only logging, but it didn't occur")
	}
}
