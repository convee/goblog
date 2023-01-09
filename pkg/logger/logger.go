package logger

import (
	"fmt"
	"go.uber.org/zap"
)

// log is A global variable so that log functions can be directly accessed
var log *zap.Logger

// Fields Type to pass when we want to call WithFields for structured logging
type Fields map[string]interface{}

// Logger is a contract for the logger
type Logger interface {
	Info(args ...interface{})
	Warn(args ...interface{})
	Error(args ...interface{})
}

// Init init log
func Init(cfg *Config) *zap.Logger {
	var err error
	// new zap logger
	log, err = newZapLogger(cfg)
	if err != nil {
		_ = fmt.Errorf("init newZapLogger err: %v", err)
	}
	return log
}

// GetLogger return a log
func GetLogger() *zap.Logger {
	return log
}

// Info logger
func Info(msg string, fields ...zap.Field) {
	//requestId := new(http.Request).Header.Get("request_id")
	//fields = append(fields, zap.String("request_id", requestId))
	log.Info(msg, fields...)
}

// Warn logger
func Warn(msg string, fields ...zap.Field) {
	//requestId := new(http.Request).Header.Get("request_id")
	//fields = append(fields, zap.String("request_id", requestId))
	log.Warn(msg, fields...)
}

// Error logger
func Error(msg string, fields ...zap.Field) {
	//requestId := new(http.Request).Header.Get("request_id")
	//fields = append(fields, zap.String("request_id", requestId))
	log.Error(msg, fields...)
}
