package config

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

func InitLogger(logsDir string) *zap.Logger {

	w := zapcore.AddSync(&lumberjack.Logger{
		Filename: logsDir,
		MaxSize: 20,
		MaxBackups: 3,
		MaxAge: 28,
	})

	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
		w,
		zap.InfoLevel,
	)

	logger := zap.New(core)

	return logger

}
