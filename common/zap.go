package common

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var(
	logger *zap.SugaredLogger
)

func init()  {
	// Log file name
	fileName := "micro.log"
	syncWriter:= zapcore.AddSync(
		&lumberjack.Logger{
			Filename:   fileName, // File name
			MaxSize:    512,//MB
			//MaxAge:     0,
			MaxBackups: 0, // Maximum backups
			LocalTime:  true,
			Compress:   true, // Enable compression
		})
	// Encoding
	encoder:=zap.NewProductionEncoderConfig()
	// Time format
	encoder.EncodeTime = zapcore.ISO8601TimeEncoder
	core:= zapcore.NewCore(
		// Encoder
		zapcore.NewJSONEncoder(encoder),
		syncWriter,
		//
		zap.NewAtomicLevelAt(zap.DebugLevel))
	log := zap.New(
		core,
		zap.AddCaller(),
		zap.AddCallerSkip(1))
	logger = log.Sugar()
}

func Debug(args ...interface{})  {
	logger.Debug(args)
}

func Debugf(template string, args ...interface{}) {
	logger.Debugf(template, args...)
}

func Info(args ...interface{}) {
	logger.Info(args...)
}

func Infof(template string, args ...interface{}) {
	logger.Infof(template, args...)
}

func Warn(args ...interface{}) {
	logger.Warn(args...)
}

func Warnf(template string, args ...interface{}) {
	logger.Warnf(template, args...)
}

func Error(args ...interface{}) {
	logger.Error(args...)
}

func Errorf(template string, args ...interface{}) {
	logger.Errorf(template, args...)
}

func DPanic(args ...interface{}) {
	logger.DPanic(args...)
}

func DPanicf(template string, args ...interface{}) {
	logger.DPanicf(template, args...)
}

func Panic(args ...interface{}) {
	logger.Panic(args...)
}

func Panicf(template string, args ...interface{}) {
	logger.Panicf(template, args...)
}

func Fatal(args ...interface{}) {
	logger.Fatal(args...)
}

func Fatalf(template string, args ...interface{}) {
	logger.Fatalf(template, args...)
}