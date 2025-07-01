package utils

import (
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	SystemLogger zerolog.Logger
	AccessLogger zerolog.Logger
	AuthLogger   zerolog.Logger
	BotLogger    zerolog.Logger
)

func InitLogger() {
	zerolog.TimeFieldFormat = time.RFC3339

	logDir := "logs"
	_ = os.MkdirAll(logDir, 0755)

	SystemLogger = newRotatingLogger(logDir + "/system.log")
	AccessLogger = newRotatingLogger(logDir + "/access.log")
	AuthLogger = newRotatingLogger(logDir + "/auth.log")
	BotLogger = newRotatingLogger(logDir + "/bot.log")

	// Default global logger = System
	log.Logger = SystemLogger
}

func newRotatingLogger(path string) zerolog.Logger {
	writer := &lumberjack.Logger{
		Filename:   path,
		MaxSize:    10, // MB
		MaxBackups: 5,
		MaxAge:     28,
		Compress:   true,
	}
	return zerolog.New(writer).With().Timestamp().Logger()
}
