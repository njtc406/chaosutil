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

// NewDefaultLogger 创建一个通用日志对象
// filePath 日志文件名(最终文件名会是 filePath_20060102150405.log)(filePath为空且开启标准输出的情况下默认输出到stdout,否则无任何输出)
// maxAge 最大存放时间(过期会自动删除)
// rotationTime 自动切分间隔(到期日志自动切换新文件)
// level 日志级别(小于设置级别的信息都会被记录打印)
// withCaller 是否需要调用者信息
// fullCaller 如果需要打印调用者信息,那么这个参数可以设置调用者信息的详细程度
// withColors 是否需要信息的颜色(基本上只能用于linux的前台打印)
// openStdout 是否开启标准输出(如果filePath为空,且openStdout未开启,那么将不会有任何日志信息被记录)
func NewDefaultLogger(filePath string, maxAge, rotationTime time.Duration, level uint32, withCaller, fullCaller, withColors, openStdout bool) *DefaultLogger {
	logger := &DefaultLogger{}
	if len(filePath) > 0 {
		if rotationTime < time.Second*60 || rotationTime > time.Hour*24 {
			panic("rotationTime must >= 1min and <= 24hour")
		}
		pattern := "_%Y%m%d.log"
		if rotationTime < time.Minute*60 {
			pattern = "_%Y%m%d%H%M.log"
		} else if rotationTime < time.Hour*24 {
			pattern = "_%Y%m%d%H.log"
		}
		if w, err := rotateNew(
			filePath,
			WithMaxAge(maxAge),
			WithRotationTime(rotationTime),
			WithPattern(pattern),
		); err != nil {
			panic(err)
		} else {
			logger.Writer = w
		}
	} else {
		logger.Writer = os.Stdout
	}

	if level > 6 {
		panic("log level must <= 6")
	}

	logger.Logger = New(
		WithLevel(logrus.Level(level)),
		WithCaller(withCaller),
		WithColor(withColors),
		WithOut(logger.Writer),
		WithFullCaller(fullCaller),
	)
	if openStdout && len(filePath) > 0 {
		logger.Logger.SetOutput(io.MultiWriter(os.Stdout, logger.Writer))
	}
	// 由于是追加模式,所以默认为无锁
	logger.Logger.SetNoLock()

	return logger
}

// Close 释放日志对象
func (d *DefaultLogger) Close() {
	d.Logger.SetOutput(os.Stdout)
	if err := d.Writer.Close(); err != nil {
		_, _ = fmt.Fprintln(os.Stdout, err)
	}
	d.Writer.Close()
	d.Writer = nil
	d.Logger = nil
}
