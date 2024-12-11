package gopherlogger

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"sync"
	"time"
)

// LogLevel represents the severity of a log entry, providing a way to categorize and filter log messages
// based on their importance and urgency.
//
// The levels are hierarchical, allowing filtering of less critical messages:
// - DEBUG: Lowest level, used for detailed tracing and diagnostics
// - INFO: General information about application flow
// - WARN: Potential issues that don't prevent application from functioning
// - ERROR: Significant problems that may impact functionality
// - FATAL: Critical errors that will cause application termination
//
// Example usage:
//
//	logger, _ := NewGopherLogger(LoggerConfig{LogLevel: INFO})
//	// This logger will only output INFO, WARN, ERROR, and FATAL messages
//	logger.Debug("This won't be logged")
//	logger.Info("Application started successfully")
type LogLevel int

const (
	DEBUG LogLevel = iota
	INFO
	WARN
	ERROR
	FATAL
)

// String returns a human-readable string representation of the LogLevel.
//
// This method allows easy conversion of LogLevel to its corresponding string name,
// which is useful for logging and debugging purposes.
//
// Example:
//
//	level := DEBUG
//	fmt.Println(level.String())  // Outputs: "DEBUG"
func (l LogLevel) String() string {
	return [...]string{"DEBUG", "INFO", "WARN", "ERROR", "FATAL"}[l]
}

// ParseLogLevel converts a string representation of a log level into the corresponding LogLevel type.
// It returns an error if the string does not correspond to a valid LogLevel.
func ParseLogLevel(level string) (LogLevel, error) {
	switch strings.ToUpper(level) {
	case "DEBUG":
		return DEBUG, nil
	case "INFO":
		return INFO, nil
	case "WARN":
		return WARN, nil
	case "ERROR":
		return ERROR, nil
	case "FATAL":
		return FATAL, nil
	default:
		return DEBUG, errors.New("invalid log level: " + level)
	}
}

// Logger is a sophisticated logging structure providing advanced file rotation and
// multi-destination logging capabilities. It supports:
// - Configurable log levels
// - File size-based log rotation
// - Concurrent-safe logging
// - Buffered logging for performance
// - Output to both file and terminal
//
// The logger manages log files efficiently, creating new log files when size thresholds are reached
// and maintaining a configurable number of backup log files.
type Logger struct {
	mu           sync.RWMutex
	level        LogLevel
	baseFilename string
	maxBytes     int64
	backupCount  int
	file         *os.File
	multiWriter  io.Writer
	bufferPool   *sync.Pool
}

// LoggerConfig provides a flexible configuration mechanism for customizing logger behavior.
//
// Fields allow precise control over logging characteristics:
// - Filename: Base name for log files
// - MaxBytes: Maximum file size before rotation
// - BackupCount: Number of historical log files to retain
// - LogLevel: Minimum severity of logs to record
//
// Example:
//
//	config := LoggerConfig{
//	    Filename:    "application.log",
//	    MaxBytes:    5 * 1024 * 1024,  // 5 MB
//	    BackupCount: 3,
//	    LogLevel:    INFO,
//	}
type LoggerConfig struct {
	Filename    string
	MaxBytes    int64
	BackupCount int
	LogLevel    LogLevel
}

