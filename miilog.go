package miilog

import (
	"net"
	"net/http"
	"net/url"
	"os"
	"path"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// SetLoggerProductionWithLokiMust set a logger for global.
// Output: stdout with JSON and a Grafana Loki Server using protocol buffers.
func SetLoggerProductionWithLokiMust(lokiURL, tenantID, labels string) {
	client := &http.Client{
		Transport: &http.Transport{
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
		},
		Timeout: 5 * time.Second,
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
			zapcore.NewJSONEncoder(
				zapcore.EncoderConfig{
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
				},
			),
			zapcore.NewMultiWriteSyncer(os.Stdout, lokiSyncer),
			zap.NewAtomicLevelAt(zap.InfoLevel),
		),
		zap.AddCaller(),
		zap.AddStacktrace(zap.ErrorLevel),
	)
	zap.ReplaceGlobals(logger)
	setWrappers()
}

// SetLoggerProductionMust set a logger for global.
// Output: stdout with JSON.
func SetLoggerProductionMust() {
	cfg := zap.Config{
		Level:       zap.NewAtomicLevelAt(zap.InfoLevel),
		Development: false,
		Sampling: &zap.SamplingConfig{
			Initial:    100,
			Thereafter: 100,
		},
		Encoding: "json",
		EncoderConfig: zapcore.EncoderConfig{
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
		},
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
	}
	logger, err := cfg.Build()
	if err != nil {
		panic(err)
	}
	zap.ReplaceGlobals(logger)
	setWrappers()
}

// SetLoggerDevelopmentMust set a logger for global.
// Output: stdout with TextLine.
func SetLoggerDevelopmentMust() {
	cfg := zap.Config{
		Level:       zap.NewAtomicLevelAt(zap.InfoLevel),
		Development: false,
		Sampling: &zap.SamplingConfig{
			Initial:    100,
			Thereafter: 100,
		},
		Encoding: "console",
		EncoderConfig: zapcore.EncoderConfig{
			TimeKey:        "ts",
			LevelKey:       "level",
			NameKey:        "logger",
			CallerKey:      "caller",
			FunctionKey:    zapcore.OmitKey,
			MessageKey:     "msg",
			StacktraceKey:  "stacktrace",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.CapitalLevelEncoder,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.SecondsDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		},
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
	}
	logger, err := cfg.Build()
	if err != nil {
		panic(err)
	}
	zap.ReplaceGlobals(logger)
	setWrappers()
}

var (
	Sync   func() error
	Debug  func(args ...interface{})
	Debugw func(msg string, keysAndValues ...interface{})
	Debugf func(template string, args ...interface{})
	Info   func(args ...interface{})
	Infow  func(msg string, keysAndValues ...interface{})
	Infof  func(template string, args ...interface{})
	Warn   func(args ...interface{})
	Warnw  func(msg string, keysAndValues ...interface{})
	Warnf  func(template string, args ...interface{})
	Error  func(args ...interface{})
	Errorw func(msg string, keysAndValues ...interface{})
	Errorf func(template string, args ...interface{})
	Fatal  func(args ...interface{})
	Fatalw func(msg string, keysAndValues ...interface{})
	Fatalf func(template string, args ...interface{})
)

func setWrappers() {
	Sync = zap.S().Sync
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
}
