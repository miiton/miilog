package miilog

import "testing"

func Test_SetLoggerProductionMust(t *testing.T) {
	SetLoggerProductionMust()
	Debug("hoge")
	Info("hoge")
	Warn("hoge")
	Error("hoge")
}

func Test_SetLoggerDevelopmentMust(t *testing.T) {
	SetLoggerDevelopmentMust()
	Debug("hoge")
	Info("hoge")
	Warn("hoge")
	Error("hoge")
}
