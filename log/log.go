package log

import (
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

func NewLogger(logLevel string, logFilename string, logMaxSize int, logMaxBackups int, logMaxAge int, logCompress bool) *zap.Logger {
	var ll zapcore.Level
	switch logLevel {
	case "DEBUG":
		ll = zap.DebugLevel
	case "INFO":
		ll = zap.InfoLevel
	case "WARN", "WARNING":
		ll = zap.WarnLevel
	case "ERR", "ERROR":
		ll = zap.WarnLevel
	case "DPANIC":
		ll = zap.DPanicLevel
	case "PANIC":
		ll = zap.PanicLevel
	case "FATAL":
		ll = zap.FatalLevel
	}

	var ws zapcore.WriteSyncer
	switch logFilename {
	case "", os.Stderr.Name():
		ws = zapcore.AddSync(os.Stderr)
	case os.Stdout.Name():
		ws = zapcore.AddSync(os.Stdout)
	default:
		ws = zapcore.AddSync(
			&lumberjack.Logger{
				Filename:   logFilename,
				MaxSize:    logMaxSize, // megabytes
				MaxBackups: logMaxBackups,
				MaxAge:     logMaxAge, // days
				Compress:   logCompress,
			},
		)
	}

	ec := zap.NewProductionEncoderConfig()
	ec.TimeKey = "_timestamp_"
	ec.LevelKey = "_level_"
	ec.NameKey = "_name_"
	ec.CallerKey = "_caller_"
	ec.MessageKey = "_message_"
	ec.StacktraceKey = "_stacktrace_"
	ec.EncodeTime = zapcore.ISO8601TimeEncoder
	ec.EncodeCaller = zapcore.ShortCallerEncoder

	logger := zap.New(
		zapcore.NewCore(
			zapcore.NewJSONEncoder(ec),
			ws,
			ll,
		),
		zap.AddCaller(),
	).Named("ethpayserver")

	return logger
}
