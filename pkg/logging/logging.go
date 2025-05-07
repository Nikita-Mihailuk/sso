package logging

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

var logger *zap.Logger

func GetLogger() *zap.Logger {
	return logger
}

func init() {
	consoleEncoder := zapcore.NewConsoleEncoder(zapcore.EncoderConfig{
		TimeKey:          "time",
		LevelKey:         "level",
		CallerKey:        "caller",
		MessageKey:       "msg",
		EncodeLevel:      zapcore.CapitalColorLevelEncoder,
		EncodeCaller:     zapcore.FullCallerEncoder,
		EncodeTime:       zapcore.ISO8601TimeEncoder,
		ConsoleSeparator: " | ",
	})

	fileEncoder := zapcore.NewJSONEncoder(zapcore.EncoderConfig{
		TimeKey:      "time",
		LevelKey:     "level",
		CallerKey:    "caller",
		MessageKey:   "msg",
		EncodeCaller: zapcore.FullCallerEncoder,
		EncodeTime:   zapcore.ISO8601TimeEncoder,
		EncodeLevel:  zapcore.LowercaseLevelEncoder,
	})

	logFile, err := os.OpenFile("pkg/logging/all.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		panic("log file is not find")
	}

	fileWriter := zapcore.AddSync(logFile)
	consoleWriter := zapcore.AddSync(os.Stdout)

	core := zapcore.NewTee(
		zapcore.NewCore(consoleEncoder, consoleWriter, zapcore.DebugLevel),
		zapcore.NewCore(fileEncoder, fileWriter, zapcore.DebugLevel),
	)

	logger = zap.New(core, zap.AddCaller())
}
