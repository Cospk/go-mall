package logger

import (
	"github.com/Cospk/go-mall/pkg/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
)

var (
	Logger *zap.Logger
)

func InitLogger() {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoder := zapcore.NewJSONEncoder(encoderConfig)

	var cores []zapcore.Core
	if config.AppConfig.App.Env == "dev" {
		// 开发环境：控制台和文件都要日志，且是debug级别
		cores = append(
			cores,
			zapcore.NewCore(encoder, getFileLogWriter(), zapcore.DebugLevel),
			zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), zapcore.DebugLevel),
		)
	} else {
		// 生产环境: 只要文件日志，且是info级别
		cores = append(cores, zapcore.NewCore(encoder, getFileLogWriter(), zapcore.InfoLevel))
	}

	core := zapcore.NewTee(cores...)
	Logger = zap.New(core)
}

func getFileLogWriter() (writeSyncer zapcore.WriteSyncer) {
	// 使用lumberjack 实现logger rotate
	lumberJackLogger := &lumberjack.Logger{
		Filename:  config.AppConfig.App.Log.FilePath,
		MaxSize:   config.AppConfig.App.Log.FileMaxSize,
		MaxAge:    config.AppConfig.App.Log.BackUpFileMaxAge,
		Compress:  false,
		LocalTime: true,
	}
	return zapcore.AddSync(lumberJackLogger)
}