// NewGopherLogger creates a new rotating file logger with terminal output and advanced configuration.
//
// This function:
// - Creates logs directory if it doesn't exist
// - Generates timestamped log filenames
// - Sets up file and console logging
// - Configures log rotation parameters
//
// Parameters:
//   - config: Configuration settings for the logger
//
// Returns:
//   - *Logger: Configured logger instance
//   - error: Any error encountered during logger setup
//
// Example:
//
//	logger, err := NewGopherLogger(LoggerConfig{
//	    Filename:    "app.log",
//	    MaxBytes:    10*1024*1024,
//	    BackupCount: 5,
//	    LogLevel:    DEBUG,
//	})
//	if err != nil {
//	    panic(err)
//	}
//	defer logger.Close()
func NewGopherLogger(config LoggerConfig) (*Logger, error) {
	// Ensure logs directory exists in current working directory
	logsDir := "logs"
	if err := os.MkdirAll(logsDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create logs directory: %v", err)
	}

	// Generate log filename with unix timestamp
	if config.Filename == "" {
		config.Filename = "gopherlogger.log"
	}
	logFilePath := filepath.Join(logsDir, config.Filename)

	// Default configuration if not specified
	if config.MaxBytes == 0 {
		config.MaxBytes = 10 * 1024 * 1024 // 10 MB
	}
	if config.BackupCount == 0 {
		config.BackupCount = 5
	}

	// Open the log file
	file, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to open log file: %v", err)
	}

	// Create multi-writer to output to both terminal and file
	multiWriter := io.MultiWriter(os.Stdout, file)

	// Create logger with buffer pool for efficiency
	logger := &Logger{
		baseFilename: logFilePath,
		maxBytes:     config.MaxBytes,
		backupCount:  config.BackupCount,
		file:         file,
		level:        config.LogLevel,
		multiWriter:  multiWriter,
		bufferPool: &sync.Pool{
			New: func() interface{} {
				return new(bytes.Buffer)
			},
		},
	}

	return logger, nil
}

// formatLogMessage creates a comprehensive, high-performance log message format.
//
// Key features:
// - Includes timestamp
// - Shows log level
// - Captures source file and line number
// - Includes calling function name
//
// Designed to provide maximum context with minimal performance overhead.
//
// Parameters:
//   - level: Severity of the log message
//   - message: Actual log content
//
// # Returns formatted log message as a string
//
// Example:
//
//	formattedMsg := logger.formatLogMessage(INFO, "User logged in")
//	// Outputs something like: "2024/01/15 14:30:45.123456 [INFO] main.go:42 (main): User logged in"
func (l *Logger) formatLogMessage(level LogLevel, message string) string {
	// Get caller information with minimal overhead
	pc, file, line, ok := runtime.Caller(2)
	if !ok {
		file = "unknown"
		line = 0
	}

	// Use runtime.FuncForPC to get function name (more efficient than filepath)
	funcName := "unknown"
	if fn := runtime.FuncForPC(pc); fn != nil {
		funcName = filepath.Base(fn.Name())
	}

	// Trim file path to filename
	file = filepath.Base(file)

	// Format timestamp
	now := time.Now()
	return fmt.Sprintf("%s [%s] %s:%d (%s): %s\n",
		now.Format("2006/01/02 15:04:05.000000"),
		level,
		file,
		line,
		funcName,
		message,
	)
}

// log is an internal method handling logging with high performance and thread safety.
//
// Core responsibilities:
// - Level-based filtering
// - Message formatting
// - Efficient buffer management
// - Thread-safe writing
// - Automatic log rotation
//
// Designed for minimal allocation and maximum concurrency support.
//
// Parameters:
//   - level: Log message severity
//   - message: Content to be logged
//
// Does not return any value, writes directly to configured outputs
func (l *Logger) log(level LogLevel, message string) {
	// Quick level check
	if level < l.level {
		return
	}

	// Format message
	formattedMsg := l.formatLogMessage(level, message)

	// Use buffer from pool to reduce allocations
	buf := l.bufferPool.Get().(*bytes.Buffer)
	defer l.bufferPool.Put(buf)
	buf.Reset()
	buf.WriteString(formattedMsg)

	// Thread-safe writing
	l.mu.Lock()
	defer l.mu.Unlock()

	// Check file size and rotate if needed
	if err := l.checkFileSize(); err != nil {
		// Log rotation error (consider using a lightweight error tracking)
		fmt.Fprintf(os.Stderr, "Log rotation error: %v\n", err)
	}

	// Write to multi-writer
	l.multiWriter.Write(buf.Bytes())
}

