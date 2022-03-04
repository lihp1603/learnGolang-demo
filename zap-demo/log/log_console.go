package log

import "os"

var std = NewLogger(os.Stderr, InfoLevel, WithCaller(true), AddCallerSkip(1))

var (
	Info   = std.Info
	Warn   = std.Warn
	Error  = std.Error
	DPanic = std.DPanic
	Panic  = std.Panic
	Fatal  = std.Fatal
	Debug  = std.Debug
)

func Default() *Logger {
	return std
}

// not safe for concurrent use
//func ResetDefault(l *Logger) {
//	std = l
//	Info = std.Info
//	Warn = std.Warn
//	Error = std.Error
//	DPanic = std.DPanic
//	Panic = std.Panic
//	Fatal = std.Fatal
//	Debug = std.Debug
//}

func Sync() error {
	if std != nil {
		return std.Sync()
	}
	return nil
}
