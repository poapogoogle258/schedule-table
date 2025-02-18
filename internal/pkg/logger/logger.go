package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

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

func Debug(message string, fields ...zap.Field) {
	logger.Debug(message, fields...)
}

func Info(message string, fields ...zap.Field) {
	logger.Info(message, fields...)
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
	file := mustOpenFile("../../info.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	highPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zap.ErrorLevel
	})

	allPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return true
	})

	fileEncoder := zapcore.NewJSONEncoder(encoderConfig)
	consoleEncoder := zapcore.NewConsoleEncoder(encoderConfig)

	core := zapcore.NewTee(
		zapcore.NewCore(fileEncoder, zapcore.AddSync(file), highPriority),
		zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), allPriority),
	)

	logger = zap.New(core)

}

func Sync() {
	logger.Sync()
	file.Close()
}
