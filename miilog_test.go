package miilog

import (
	"testing"
)

func Test_SetLoggerProductionMust(t *testing.T) {
	SetLoggerProductionMust()
	defer Sync()
	Debug("hoge")
	Info("hoge")
	Warn("hoge")
	Error("hoge")
}

func Test_SetLoggerDevelopmentMust(t *testing.T) {
	SetLoggerDevelopmentMust()
	defer Sync()
	Debug("hoge")
	Debug("hoge")
	Debug("hoge")
	Debug("hoge")
	Debug("hoge")
	Info("hoge")
	Warn("hoge")
	Error("hoge")
}

func Test_SetLoggerProductionWithFileAndLokiMust(t *testing.T) {
	SetLoggerProductionWithFileAndLokiMust("tmp/test.log", "http://localhost:3100", "MYAWESOMETENANT", "{host=\"localhost\", job=\"gotest\"}")
	defer Sync()
	Debug("hoge")
	Info("hoge")
	Warn("hoge")
	Error("hoge")
}

func Test_SetLoggerProductionWithLokiMust(t *testing.T) {
	SetLoggerProductionWithLokiMust("http://localhost:3100", "MYAWESOMETENANT", "{host=\"localhost\", job=\"gotest\"}")
	defer Sync()
	Info("foo")
	Info("bar")
	Info("baz")
	Info("hoge")
	Info("fuga")
	Info("piyo")
	Info("foobar")
	Info("bazbaz")
	Info("hogehoge")
	Info("fugafuga")
	Info("piyopiyo")
	Info("hogefuga")
	Info("hogepiyo")
	Info("fugapiyo")
	Info("hogefugapiyo")
}
