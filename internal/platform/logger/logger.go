package logger

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

// LogLevel represents the severity of a log message
type LogLevel int

const (
	DEBUG LogLevel = iota
	INFO
	WARN
	ERROR
	FATAL
)

/*
String converts a LogLevel enum value to its string representation.
Returns "DEBUG", "INFO", "WARN", "ERROR", "FATAL", or "UNKNOWN" for invalid levels.
*/
func (l LogLevel) String() string {
	switch l {
	case DEBUG:
		return "DEBUG"
	case INFO:
		return "INFO"
	case WARN:
		return "WARN"
	case ERROR:
		return "ERROR"
	case FATAL:
		return "FATAL"
	default:
		return "UNKNOWN"
	}
}

// Logger is a structured logger
type Logger struct {
	level  LogLevel
	logger *log.Logger
}

/*
New creates a new logger instance with the specified log level.
The level parameter should be one of: "debug", "info", "warn", "error", "fatal".
Messages below the specified level will not be logged.
Example: New("info") will log INFO, WARN, ERROR, and FATAL, but not DEBUG.
*/
func New(level string) *Logger {
	return &Logger{
		level:  parseLogLevel(level),
		logger: log.New(os.Stdout, "", 0),
	}
}

/*
Debug logs a debug-level message with optional key-value field pairs.
Debug messages are only logged if the logger level is set to DEBUG.
Useful for detailed diagnostic information during development.
Example: logger.Debug("User query", "user_id", 123, "query_time_ms", 45)
*/
func (l *Logger) Debug(msg string, fields ...interface{}) {
	if l.level <= DEBUG {
		l.log(DEBUG, msg, fields...)
	}
}

/*
Info logs an informational message with optional key-value field pairs.
Info messages are logged if the logger level is DEBUG or INFO.
Useful for general application flow information.
Example: logger.Info("Server started", "port", 6969)
*/
func (l *Logger) Info(msg string, fields ...interface{}) {
	if l.level <= INFO {
		l.log(INFO, msg, fields...)
	}
}

/*
Warn logs a warning message with optional key-value field pairs.
Warning messages indicate potentially harmful situations that don't prevent operation.
Example: logger.Warn("Slow query detected", "duration_ms", 5000)
*/
func (l *Logger) Warn(msg string, fields ...interface{}) {
	if l.level <= WARN {
		l.log(WARN, msg, fields...)
	}
}

/*
Error logs an error message with optional key-value field pairs.
Error messages indicate failures that should be investigated.
Example: logger.Error("Database connection failed", "error", err.Error())
*/
func (l *Logger) Error(msg string, fields ...interface{}) {
	if l.level <= ERROR {
		l.log(ERROR, msg, fields...)
	}
}

/*
Fatal logs a fatal error message and terminates the application with exit code 1.
Use this only for unrecoverable errors that prevent the application from continuing.
Example: logger.Fatal("Failed to load configuration", "error", err.Error())
*/
func (l *Logger) Fatal(msg string, fields ...interface{}) {
	l.log(FATAL, msg, fields...)
	os.Exit(1)
}

// log formats and writes the log message
func (l *Logger) log(level LogLevel, msg string, fields ...interface{}) {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	levelStr := level.String()

	// Format fields as key=value pairs
	var fieldStr string
	if len(fields) > 0 {
		fieldStr = " " + formatFields(fields...)
	}

	// Color codes for terminal output
	colorCode := getColorCode(level)
	resetCode := "\033[0m"

	output := fmt.Sprintf("%s[%s] %s%-5s%s %s%s",
		colorCode,
		timestamp,
		colorCode,
		levelStr,
		resetCode,
		msg,
		fieldStr,
	)

	l.logger.Println(output)
}

// formatFields converts field pairs to key=value format
func formatFields(fields ...interface{}) string {
	if len(fields)%2 != 0 {
		return fmt.Sprintf("invalid_fields=%v", fields)
	}

	var parts []string
	for i := 0; i < len(fields); i += 2 {
		key := fmt.Sprintf("%v", fields[i])
		value := fmt.Sprintf("%v", fields[i+1])
		parts = append(parts, fmt.Sprintf("%s=%s", key, value))
	}

	return strings.Join(parts, " ")
}

// parseLogLevel converts string to LogLevel
func parseLogLevel(level string) LogLevel {
	switch strings.ToLower(level) {
	case "debug":
		return DEBUG
	case "info":
		return INFO
	case "warn", "warning":
		return WARN
	case "error":
		return ERROR
	case "fatal":
		return FATAL
	default:
		return INFO
	}
}

// getColorCode returns ANSI color code for log level
func getColorCode(level LogLevel) string {
	switch level {
	case DEBUG:
		return "\033[36m" // Cyan
	case INFO:
		return "\033[32m" // Green
	case WARN:
		return "\033[33m" // Yellow
	case ERROR:
		return "\033[31m" // Red
	case FATAL:
		return "\033[35m" // Magenta
	default:
		return "\033[0m" // Reset
	}
}
