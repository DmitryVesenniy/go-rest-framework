package logger

type LoggerInterface interface {
	Info(message string)
	Warning(message string)
	Error(message string)
}
