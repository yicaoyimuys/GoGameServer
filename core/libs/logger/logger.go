package logger

import (
	"os"

	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	fileLogger    *zap.Logger
	consoleLogger *zap.Logger
	conf          option
)

type option struct {
	debug bool
	both  bool
	file  bool
	name  string
}

// 默认参数
func defaultOption() option {
	return option{
		debug: false,
		both:  true,
		file:  true,
		name:  "log",
	}
}

// Option 参数
type Option func(*option)

// Init Init
func Init(opts ...Option) {
	conf := defaultOption()
	for _, opt := range opts {
		opt(&conf)
	}

	var logLevel = zapcore.InfoLevel
	if conf.debug {
		logLevel = zapcore.DebugLevel
	}

	if conf.both {
		createConsoleLogger(conf.name, logLevel)
	}

	if conf.file {
		createFileLogger(conf.name, logLevel)
	}
}

func createFileLogger(name string, logLevel zapcore.Level) {
	lumberJackLogger := &lumberjack.Logger{
		Filename:   "./logs/" + name + ".log",
		MaxSize:    10,
		MaxBackups: 5,
		MaxAge:     30,
		Compress:   false,
	}
	writeSyncer := zapcore.AddSync(lumberJackLogger)

	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	encoder := zapcore.NewJSONEncoder(encoderConfig)

	fileLogger = zap.New(zapcore.NewCore(encoder, writeSyncer, logLevel))
	fileLogger = fileLogger.Named(name)
}

func createConsoleLogger(fileName string, logLevel zapcore.Level) {
	writeSyncer := os.Stderr

	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	encoder := zapcore.NewConsoleEncoder(encoderConfig)

	consoleLogger = zap.New(zapcore.NewCore(encoder, writeSyncer, logLevel))
	consoleLogger = consoleLogger.Named(fileName)
}

// WithDebug 是否开启Debug参数
func WithDebug(debug bool) Option {
	return func(o *option) {
		o.debug = debug
	}
}

// WithBoth 是否开启Both参数
func WithBoth(both bool) Option {
	return func(o *option) {
		o.both = both
	}
}

// WithFile 是否开启File参数
func WithFile(file bool) Option {
	return func(o *option) {
		o.file = file
	}
}

// WithName 设置name参数
func WithName(name string) Option {
	return func(o *option) {
		o.name = name
	}
}

// Error Error
func Error(msg string, fields ...zap.Field) {
	if consoleLogger != nil {
		defer consoleLogger.Sync()
		consoleLogger.Error(msg, fields...)
	}

	if fileLogger != nil {
		defer fileLogger.Sync()
		fileLogger.Error(msg, fields...)
	}
}

// Warn Warn
func Warn(msg string, fields ...zap.Field) {
	if consoleLogger != nil {
		defer consoleLogger.Sync()
		consoleLogger.Warn(msg, fields...)
	}

	if fileLogger != nil {
		defer fileLogger.Sync()
		fileLogger.Warn(msg, fields...)
	}
}

// Info Info
func Info(msg string, fields ...zap.Field) {
	if consoleLogger != nil {
		defer consoleLogger.Sync()
		consoleLogger.Info(msg, fields...)
	}

	if fileLogger != nil {
		defer fileLogger.Sync()
		fileLogger.Info(msg, fields...)
	}
}

// Debug Debug
func Debug(msg string, fields ...zap.Field) {
	if consoleLogger != nil {
		defer consoleLogger.Sync()
		consoleLogger.Debug(msg, fields...)
	}

	if fileLogger != nil {
		defer fileLogger.Sync()
		fileLogger.Debug(msg, fields...)
	}
}
