package miilog

import (
	"net"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var logLevel zapcore.Level

var jsonEncoderConfig = zapcore.EncoderConfig{
	TimeKey:        "ts",
	LevelKey:       "level",
	NameKey:        "logger",
	CallerKey:      "caller",
	FunctionKey:    zapcore.OmitKey,
	MessageKey:     "msg",
	StacktraceKey:  "stacktrace",
	LineEnding:     zapcore.DefaultLineEnding,
	EncodeLevel:    zapcore.LowercaseLevelEncoder,
	EncodeTime:     zapcore.RFC3339NanoTimeEncoder,
	EncodeDuration: zapcore.SecondsDurationEncoder,
	EncodeCaller:   zapcore.ShortCallerEncoder,
}

var consoleEncoderConfig = zapcore.EncoderConfig{
	TimeKey:          "ts",
	LevelKey:         "level",
	NameKey:          "logger",
	CallerKey:        "caller",
	FunctionKey:      zapcore.OmitKey,
	MessageKey:       "msg",
	StacktraceKey:    "stacktrace",
	LineEnding:       zapcore.DefaultLineEnding,
	EncodeLevel:      zapcore.CapitalColorLevelEncoder,
	EncodeTime:       zapcore.TimeEncoderOfLayout("15:04:05.000000"),
	EncodeDuration:   zapcore.SecondsDurationEncoder,
	EncodeCaller:     zapcore.ShortCallerEncoder,
	ConsoleSeparator: " ",
}

var transportConfig = &http.Transport{
	DialContext: (&net.Dialer{
		Timeout:   5 * time.Second,
		KeepAlive: 30 * time.Second,
	}).DialContext,
	MaxIdleConns:          100,
	MaxConnsPerHost:       100,
	IdleConnTimeout:       30 * time.Second,
	TLSHandshakeTimeout:   5 * time.Second,
	ResponseHeaderTimeout: 5 * time.Second,
	ExpectContinueTimeout: 1 * time.Second,
}

// SetLoggerProductionWithLokiMust set a logger for global.
// Output: stdout with JSON and a Grafana Loki Server using protocol buffers.
func SetLoggerProductionWithLokiMust(lokiURL, tenantID, labels string) {
	client := &http.Client{
		Transport: transportConfig,
		Timeout:   5 * time.Second,
	}
	u, err := url.Parse(lokiURL)
	if err != nil {
		panic(err)
	}
	u.Path = path.Join(u.Path, "loki", "api", "v1", "push")
	lokiSyncer := &LokiSyncer{
		URL:      u.String(),
		TenantID: tenantID,
		Labels:   labels,
		Client:   client,
	}
	logger := zap.New(
		zapcore.NewCore(
			zapcore.NewJSONEncoder(jsonEncoderConfig),
			zapcore.NewMultiWriteSyncer(os.Stdout, lokiSyncer),
			zap.NewAtomicLevelAt(logLevel),
		),
		zap.AddCaller(),
		zap.AddStacktrace(zap.ErrorLevel),
	)
	zap.ReplaceGlobals(logger)
	SetWrappers()
}

// SetLoggerProductionWithFileAndLokiMust set a logger for global.
// Output: stdout with JSONFile and a Grafana Loki Server using protocol buffers.
func SetLoggerProductionWithFileAndLokiMust(filePath, lokiURL, tenantID, labels string) {
	f, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		panic(err)
	}
	client := &http.Client{
		Transport: transportConfig,
		Timeout:   5 * time.Second,
	}
	u, err := url.Parse(lokiURL)
	if err != nil {
		panic(err)
	}
	u.Path = path.Join(u.Path, "loki", "api", "v1", "push")
	lokiSyncer := &LokiSyncer{
		URL:      u.String(),
		TenantID: tenantID,
		Labels:   labels,
		Client:   client,
	}
	logger := zap.New(
		zapcore.NewCore(
			zapcore.NewJSONEncoder(jsonEncoderConfig),
			zapcore.NewMultiWriteSyncer(os.Stdout, zapcore.AddSync(f), lokiSyncer),
			zap.NewAtomicLevelAt(logLevel),
		),
		zap.AddCaller(),
		zap.AddStacktrace(zap.ErrorLevel),
	)
	zap.ReplaceGlobals(logger)
	SetWrappers()
}

