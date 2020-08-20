package logger

import (
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"sync"
)

var (
	Log      *zap.Logger
	onceInit sync.Once
)

func getLogWriter(filenamewithPath string) zapcore.WriteSyncer {
	lumberJackLogger := &lumberjack.Logger{
		Filename:   filenamewithPath,
		MaxSize:    1,
		MaxBackups: 5,
		MaxAge:     30,
		Compress:   false,
	}
	return zapcore.AddSync(lumberJackLogger)
}

func getEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	return zapcore.NewJSONEncoder(encoderConfig)
}

func GetLogLevel(loglevelString string) int {
	switch loglevelString {
	case "INFO":
		return 0
	case "WARN":
		return 1
	case "ERROR":
		return 2
	case "DPANIC":
		return 3
	case "PANIC":
		return 4
	default:
		return -1
	}
	return -1

}
func InitLogger(lvl int, filenamewithPath string) error {

	var err error
	onceInit.Do(func() {
		writeSyncher := zapcore.NewMultiWriteSyncer(getLogWriter(filenamewithPath), os.Stdout)
		//		writeSyncher := getLogWriter(filenamewithPath)

		ecfg := getEncoder()

		core := zapcore.NewCore(ecfg, writeSyncher, zapcore.Level(lvl))

		Log = zap.New(core, zap.AddCaller())
		zap.RedirectStdLog(Log)
	})
	return err
}
