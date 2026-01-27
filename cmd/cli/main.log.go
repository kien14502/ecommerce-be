package cli

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func main() {
	sugar := zap.NewExample().Sugar()
	sugar.Info("This is a log message from main.log.go")

	encoder := getEncoderLog()
	writerSync := getWriterSync()
	core := zapcore.NewCore(encoder, writerSync, zapcore.InfoLevel)
	logger := zap.New(core, zap.AddCaller())
	defer logger.Sync()
}

func getEncoderLog() zapcore.Encoder {
	encodeConfig := zap.NewProductionEncoderConfig()
	encodeConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encodeConfig.TimeKey = "time"
	encodeConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	encodeConfig.LevelKey = "level"
	encodeConfig.NameKey = "logger"
	encodeConfig.CallerKey = "caller"
	return zapcore.NewJSONEncoder(encodeConfig)
}

func getWriterSync() zapcore.WriteSyncer {
	file, _ := os.OpenFile("./log/log.txt", os.O_RDWR, os.ModePerm)
	syncFile := zapcore.AddSync(zapcore.Lock(file))
	syncConsole := zapcore.AddSync(os.Stdout)
	return zapcore.NewMultiWriteSyncer(syncFile, syncConsole)
}