// SetLoggerProductionMust set a logger for global.
// Output: stdout with JSON.
func SetLoggerProductionMust() {
	cfg := zap.Config{
		Level:       zap.NewAtomicLevelAt(logLevel),
		Development: false,
		Sampling: &zap.SamplingConfig{
			Initial:    100,
			Thereafter: 100,
		},
		Encoding:         "json",
		EncoderConfig:    jsonEncoderConfig,
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
	}
	logger, err := cfg.Build()
	if err != nil {
		panic(err)
	}
	zap.ReplaceGlobals(logger)
	SetWrappers()
}

// SetLoggerDevelopmentMust set a logger for global.
// Output: stdout with TextLine.
func SetLoggerDevelopmentMust() {
	cfg := zap.Config{
		Level:       zap.NewAtomicLevelAt(logLevel),
		Development: true,
		Sampling: &zap.SamplingConfig{
			Initial:    100,
			Thereafter: 100,
		},
		Encoding:         "console",
		EncoderConfig:    consoleEncoderConfig,
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
	}
	logger, err := cfg.Build()
	if err != nil {
		panic(err)
	}
	zap.ReplaceGlobals(logger)
	SetWrappers()
}

var (
	Sync    func() error
	Debug   func(args ...interface{})
	Debugw  func(msg string, keysAndValues ...interface{})
	Debugf  func(template string, args ...interface{})
	Info    func(args ...interface{})
	Infow   func(msg string, keysAndValues ...interface{})
	Infof   func(template string, args ...interface{})
	Warn    func(args ...interface{})
	Warnw   func(msg string, keysAndValues ...interface{})
	Warnf   func(template string, args ...interface{})
	Error   func(args ...interface{})
	Errorw  func(msg string, keysAndValues ...interface{})
	Errorf  func(template string, args ...interface{})
	Fatal   func(args ...interface{})
	Fatalw  func(msg string, keysAndValues ...interface{})
	Fatalf  func(template string, args ...interface{})
	Panic   func(args ...interface{})
	Panicw  func(msg string, keysAndValues ...interface{})
	Panicf  func(template string, args ...interface{})
	DPanic  func(args ...interface{})
	DPanicw func(msg string, keysAndValues ...interface{})
	DPanicf func(template string, args ...interface{})
)

func SetWrappers(args ...interface{}) {
	Sync = zap.S().Sync
	if len(args) > 0 {
		Debug = zap.S().With(args...).Debug
		Debugw = zap.S().With(args...).Debugw
		Debugf = zap.S().With(args...).Debugf
		Info = zap.S().With(args...).Info
		Infow = zap.S().With(args...).Infow
		Infof = zap.S().With(args...).Infof
		Warn = zap.S().With(args...).Warn
		Warnw = zap.S().With(args...).Warnw
		Warnf = zap.S().With(args...).Warnf
		Error = zap.S().With(args...).Error
		Errorw = zap.S().With(args...).Errorw
		Errorf = zap.S().With(args...).Errorf
		Fatal = zap.S().With(args...).Fatal
		Fatalw = zap.S().With(args...).Fatalw
		Fatalf = zap.S().With(args...).Fatalf
		Panic = zap.S().With(args...).Panic
		Panicw = zap.S().With(args...).Panicw
		Panicf = zap.S().With(args...).Panicf
		DPanic = zap.S().With(args...).DPanic
		DPanicw = zap.S().With(args...).DPanicw
		DPanicf = zap.S().With(args...).DPanicf
		return
	}
	Debug = zap.S().Debug
	Debugw = zap.S().Debugw
	Debugf = zap.S().Debugf
	Info = zap.S().Info
	Infow = zap.S().Infow
	Infof = zap.S().Infof
	Warn = zap.S().Warn
	Warnw = zap.S().Warnw
	Warnf = zap.S().Warnf
	Error = zap.S().Error
	Errorw = zap.S().Errorw
	Errorf = zap.S().Errorf
	Fatal = zap.S().Fatal
	Fatalw = zap.S().Fatalw
	Fatalf = zap.S().Fatalf
	Panic = zap.S().Panic
	Panicw = zap.S().Panicw
	Panicf = zap.S().Panicf
	DPanic = zap.S().DPanic
	DPanicw = zap.S().DPanicw
	DPanicf = zap.S().DPanicf
}

func init() {
	if strings.ToLower(os.Getenv("LOG_LEVEL")) == "debug" {
		logLevel = zap.DebugLevel
	} else {
		logLevel = zap.InfoLevel
	}
}
