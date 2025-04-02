package loggers

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func InitLogger() (*zap.Logger, error) {
	config := zap.NewProductionConfig() //Создает новую конфигурацию журнала с настройками
	// по умолчанию для производственной среды.
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder //установка формата времени

	config.OutputPaths = []string{"./logger/logger.log"} //установка места где будет записываться соо.

	return config.Build()
}