// checkFileSize determines if log file rotation is necessary.
//
// Compares current log file size against configured maximum size threshold.
// Triggers log rotation if size limit is exceeded.
//
// Returns:
//   - nil if file size is acceptable
//   - error if file stat or rotation fails
func (l *Logger) checkFileSize() error {
	fi, err := l.file.Stat()
	if err != nil {
		return err
	}

	// If file size exceeds max bytes, rotate
	if fi.Size() >= l.maxBytes {
		return l.rotateLogFiles()
	}

	return nil
}

// rotateLogFiles manages log file rotation process.
//
// Key rotation steps:
// - Close current log file
// - Rename current log file with timestamp
// - Remove excess backup files
// - Open new log file
// - Update multi-writer
//
// Ensures log history is maintained while preventing unbounded disk usage.
//
// Returns:
//   - nil on successful rotation
//   - error if any rotation step fails
func (l *Logger) rotateLogFiles() error {
	// Close current file
	if err := l.file.Close(); err != nil {
		return err
	}

	// Generate new log filename
	currentTime := time.Now().Format("20060102_150405")
	newLogFilePath := filepath.Join(filepath.Dir(l.baseFilename),
		fmt.Sprintf("%s_%s", currentTime, filepath.Base(l.baseFilename)))

	// Rename current log file
	if err := os.Rename(l.baseFilename, newLogFilePath); err != nil {
		return err
	}

	// Manage backup files
	backupFiles, _ := filepath.Glob(l.baseFilename + ".*")
	if len(backupFiles) > l.backupCount {
		// Sort and remove oldest backup files
		sort.Strings(backupFiles)
		for _, oldFile := range backupFiles[:len(backupFiles)-l.backupCount] {
			os.Remove(oldFile)
		}
	}

	// Reopen the base log file
	file, err := os.OpenFile(l.baseFilename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return err
	}
	l.file = file

	// Update multi-writer
	l.multiWriter = io.MultiWriter(os.Stdout, file)

	return nil
}

// Debug logs a message at DEBUG level with minimal arguments.
//
// Suitable for detailed tracing and diagnostic information.
// Only logged if logger's level is set to DEBUG or lower.
//
// Example:
//
//	logger.Debug("Function entry", variableName)
func (l *Logger) Debug(v ...interface{}) {
	l.log(DEBUG, fmt.Sprint(v...))
}

func (l *Logger) Info(v ...interface{}) {
	l.log(INFO, fmt.Sprint(v...))
}

func (l *Logger) Warn(v ...interface{}) {
	l.log(WARN, fmt.Sprint(v...))
}

func (l *Logger) Error(v ...interface{}) {
	l.log(ERROR, fmt.Sprint(v...))
}

func (l *Logger) Fatal(v ...interface{}) {
	l.log(FATAL, fmt.Sprint(v...))
	os.Exit(1)
}

// Formatted logging methods
func (l *Logger) Debugf(format string, v ...interface{}) {
	l.log(DEBUG, fmt.Sprintf(format, v...))
}

func (l *Logger) Infof(format string, v ...interface{}) {
	l.log(INFO, fmt.Sprintf(format, v...))
}

func (l *Logger) Warnf(format string, v ...interface{}) {
	l.log(WARN, fmt.Sprintf(format, v...))
}

func (l *Logger) Errorf(format string, v ...interface{}) {
	l.log(ERROR, fmt.Sprintf(format, v...))
}

func (l *Logger) Fatalf(format string, v ...interface{}) {
	l.log(FATAL, fmt.Sprintf(format, v...))
	os.Exit(1)
}

// Close safely closes the log file and releases associated resources.
//
// Should be called when logger is no longer needed, typically via defer.
//
// Example:
//
//	logger, _ := NewGopherLogger(config)
//	defer logger.Close()  // Ensures log file is properly closed
func (l *Logger) Close() error {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.file.Close()
}
