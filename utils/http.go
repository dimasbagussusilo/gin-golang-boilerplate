package utils

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/natefinch/lumberjack"
	"io"
	"math"
	"os"
	"time"
)

type ResponseStatus string

const (
	ResponseStatusSuccess ResponseStatus = "success"
	ResponseStatusError   ResponseStatus = "error"
)

type Response struct {
	Status  ResponseStatus `json:"status"`
	Message string         `json:"message"`
	Data    any            `json:"data"`
}

func ResponseData(status ResponseStatus, message string, data any) *Response {
	return &Response{
		Status:  status,
		Message: message,
		Data:    data,
	}
}

type Pagination struct {
	TotalItems int64 `json:"total_items"`
	TotalPages int   `json:"total_pages"`
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
	HasPrev    bool  `json:"has_prev"`
	HasNext    bool  `json:"has_next"`
}

func Paginate(count int64, page int, limit int) *Pagination {
	var pagination Pagination

	pagination.TotalItems = count
	pagination.Page = page
	pagination.Limit = limit

	if limit > 0 {
		pagination.TotalPages = int(math.Ceil(float64(count) / float64(limit)))
	} else {
		pagination.TotalPages = 0
	}

	if page > 1 {
		pagination.HasPrev = true
	}

	if page < pagination.TotalPages {
		pagination.HasNext = true
	}

	return &pagination
}

type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func (w bodyLogWriter) WriteString(s string) (int, error) {
	w.body.WriteString(s)
	return w.ResponseWriter.WriteString(s)
}

func Logger() gin.HandlerFunc {
	// Create a new log file
	logFile, err := os.OpenFile(
		"./gin.log",
		os.O_APPEND|os.O_CREATE|os.O_WRONLY,
		0666,
	)
	if err != nil {
		panic(fmt.Sprintf("Error while open log file: %e", err))
	}

	// Set up log rotation every 5 days
	rotationTime := time.Now().AddDate(0, 0, 5)
	lumberjackLogger := &lumberjack.Logger{
		Filename:   "./gin.log",
		MaxSize:    10,   // Max size of each log file in megabytes
		MaxBackups: 30,   // Max number of log files to keep
		MaxAge:     5,    // Max number of days to keep old log files
		Compress:   true, // Compress the rotated log files
	}

	return func(c *gin.Context) {
		requestID := uuid.New().String()
		c.Set("RequestID", requestID)

		// Log the request time, method, and URL
		startTime := time.Now()

		buf, _ := io.ReadAll(c.Request.Body)
		c.Request.Body = io.NopCloser(bytes.NewBuffer(buf))

		blw := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = blw

		c.Next()

		endTime := time.Now()

		logMessage := fmt.Sprintf(
			"{\"date_time\":\"%s\", \"request_id\":\"%s\", \"client_ip\":\"%s\", \"request_method\":\"%s\", \"endpoint\":\"%s\", \"response_status\":%d, \"latency\":\"%s\", \"response_body\":\"%s\"}\n",
			endTime.Format("2006/01/02-15:04:05"),
			requestID,
			c.ClientIP(),
			c.Request.Method,
			c.Request.URL.String(),
			c.Writer.Status(),
			endTime.Sub(startTime).String(),
			blw.body.String(),
		)

		// Write the log message to the file and to stdout
		_, err := logFile.WriteString(logMessage)
		if err != nil {
			panic(fmt.Sprintf("Error while write log file: %e", err))
		}

		// Rotate the log file if necessary
		if time.Now().After(rotationTime) {
			err := logFile.Close()
			if err != nil {
				panic(fmt.Sprintf("Error while close log file: %e", err))
			}
			err = lumberjackLogger.Rotate()
			if err != nil {
				panic(fmt.Sprintf("Error while rotate log file: %e", err))
			}
			logFile, err = os.OpenFile(
				"./gin.log",
				os.O_APPEND|os.O_CREATE|os.O_WRONLY,
				0666,
			)
			if err != nil {
				panic(err)
			}
			rotationTime = time.Now().AddDate(0, 0, 5)
		}
	}
}
