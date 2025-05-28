package log

import (
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"abchain_scan/config"
)

var Logger *zap.Logger

func init() {
	InitLoggerForTest()
}

func InitLoggerForTest() {
	Logger, _ = zap.NewDevelopment()
	return
}

func InitLogger() {
	if !config.G.Log.Async {
		Logger, _ = zap.NewDevelopment()
		return
	}

	buffer := &zapcore.BufferedWriteSyncer{
		Size:          config.G.Log.AsyncBufferSizeByByte,
		FlushInterval: time.Second * time.Duration(config.G.Log.AsyncFlushIntervalBySecond),
		WS:            os.Stdout,
	}
	writeSyncer := zapcore.AddSync(buffer)

	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		writeSyncer,
		zapcore.DebugLevel,
	)

	Logger = zap.New(core)
}
