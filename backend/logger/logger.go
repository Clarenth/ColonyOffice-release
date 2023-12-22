package logger

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

// CustomLogger instances a logger using a customizedd gin.LoggerWithFormatter
func CustomLogger(logsDir string) gin.HandlerFunc {

	const layout string = "2006-01-02"
	t := time.Now()
	logFileName := t.Format(layout) + ".log"

	// fileOutput opens the log file in the logs folder destination for writing
	fileOutput, err := os.OpenFile(logsDir+"/"+logFileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
	}

	loggerActivityWriter := io.MultiWriter(fileOutput, os.Stdout)

	return gin.LoggerWithConfig(gin.LoggerConfig{
		Formatter: customGinFormatter,
		Output:    loggerActivityWriter,
	})
}

var customGinFormatter = func(param gin.LogFormatterParams) string {
	const logTime string = "2006-01-02-T15:04:05-MST"
	if param.Latency > time.Minute {
		// Truncate in a golang < 1.8 safe way
		param.Latency = param.Latency - param.Latency%time.Second
	}
	return fmt.Sprintf("[GIN]|Date:%v|Status:%1d|IP:%5v|Lat:%2s|%s %v|Agent:%s|\n",
		// "[GIN]|Date:%v|Status:%1d|IP:%5v|Lat:%2s|Content-Length:%d|Keys:%s|%s %v|Agent:%s|Request:%s %s|error:%s|\n"
		param.TimeStamp.Format(logTime),
		param.StatusCode,
		param.ClientIP,
		param.Latency,
		param.Method,
		param.Request.URL,
		param.Request.UserAgent(),
	)
}

// createLoggerFolder creates an activity logs folder for the server.
// In the future this will make a folder for each unique user, and each session connection.
func createLoggerDir(logsDir string) string {
	// Checks to see if the logs folder exists using os.Stat, and creates the folder
	dirName := logsDir
	//dirName := "./logs"
	if _, err := os.Stat(dirName); errors.Is(err, os.ErrNotExist) {
		err := os.MkdirAll(dirName, 0700)
		if err != nil {
			log.Println(err)
		}
	}
	return dirName
}

// createLogFileNameFromDate derives the file name from the current date
func createLogFileNameFromDate() string {
	// This currently only creates a file when the server activates. In future
	// this will create based on a user session connecting each time .
	const layout string = "2006-01-02"
	t := time.Now()
	return t.Format(layout) + ".log"
}
