package logger

import (
	"os"

	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Field = zap.Field

var (
	Logger		*zap.Logger
	String		= zap.String
	Any			= zap.Any
	Int			= zap.Int
	Float32 	= zap.Float32
)

// logpath = log file path
// loglevel = log level
func InitLogger(logpath, loglevel string) {
	// log splitting
	hook := lumberjack.Logger {
		Filename: 		logpath,	// log file path, default os.TempDir()
		MaxSize: 		100,		// each log file saves 100MB, the default is 100MB
		MaxBackups: 	30,			// keep 30 backups, the default is unlimited
		MaxAge: 		7,			// keep for 7 days, default unlimited
		Compress: 		true,		// whether to compress, default not to compress
	}

	write := zapcore.AddSync(&hook)

	// set log level
	// debug can print out info debug warn
	// info level can print warn info
	// warn can only print warn
	// debug -> info -> warn -> error
	var level zapcore.Level
	switch loglevel {
		case "debug":
			level = zap.DebugLevel
		case "info":
			level = zap.InfoLevel
		case "error":
			level = zap.ErrorLevel
		case "warn":
			level = zap.WarnLevel
		default:
			level = zap.InfoLevel
	}

	encoderConfig := zapcore.EncoderConfig {
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "linenum",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,  // Lowercase encoder
		EncodeTime:     zapcore.ISO8601TimeEncoder,     // ISO8601 UTC time format
		EncodeDuration: zapcore.SecondsDurationEncoder, //
		EncodeCaller:   zapcore.FullCallerEncoder,      // Full path encoder
		EncodeName:     zapcore.FullNameEncoder,
	}

	// set log level
	atomicLevel := zap.NewAtomicLevel()
	atomicLevel.SetLevel(level)

	var writes = []zapcore.WriteSyncer {write}

	// if it is a development environment, also output on the console
	if level == zap.DebugLevel {
		writes = append(writes, zapcore.AddSync(os.Stdout))
	}

	core := zapcore.NewCore(
		zapcore.NewConsoleEncoder(encoderConfig),
		// zapcore.NewJSONEncoder(encoderConfig),
		zapcore.NewMultiWriteSyncer(writes...),		// print to console and file
		// write,
		level,
	)

	// open development mode, stack trace
	caller := zap.AddCaller()
	// open file and line number
	development := zap.Development()
	// set initialization fields, such as: add a server name
	filed := zap.Fields(zap.String("application", "gin-chat-svc"))
	// construction log
	Logger = zap.New(core, caller, development, filed)
	Logger.Info("logger init success")
}