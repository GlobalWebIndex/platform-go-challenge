package utils

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

var (
	InfoLogger  *log.Logger
	WarnLogger  *log.Logger
	ErrorLogger *log.Logger
	DebugLogger *log.Logger
)

type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
	body       bytes.Buffer
}

func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

func (lrw *loggingResponseWriter) Write(b []byte) (int, error) {
	lrw.body.Write(b)
	return lrw.ResponseWriter.Write(b)
}

func LogHandler(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Log request
		requestBody, _ := io.ReadAll(r.Body)
		r.Body = io.NopCloser(bytes.NewBuffer(requestBody))
		InfoLogger.Printf("Request: %s %s\nHeaders: %v\nBody: %s", r.Method, r.URL, r.Header, string(requestBody))

		// Create custom response writer
		lrw := &loggingResponseWriter{ResponseWriter: w, statusCode: http.StatusOK}

		// Call the original handler
		handler.ServeHTTP(lrw, r)

		// Log response
		InfoLogger.Printf("Response: Status: %d\nHeaders: %v\nBody: %s", lrw.statusCode, lrw.Header(), lrw.body.String())
	}
}

func SetupLogger() (*os.File, error) {
	// Create logs directory if it doesn't exist
	err := os.MkdirAll("logs", 0755)
	if err != nil {
		return nil, fmt.Errorf("failed to create logs directory: %v", err)
	}

	// Open a file for writing logs
	currentDate := time.Now().Format("2006-01-02")
	logFileName := fmt.Sprintf("app-%s.log", currentDate)
	logFile, err := os.OpenFile(filepath.Join("logs", logFileName), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to open log file: %v", err)
	}

	// Create multi-writer for both file and stdout
	multiWriter := io.MultiWriter(os.Stdout, logFile)

	// Set up loggers for different severity levels
	InfoLogger = log.New(multiWriter, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	WarnLogger = log.New(multiWriter, "WARN: ", log.Ldate|log.Ltime|log.Lshortfile)
	ErrorLogger = log.New(multiWriter, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
	DebugLogger = log.New(multiWriter, "DEBUG: ", log.Ldate|log.Ltime|log.Lshortfile)

	return logFile, nil
}
