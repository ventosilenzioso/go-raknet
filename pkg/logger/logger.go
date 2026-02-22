package logger

import (
	"fmt"
	"log"
	"os"
	"time"
)

// ANSI color codes
const (
	ColorReset  = "\033[0m"
	ColorRed    = "\033[31m"
	ColorGreen  = "\033[32m"
	ColorYellow = "\033[33m"
	ColorWhite  = "\033[37m"
	ColorCyan   = "\033[36m"
	ColorGray   = "\033[90m"
)

// Log levels
const (
	LevelDebug = iota
	LevelInfo
	LevelWarn
	LevelError
	LevelSuccess
)

// Logger represents a colored logger
type Logger struct {
	level      int
	timeFormat string
	showTime   bool
}

var defaultLogger *Logger

func init() {
	defaultLogger = &Logger{
		level:      LevelInfo,
		timeFormat: "15:04:05",
		showTime:   true,
	}
}

// SetLevel sets the minimum log level
func SetLevel(level int) {
	defaultLogger.level = level
}

// SetTimeFormat sets the time format for logs
func SetTimeFormat(format string) {
	defaultLogger.timeFormat = format
}

// ShowTime enables or disables timestamp in logs
func ShowTime(show bool) {
	defaultLogger.showTime = show
}

// formatMessage formats a log message with color and timestamp
func (l *Logger) formatMessage(level, color, prefix, message string) string {
	timestamp := ""
	if l.showTime {
		timestamp = fmt.Sprintf("%s[%s]%s ", ColorGray, time.Now().Format(l.timeFormat), ColorReset)
	}
	return fmt.Sprintf("%s%s[%s]%s %s", timestamp, color, prefix, ColorReset, message)
}

// Debug logs a debug message (gray)
func Debug(format string, args ...interface{}) {
	if defaultLogger.level <= LevelDebug {
		msg := fmt.Sprintf(format, args...)
		log.Println(defaultLogger.formatMessage("DEBUG", ColorGray, "DEBUG", msg))
	}
}

// Info logs an informational message (white)
func Info(format string, args ...interface{}) {
	if defaultLogger.level <= LevelInfo {
		msg := fmt.Sprintf(format, args...)
		log.Println(defaultLogger.formatMessage("INFO", ColorWhite, "INFO", msg))
	}
}

// Warn logs a warning message (yellow)
func Warn(format string, args ...interface{}) {
	if defaultLogger.level <= LevelWarn {
		msg := fmt.Sprintf(format, args...)
		log.Println(defaultLogger.formatMessage("WARN", ColorYellow, "WARN", msg))
	}
}

// Error logs an error message (red)
func Error(format string, args ...interface{}) {
	if defaultLogger.level <= LevelError {
		msg := fmt.Sprintf(format, args...)
		log.Println(defaultLogger.formatMessage("ERROR", ColorRed, "ERROR", msg))
	}
}

// Success logs a success message (green)
func Success(format string, args ...interface{}) {
	if defaultLogger.level <= LevelSuccess {
		msg := fmt.Sprintf(format, args...)
		log.Println(defaultLogger.formatMessage("SUCCESS", ColorGreen, "SUCCESS", msg))
	}
}

// Fatal logs a fatal error and exits
func Fatal(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	log.Println(defaultLogger.formatMessage("FATAL", ColorRed, "FATAL", msg))
	os.Exit(1)
}
// InfoCyan logs an info message in cyan (for special highlights)
func InfoCyan(format string, args ...interface{}) {
	if defaultLogger.level <= LevelInfo {
		msg := fmt.Sprintf(format, args...)
		log.Println(defaultLogger.formatMessage("INFO", ColorCyan, "INFO", msg))
	}
}

// Section prints a section header
func Section(title string) {
	border := "═══════════════════════════════════════════════════════════"
	fmt.Printf("\n%s╔%s╗%s\n", ColorCyan, border, ColorReset)
	fmt.Printf("%s║%s %-57s %s║%s\n", ColorCyan, ColorReset, title, ColorCyan, ColorReset)
	fmt.Printf("%s╚%s╝%s\n\n", ColorCyan, border, ColorReset)
}

// Banner prints the application banner
func Banner(title, version string) {
	banner := `
╔═══════════════════════════════════════════════════════════╗
║                                                           ║
║   ███████╗ █████╗       ███╗   ███╗██████╗               ║
║   ██╔════╝██╔══██╗      ████╗ ████║██╔══██╗              ║
║   ███████╗███████║█████╗██╔████╔██║██████╔╝              ║
║   ╚════██║██╔══██║╚════╝██║╚██╔╝██║██╔═══╝               ║
║   ███████║██║  ██║      ██║ ╚═╝ ██║██║                   ║
║   ╚══════╝╚═╝  ╚═╝      ╚═╝     ╚═╝╚═╝                   ║
║                                                           ║
║              %s%-37s%s║
║                    %sVersion %-7s%s                      ║
║                                                           ║
╚═══════════════════════════════════════════════════════════╝
`
	fmt.Printf(banner, ColorCyan, title, ColorReset, ColorGreen, version, ColorReset)
}
