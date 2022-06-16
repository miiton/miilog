package miilog

import "testing"

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
