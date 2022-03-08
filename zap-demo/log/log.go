package log

import (
	"io"
	"time"

	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Option = zap.Option

var (
	WithCaller    = zap.WithCaller
	AddStacktrace = zap.AddStacktrace
	AddCallerSkip = zap.AddCallerSkip
)

type Level = zapcore.Level

const (
	DebugLevel  Level = zap.DebugLevel  // -1
	InfoLevel   Level = zap.InfoLevel   // 0, default level
	WarnLevel   Level = zap.WarnLevel   // 1
	ErrorLevel  Level = zap.ErrorLevel  // 2
	DPanicLevel Level = zap.DPanicLevel // 3, used in development log
	// PanicLevel logs a message, then panics
	PanicLevel Level = zap.PanicLevel // 4
	// FatalLevel logs a message, then calls os.Exit(1).
	FatalLevel Level = zap.FatalLevel // 5
)

// function variables for all field types
// in github.com/uber-go/zap/field.go
type Field = zap.Field

var (
	Skip        = zap.Skip
	Binary      = zap.Binary
	Bool        = zap.Bool
	Boolp       = zap.Boolp
	ByteString  = zap.ByteString
	Complex128  = zap.Complex128
	Complex128p = zap.Complex128p
	Complex64   = zap.Complex64
	Complex64p  = zap.Complex64p
	Float64     = zap.Float64
	Float64p    = zap.Float64p
	Float32     = zap.Float32
	Float32p    = zap.Float32p
	Int         = zap.Int
	Intp        = zap.Intp
	Int64       = zap.Int64
	Int64p      = zap.Int64p
	Int32       = zap.Int32
	Int32p      = zap.Int32p
	Int16       = zap.Int16
	Int16p      = zap.Int16p
	Int8        = zap.Int8
	Int8p       = zap.Int8p
	String      = zap.String
	Stringp     = zap.Stringp
	Uint        = zap.Uint
	Uintp       = zap.Uintp
	Uint64      = zap.Uint64
	Uint64p     = zap.Uint64p
	Uint32      = zap.Uint32
	Uint32p     = zap.Uint32p
	Uint16      = zap.Uint16
	Uint16p     = zap.Uint16p
	Uint8       = zap.Uint8
	Uint8p      = zap.Uint8p
	Uintptr     = zap.Uintptr
	Uintptrp    = zap.Uintptrp
	Reflect     = zap.Reflect
	Namespace   = zap.Namespace
	Stringer    = zap.Stringer
	Time        = zap.Time
	Timep       = zap.Timep
	Stack       = zap.Stack
	StackSkip   = zap.StackSkip
	Duration    = zap.Duration
	Durationp   = zap.Durationp
	Any         = zap.Any
)

type Logger struct {
	log   *zap.SugaredLogger // zap ensure that zap.Logger is safe for concurrent use
	level Level
}

// New create a new logger (not support log rotating).
func NewLogger(writer io.Writer, level Level, opts ...Option) *Logger {
	if writer == nil {
		panic("the writer is nil")
	}
	cfg := zap.NewProductionConfig()
	cfg.EncoderConfig.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.Format("2006-01-02 15:04:05.000"))
	}
	//enc := zapcore.NewJSONEncoder(cfg.EncoderConfig)
	enc := zapcore.NewConsoleEncoder(cfg.EncoderConfig)
	core := zapcore.NewCore(
		enc,
		zapcore.AddSync(writer),
		zapcore.Level(level),
	)
	zap_logger:=zap.New(core, opts...)
	logger := &Logger{
		log:   zap_logger.Sugar(),
		level: level,
	}
	return logger
}

func (l *Logger) Debug(fmt string,args ...interface{}) {
	l.log.Debugf(fmt,args...)
}

func (l *Logger) Info(fmt string,args ...interface{}) {
	l.log.Infof(fmt,args...)
}

func (l *Logger) Warn(fmt string,args ...interface{}) {
	l.log.Warnf(fmt,args...)
}

func (l *Logger) Error(fmt string,args ...interface{}) {
	l.log.Errorf(fmt,args...)
}

func (l *Logger) DPanic(fmt string,args ...interface{}) {
	l.log.DPanicf(fmt,args...)
}

func (l *Logger) Panic(fmt string,args ...interface{}) {
	l.log.Panicf(fmt,args...)
}

func (l *Logger) Fatal(fmt string,args ...interface{}) {
	l.log.Fatalf(fmt,args...)
}

func (l *Logger) With(args ...interface{}) *Logger {
	log := l.log.With(args...)
	l.log = log
	return l
}

func (l *Logger) Sync() error {
	return l.log.Sync()
}

type RotateOptions struct {
	MaxSize    int
	MaxAge     int
	MaxBackups int
	Compress   bool
}

type LevelEnablerFunc func(lvl Level) bool

type TeeOption struct {
	Filename string
	Ropt     RotateOptions
	Lef      LevelEnablerFunc
}

func NewLoggerTeeWithRotate(tops []TeeOption, opts ...Option) *Logger {
	var cores []zapcore.Core
	cfg := zap.NewProductionConfig()
	cfg.EncoderConfig.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.Format("2006-01-02T15:04:05.000Z0700"))
	}
	//enc_cfg := zapcore.NewJSONEncoder(cfg.EncoderConfig)
	enc_cfg := zapcore.NewConsoleEncoder(cfg.EncoderConfig)
	for _, top := range tops {
		top := top

		lv := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
			return top.Lef(Level(lvl))
		})

		w := zapcore.AddSync(&lumberjack.Logger{
			Filename:   top.Filename,
			MaxSize:    top.Ropt.MaxSize,
			MaxBackups: top.Ropt.MaxBackups,
			MaxAge:     top.Ropt.MaxAge,
			Compress:   top.Ropt.Compress,
		})

		core := zapcore.NewCore(
			enc_cfg,
			zapcore.AddSync(w),
			lv,
		)
		cores = append(cores, core)
	}
	zap_logger := zap.New(zapcore.NewTee(cores...), opts...)
	logger := &Logger{
		log: zap_logger.Sugar(),
	}
	return logger
}
