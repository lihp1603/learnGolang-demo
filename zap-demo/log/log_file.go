package log

import "os"

func NewFileLoger(logpath string, level Level, opts ...Option) *Logger {
	file, err := os.OpenFile(logpath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	return NewLogger(file, level, opts...)
}
