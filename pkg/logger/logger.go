package logger

import (
	"os"
	"sgin/pkg/config"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

type Logger struct {
	*zap.SugaredLogger
}

func (l *Logger) With(args ...interface{}) *Logger {
	return &Logger{l.SugaredLogger.With(args...)}
}

// WithOptions wraps zap.WithOptions and preserves *Logger type
func (l *Logger) WithOptions(opts ...zap.Option) *Logger {
	if l == nil || l.SugaredLogger == nil {
		return l
	}
	base := l.SugaredLogger.Desugar().WithOptions(opts...)
	return &Logger{base.Sugar()}
}

// Named returns a new named logger
func (l *Logger) Named(name string) *Logger {
	if l == nil || l.SugaredLogger == nil {
		return l
	}
	return &Logger{l.SugaredLogger.Named(name)}
}

// Desugar exposes the underlying zap.Logger
func (l *Logger) Desugar() *zap.Logger {
	if l == nil || l.SugaredLogger == nil {
		return zap.NewNop()
	}
	return l.SugaredLogger.Desugar()
}

// Convenience print methods (compat with std log-like interfaces)
func (l *Logger) Print(args ...interface{})   { l.SugaredLogger.Info(args...) }
func (l *Logger) Println(args ...interface{}) { l.SugaredLogger.Info(args...) }
func (l *Logger) Debug(args ...interface{})   { l.SugaredLogger.Debug(args...) }
func (l *Logger) Info(args ...interface{})    { l.SugaredLogger.Info(args...) }
func (l *Logger) Warn(args ...interface{})    { l.SugaredLogger.Warn(args...) }
func (l *Logger) Error(args ...interface{})   { l.SugaredLogger.Error(args...) }
func (l *Logger) DPanic(args ...interface{})  { l.SugaredLogger.DPanic(args...) }
func (l *Logger) Panic(args ...interface{})   { l.SugaredLogger.Panic(args...) }
func (l *Logger) Fatal(args ...interface{})   { l.SugaredLogger.Fatal(args...) }

func (l *Logger) Debugf(template string, args ...interface{}) {
	l.SugaredLogger.Debugf(template, args...)
}
func (l *Logger) Infof(template string, args ...interface{}) {
	l.SugaredLogger.Infof(template, args...)
}
func (l *Logger) Warnf(template string, args ...interface{}) {
	l.SugaredLogger.Warnf(template, args...)
}
func (l *Logger) Errorf(template string, args ...interface{}) {
	l.SugaredLogger.Errorf(template, args...)
}
func (l *Logger) DPanicf(template string, args ...interface{}) {
	l.SugaredLogger.DPanicf(template, args...)
}
func (l *Logger) Panicf(template string, args ...interface{}) {
	l.SugaredLogger.Panicf(template, args...)
}
func (l *Logger) Fatalf(template string, args ...interface{}) {
	l.SugaredLogger.Fatalf(template, args...)
}

func (l *Logger) Debugw(msg string, keysAndValues ...interface{}) {
	l.SugaredLogger.Debugw(msg, keysAndValues...)
}
func (l *Logger) Infow(msg string, keysAndValues ...interface{}) {
	l.SugaredLogger.Infow(msg, keysAndValues...)
}
func (l *Logger) Warnw(msg string, keysAndValues ...interface{}) {
	l.SugaredLogger.Warnw(msg, keysAndValues...)
}
func (l *Logger) Errorw(msg string, keysAndValues ...interface{}) {
	l.SugaredLogger.Errorw(msg, keysAndValues...)
}
func (l *Logger) DPanicw(msg string, keysAndValues ...interface{}) {
	l.SugaredLogger.DPanicw(msg, keysAndValues...)
}
func (l *Logger) Panicw(msg string, keysAndValues ...interface{}) {
	l.SugaredLogger.Panicw(msg, keysAndValues...)
}
func (l *Logger) Fatalw(msg string, keysAndValues ...interface{}) {
	l.SugaredLogger.Fatalw(msg, keysAndValues...)
}

func NewLogger(cfg config.LogConfig) *Logger {
	encoder := getEncoder(cfg.Format)

	var logLevel zapcore.Level
	if err := logLevel.Set(cfg.Level); err != nil {
		logLevel = zap.InfoLevel
	}

	cores := []zapcore.Core{}

	// 文件输出（当配置了文件名时）
	if cfg.Filename != "" {
		fileWS := getLogWriter(cfg)
		cores = append(cores, zapcore.NewCore(encoder, fileWS, logLevel))
	}

	// 控制台输出（可与文件并存）
	if cfg.ShowConsole || cfg.Filename == "" {
		consoleEncoder := zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())
		consoleWS := zapcore.Lock(os.Stdout)
		cores = append(cores, zapcore.NewCore(consoleEncoder, consoleWS, logLevel))
	}

	var core zapcore.Core
	if len(cores) == 1 {
		core = cores[0]
	} else {
		core = zapcore.NewTee(cores...)
	}

	// 可选日志采样
	if cfg.EnableSampling {
		initial := cfg.SamplingInitial
		thereafter := cfg.SamplingThereafter
		if initial <= 0 {
			initial = 100
		}
		if thereafter <= 0 {
			thereafter = 100
		}
		core = zapcore.NewSamplerWithOptions(core, time.Second, initial, thereafter)
	}

	// 选择堆栈等级
	var opts []zap.Option
	opts = append(opts, zap.AddCaller())
	if stLevel, err := parseLevel(cfg.StacktraceLevel); err == nil {
		opts = append(opts, zap.AddStacktrace(stLevel))
	}

	logger := zap.New(core, opts...)
	return &Logger{logger.Sugar()}
}

func getEncoder(format string) zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeDuration = zapcore.SecondsDurationEncoder

	switch format {
	case "json":
		return zapcore.NewJSONEncoder(encoderConfig)
	case "console":
		return zapcore.NewConsoleEncoder(encoderConfig)
	default:
		return zapcore.NewJSONEncoder(encoderConfig)
	}
}

func getLogWriter(config config.LogConfig) zapcore.WriteSyncer {
	lumberJackLogger := &lumberjack.Logger{
		Filename:   config.Filename,
		MaxSize:    config.MaxSize, // megabytes
		MaxBackups: 5,
		MaxAge:     config.MaxAge,   //days
		Compress:   config.Compress, // disabled by default
		LocalTime:  true,
	}

	return zapcore.AddSync(lumberJackLogger)
}

// Printf 实现了 gorm.io/gorm/logger.Writer 接口的方法
func (l *Logger) Printf(format string, args ...interface{}) {
	l.Infof(format, args...)
}

// Write 实现了 io.Writer 接口的方法
func (l *Logger) Write(p []byte) (n int, err error) {
	l.Info(string(p))
	return len(p), nil
}

func parseLevel(s string) (zapcore.Level, error) {
	if s == "" {
		return zapcore.ErrorLevel, nil
	}
	var lv zapcore.Level
	err := lv.Set(s)
	return lv, err
}

// Sync 将缓冲区刷新到下层写入器
func (l *Logger) Sync() error {
	if l == nil || l.SugaredLogger == nil {
		return nil
	}
	return l.SugaredLogger.Sync()
}
