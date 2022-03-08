package main

import (
	"development/zap-demo/log"
)

func main() {
	log.Info("this is log test.")
	log.Error("log error test")
	log.Info("test fields for info.")
	log.Info("test2 for info")
	log.Info("creater:%s","lihaiping1603")

	logger := log.NewFileLoger("./zap_demo.log", log.InfoLevel, log.WithCaller(true), log.AddCallerSkip(1))
	logger.Info("this is test for zap demo log")
	logger.Info("creater:%s","lihaiping1603")

	logger2 := logger.With(log.String("zap-demo", "test"))
	logger2.Info("with field test")


	logger.Info("end test")

}
