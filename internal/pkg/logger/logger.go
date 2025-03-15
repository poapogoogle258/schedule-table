package logger

import (
	"bytes"
	"os"

	"encoding/json"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	empty = ""
	tab   = "\t"
)

func PrettyJson(data interface{}) (string, error) {
	buffer := new(bytes.Buffer)
	encoder := json.NewEncoder(buffer)
	encoder.SetIndent(empty, tab)

	err := encoder.Encode(data)
	if err != nil {
		return empty, err
	}
	return buffer.String(), nil
}

type fnHook func(data map[string]interface{}) error
type hook struct {
	message string
	fn      fnHook
}

var logger *zap.Logger
var file *os.File
var optionHooks = make([]*hook, 0)

func mustOpenFile(name string, flag int, perm os.FileMode) *os.File {
	_file, err := os.OpenFile(name, flag, perm)
	if err != nil {
		panic(err)
	}

	return _file
}

func fieldsToMap(fields ...zap.Field) map[string]interface{} {
	jsonFields := make(map[string]interface{})
	for _, field := range fields {
		jsonFields[field.Key] = field.String
	}

	return jsonFields
}

func Debug(message string, conds ...any) {
	jsonString, _ := PrettyJson(conds)
	logger.Debug(message + " " + jsonString)
}

func Info(message string, fields ...zap.Field) {
	logger.Info(message, fields...)
}

func Warn(message string, fields ...zap.Field) {
	logger.Warn(message, fields...)
}

func Error(message string, fields ...zap.Field) {
	logger.Error(message, fields...)
}

func Fatal(message string, fields ...zap.Field) {
	logger.Fatal(message, fields...)
}

func Message(message string, fields ...zap.Field) {
	logger.Info(message, fields...)

	// check hook
	for _, hook := range optionHooks {
		if message == hook.message {
			data := fieldsToMap(fields...)
			hook.fn(data)
			return
		}
	}
}

// check massage match with hooks
func AddHook(message string, callback fnHook) {
	optionHooks = append(optionHooks, &hook{
		message: message,
		fn:      callback,
	})
}

func InitLogger() {
	fileError := mustOpenFile("../../error.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	// File encoder config (full JSON logging)
	fileEncoderConfig := zap.NewProductionEncoderConfig()
	fileEncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	// Console encoder config (minimal output)
	consoleEncodeConfig := zap.NewDevelopmentEncoderConfig()
	consoleEncodeConfig.TimeKey = "time"
	consoleEncodeConfig.LevelKey = "level"
	consoleEncodeConfig.NameKey = ""       // disable logger name
	consoleEncodeConfig.CallerKey = ""     // disable caller
	consoleEncodeConfig.FunctionKey = ""   // disable function name
	consoleEncodeConfig.StacktraceKey = "" // disable stacktrace
	consoleEncodeConfig.MessageKey = "msg"
	consoleEncodeConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	consoleEncodeConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	highPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zap.ErrorLevel
	})

	ShowLogLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zap.InfoLevel
	})

	fileEncoder := zapcore.NewJSONEncoder(fileEncoderConfig)
	consoleEncoder := zapcore.NewConsoleEncoder(consoleEncodeConfig)

	core := zapcore.NewTee(
		zapcore.NewCore(fileEncoder, zapcore.AddSync(fileError), highPriority),
		zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), ShowLogLevel),
	)

	// Create logger with development options
	logger = zap.New(core,
		zap.AddStacktrace(zap.ErrorLevel), // Add stack trace for file logging
	)
}

func Sync() {
	logger.Sync()
	file.Close()
}
