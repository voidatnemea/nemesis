package logger

import (
	"fmt"
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Log *zap.Logger

func Init() error {
	start := time.Now()
	cfg := zap.NewProductionConfig()
	cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	cfg.EncoderConfig.StacktraceKey = "" // Disable stacktrace for production clean logs

	// Ensure log directory exists
	_ = os.MkdirAll("storage/logs", os.ModePerm)

	var err error
	Log, err = cfg.Build()
	if err != nil {
		return fmt.Errorf("failed to build logger: %w", err)
	}

	elapsed := time.Since(start)
	fmt.Printf("(%s)[OK][LOGGER] Logger initialized in %s\n", time.Now().Format(time.RFC3339), elapsed)
	return nil
}

func GetLogger() *zap.Logger {
	return Log
}
