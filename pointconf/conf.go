package pointconf

import (
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"github.com/Steve-Zou/grpc-point/src/point.one/pointlog"
)

type LoggerConfig struct {
	Level      string `yaml:"level"`       //debug  info  warn  error
	Encoding   string `yaml:"encoding"`    //json or console
	CallFull   bool   `yaml:"call_full"`   //whether full call path or short path, default is short
	Filename   string `yaml:"file_name"`   //log file name
	MaxSize    int    `yaml:"max_size"`    //max size of log.(MB)
	MaxAge     int    `yaml:"max_age"`     //time to keep, (day)
	MaxBackups int    `yaml:"max_backups"` //max file numbers
	LocalTime  bool   `yaml:"local_time"`  //(default UTC)
	Compress   bool   `yaml:"compress"`    //default false
	IsTest     int    `yaml:"is_test"`
}

// 日志时间格式
// func timeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
//      enc.AppendString(t.Format("2006-01-02 15:04:05"))
// }
//NewLogger create logger by config
func (lconf *LoggerConfig) NewLogger() *pointlog.ZapLogger {
	if lconf.Filename == "" {
		logger, _ := zap.NewProduction(zap.AddCallerSkip(2))
		return pointlog.NewZapLogger(logger)
	}

	enCfg := zap.NewProductionEncoderConfig()
	if lconf.CallFull {
		enCfg.EncodeCaller = zapcore.FullCallerEncoder
	}
	encoder := zapcore.NewJSONEncoder(enCfg)
	if lconf.Encoding == "console" {
		zapcore.NewConsoleEncoder(enCfg)
	}

	//zapWriter := zapcore.
	zapWriter := zapcore.AddSync(&lumberjack.Logger{
		Filename:   lconf.Filename,
		MaxSize:    lconf.MaxSize,
		MaxAge:     lconf.MaxAge,
		MaxBackups: lconf.MaxBackups,
		LocalTime:  lconf.LocalTime,
	})

	newCore := zapcore.NewCore(encoder, zapWriter, zap.NewAtomicLevelAt(convertLogLevel(lconf.Level)))
	opts := []zap.Option{zap.ErrorOutput(zapWriter)}
	opts = append(opts, zap.AddCaller(), zap.AddCallerSkip(2))
	logger := zap.New(newCore, opts...)
	return pointlog.NewZapLogger(logger)
}

func convertLogLevel(levelStr string) (level zapcore.Level) {
	switch levelStr {
	case "debug":
		level = zap.DebugLevel
	case "info":
		level = zap.InfoLevel
	case "warn":
		level = zap.WarnLevel
	case "error":
		level = zap.ErrorLevel
	}
	return
}
