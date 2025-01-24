package helpers

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/fatih/color"
)

// TODO:  to add a buffer to file log, fottow:
// https://chatgpt.com/c/67863b1c-62fc-8013-a1ca-3c0585e901b2

var (
	// Loggers
	infoLog  *log.Logger
	errorLog *log.Logger
	fatalLog *log.Logger

	// Log configuration
	infoLogFile  *os.File
	errorLogFile *os.File

	infoLogBuffer  []string
	errorLogBuffer []errorLogMsg
	bufferLock     sync.Mutex

	flushInterval   time.Duration
	stopFlushSignal chan struct{}
)

type errorLogMsg struct {
	dept int
	msg  string
}

func ConfigLogger() {
	// Get log mode from environment variables
	infoLogMode := os.Getenv("INFO_LOG_MODE")           // "console", "file", or "off"
	errorLogMode := os.Getenv("ERROR_LOG_MODE")         // "console", "file", or "off"
	flushIntervalEnv := os.Getenv("LOG_FLUSH_INTERVAL") // Flush interval in seconds

	// Default flush interval to 10 seconds
	flushInterval = 10 * time.Second
	if flushIntervalEnv != "" {
		if interval, err := time.ParseDuration(flushIntervalEnv + "s"); err == nil {
			flushInterval = interval
		}
	}

	// Initialize stop signal for flushing
	stopFlushSignal = make(chan struct{})

	// Set up info logger
	switch infoLogMode {
	case "file":
		infoLogFile = setupLogFile("dist/logs/info.log")
		infoLog = log.New(infoLogFile, "INFO\t", log.Ldate|log.Ltime)
	case "off":
		infoLog = log.New(nil, "", 0) // Discard logs
	default:
		infoLog = log.New(os.Stdout, color.New(color.FgBlue).Sprint("INFO\t"), log.Ldate|log.Ltime)
	}

	// Set up error logger
	switch errorLogMode {
	case "file":
		errorLogFile = setupLogFile("dist/logs/error.log")
		errorLog = log.New(errorLogFile, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
		fatalLog = log.New(errorLogFile, "FATAL\t", log.Ldate|log.Ltime|log.Lshortfile)
	case "off":
		errorLog = log.New(nil, "", 0) // Discard logs
		fatalLog = log.New(nil, "", 0) // Discard logs
	default:
		errorLog = log.New(os.Stderr, color.New(color.FgRed).Sprint("ERROR\t"), log.Ldate|log.Ltime|log.Lshortfile)
		fatalLog = log.New(os.Stderr, color.New(color.FgRed, color.Bold).Sprint("FATAL\t"), log.Ldate|log.Ltime|log.Lshortfile)
	}

	// Start background flush goroutine if logging to files
	if infoLogMode == "file" || errorLogMode == "file" {
		go flushLogsPeriodically()
	}
}

// LogInfo logs an informational message with formatting options.
func LogInfo(format string, args ...interface{}) {
	logMessage := fmt.Sprintf(format, args...)

	if infoLogFile != nil {
		bufferLock.Lock()
		infoLogBuffer = append(infoLogBuffer, logMessage)
		bufferLock.Unlock()

	} else {
		infoLog.Println(logMessage)
	}
}

// LogError logs error messages with context including the correct caller location.
func LogError(err error, msg string, depth ...int) {
	// Set the default depth to 2 if none is provided
	callDepth := 2
	if len(depth) > 0 {
		callDepth = depth[0]
	}

	logMessage := ""

	if err != nil {
		logMessage = fmt.Sprintf("%s:\n%v", msg, err)
	} else {
		// errorLog.Output(callDepth, fmt.Sprintf("%s:\n", msg))
		logMessage = fmt.Sprintf("%s:\n", msg)
	}

	if errorLogFile != nil {
		bufferLock.Lock()
		errorLogBuffer = append(errorLogBuffer, errorLogMsg{msg: logMessage, dept: callDepth})
		bufferLock.Unlock()
	} else {
		errorLog.Output(callDepth, logMessage)
	}
}

// LogFatal logs critical errors and terminates the application.
func LogFatal(err error, msg string, depth ...int) {
	callDepth := 2 // Default call depth
	if len(depth) > 0 {
		callDepth = depth[0]
	}

	logMessage := ""

	if err != nil {
		logMessage = fmt.Sprintf("%s:\n%v", msg, err)
	} else {
		// errorLog.Output(callDepth, fmt.Sprintf("%s:\n", msg))
		logMessage = fmt.Sprintf("%s:\n", msg)
	}

	if errorLogFile != nil {
		bufferLock.Lock()
		errorLogBuffer = append(errorLogBuffer, errorLogMsg{msg: logMessage, dept: callDepth})
		bufferLock.Unlock()
	}

	if errorLogFile != nil || infoLogFile != nil {
		flushLogs() // Flush logs immediately before exiting

	} else {
		fatalLog.Output(callDepth, logMessage)
	}

	os.Exit(1) // Terminate the application with a non-zero status code
}

// flushLogs writes buffered logs to their respective files.
func flushLogs() {
	bufferLock.Lock()
	defer bufferLock.Unlock()

	// Flush info logs
	if infoLogFile != nil && len(infoLogBuffer) > 0 {
		for _, log := range infoLogBuffer {
			infoLog.Printf(log)
		}
		infoLogBuffer = []string{}
	}

	// Flush error logs
	if errorLogFile != nil && len(errorLogBuffer) > 0 {
		for _, log := range errorLogBuffer {
			// fmt.Fprintln(errorLogFile, log)
			errorLog.Output(log.dept, log.msg)
		}
		errorLogBuffer = []errorLogMsg{}
	}
}

// flushLogsPeriodically flushes logs at regular intervals.
func flushLogsPeriodically() {
	ticker := time.NewTicker(flushInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			flushLogs()
		case <-stopFlushSignal:
			flushLogs() // Final flush on stop signal
			return
		}
	}
}

// StopLogging stops the periodic flushing and closes log files.
func StopLogging() {
	close(stopFlushSignal)
	if infoLogFile != nil {
		infoLogFile.Close()
	}
	if errorLogFile != nil {
		errorLogFile.Close()
	}
}

// PrettyPrintJSON takes any data and prints it in an indented JSON format
func PrettyPrintJSON(data any) {
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		fmt.Println("Error formatting JSON:", err)
		return
	}
	fmt.Println(string(jsonData))
}

// setupLogFile creates or opens a log file for writing and ensures the directory exists.
func setupLogFile(filename string) *os.File {
	dir := filepath.Dir(filename)
	// Ensure directory exists
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err := os.MkdirAll(dir, 0755)
		if err != nil {
			log.Fatalf("Failed to create log directory %s: %v", dir, err)
		}
	}

	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatalf("Failed to open log file %s: %v", filename, err)
	}
	return file
}
