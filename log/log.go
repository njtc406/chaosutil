package log

import (
	"bytes"
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"os"
	"sync"
	"time"
)

//type Logger = logrus.Logger

type Logger = logrus.Logger

type Fields = logrus.Fields

type Level = logrus.Level

// These are the different logging levels. You can set the logging level to log
// on your instance of rusLogger, obtained with `logrus.New()`.
const (
	// PanicLevel level, the highest level of severity. Logs and then calls panic with the
	// message passed to Debug, Info, ...
	PanicLevel Level = logrus.PanicLevel
	// FatalLevel level. Logs and then calls `Exit(1)`. It will exit even if the
	// logging level is set to Panic.
	FatalLevel Level = logrus.FatalLevel
	// ErrorLevel level. Logs. Used for errors that should definitely be noted.
	// Commonly used for hooks to send errors to an error tracking service.
	ErrorLevel Level = logrus.ErrorLevel
	// WarnLevel level. Non-critical entries that deserve eyes.
	WarnLevel Level = logrus.WarnLevel
	// InfoLevel level. General operational entries about what's going on inside the
	// application.
	InfoLevel Level = logrus.InfoLevel
	// DebugLevel level. Usually only enabled when debugging. Very verbose logging.
	DebugLevel Level = logrus.DebugLevel
	// TraceLevel level. Designates finer-grained informational events than the Debug.
	TraceLevel Level = logrus.TraceLevel
)

type Option func(*Logger)

// WithLevel sets the rusLogger level.
func WithLevel(level Level) Option {
	return func(logger *Logger) {
		logger.SetLevel(level)
	}
}

// WithOut sets the rusLogger output.
func WithOut(output io.Writer) Option {
	return func(logger *Logger) {
		logger.SetOutput(output)
	}
}

// WithColor sets the rusLogger color.
func WithColor(color bool) Option {
	return func(logger *Logger) {
		logger.Formatter.(*Formatter).SetColors(color)
	}
}

// WithCaller sets the rusLogger caller.
func WithCaller(caller bool) Option {
	return func(logger *Logger) {
		logger.SetReportCaller(caller)
		logger.Formatter.(*Formatter).SetCallerDisable(!caller)
	}
}

// WithFullCaller sets print full caller info
func WithFullCaller(full bool) Option {
	return func(logger *Logger) {
		logger.Formatter.(*Formatter).SetFullCaller(full)
	}
}

func WithTimeFormat(formatStr string) Option {
	return func(logger *Logger) {
		logger.Formatter.(*Formatter).SetTimestampFormat(formatStr)
	}
}

// New creates a new Logger object.
func New(opts ...Option) *Logger {
	logger := logrus.New()
	//AddHook(&Hook{})
	logger.SetFormatter(&Formatter{
		Mu:              new(sync.Mutex),
		TimestampFormat: "2006-01-02 15:04:05.000",
		Colors:          false,
		TrimMessages:    false,
		NoCaller:        false,
		bufPool: &defaultPool{
			pool: &sync.Pool{
				New: func() interface{} {
					return new(bytes.Buffer)
				},
			},
		},
	})
	for _, opt := range opts {
		opt(logger)
	}

	return logger
}

type ILogger interface {
	Tracef(format string, args ...interface{})
	Debugf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Printf(format string, args ...interface{})
	Warnf(format string, args ...interface{})
	Warningf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
	Fatalf(format string, args ...interface{})
	Panicf(format string, args ...interface{})

	Trace(args ...interface{})
	Debug(args ...interface{})
	Info(args ...interface{})
	Print(args ...interface{})
	Warn(args ...interface{})
	Warning(args ...interface{})
	Error(args ...interface{})
	Fatal(args ...interface{})
	Panic(args ...interface{})

	TraceFn(fn logrus.LogFunction)
	DebugFn(fn logrus.LogFunction)
	InfoFn(fn logrus.LogFunction)
	PrintFn(fn logrus.LogFunction)
	WarnFn(fn logrus.LogFunction)
	WarningFn(fn logrus.LogFunction)
	ErrorFn(fn logrus.LogFunction)
	FatalFn(fn logrus.LogFunction)
	PanicFn(fn logrus.LogFunction)

	Traceln(args ...interface{})
	Debugln(args ...interface{})
	Infoln(args ...interface{})
	Println(args ...interface{})
	Warnln(args ...interface{})
	Warningln(args ...interface{})
	Errorln(args ...interface{})
	Fatalln(args ...interface{})
	Panicln(args ...interface{})

	Exit(code int)
}

// DefaultLogger 通用的日志记录对象
type DefaultLogger struct {
	Logger *Logger
	Writer io.WriteCloser
}

// NewDefaultLogger 创建一个通过日志对象(filePath为空时,默认输出到stdout)
func NewDefaultLogger(filePath string, maxAge, rotationTime time.Duration, level int, withCaller, withColors, openStdout bool) *DefaultLogger {
	logger := &DefaultLogger{}
	if len(filePath) > 0 {
		if w, err := rotateNew(
			filePath,
			WithMaxAge(maxAge),
			WithRotationTime(rotationTime),
		); err != nil {
			panic(err)
		} else {
			logger.Writer = w
		}
	} else {
		logger.Writer = os.Stdout
	}

	logger.Logger = New(
		WithLevel(logrus.Level(level)),
		WithCaller(withCaller),
		WithColor(withColors),
		WithOut(logger.Writer),
	)
	if openStdout && len(filePath) > 0 {
		logger.Logger.SetOutput(io.MultiWriter(os.Stdout, logger.Writer))
	}
	logger.Logger.SetNoLock()

	return logger
}

// Close 释放日志对象
func (d *DefaultLogger) Close() {
	d.Logger.SetOutput(os.Stdout)
	if err := d.Writer.Close(); err != nil {
		_, _ = fmt.Fprintln(os.Stdout, err)
	}
	d.Writer = nil
	d.Logger = nil
}
