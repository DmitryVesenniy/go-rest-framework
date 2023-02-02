package logger

import (
	"log"
	"os"
)

type AppLogger struct {
	infoLogger    *log.Logger
	warningLogger *log.Logger
	errorLogger   *log.Logger
}

func (l *AppLogger) Info(message string) {
	if l.infoLogger != nil {
		l.infoLogger.Println(message)
	}
}

func (l *AppLogger) Warning(message string) {
	if l.warningLogger != nil {
		l.warningLogger.Println(message)
	}
}

func (l *AppLogger) Error(message string) {
	if l.errorLogger != nil {
		l.errorLogger.Println(message)
	}
}

func Init(filePath string) (LoggerInterface, error) {
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return nil, err
	}

	infoLogger := log.New(file, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	warningLogger := log.New(file, "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile)
	errorLogger := log.New(file, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)

	logger := &AppLogger{
		infoLogger:    infoLogger,
		warningLogger: warningLogger,
		errorLogger:   errorLogger,
	}

	return logger, nil
}
