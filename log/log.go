package log

import (
	"bytes"
	"github.com/sirupsen/logrus"
	"io"
	"os"
	"strings"
	"sync"
	"time"
)

// TODO 差一个文件日志,这个日志是用来写入一些统计日志的,所以格式上可能会和其他不太一样,只需要数据,不需要附加信息
// TODO 同时可以增加一个文件日志切分后的自动上传,远端收到文件直接分析文件内容

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

const (
	PanicLevelStr = "panic"
	FatalLevelStr = "fatal"
	ErrorLevelStr = "error"
	WarnLevelStr  = "warn"
	InfoLevelStr  = "info"
	DebugLevelStr = "debug"
	TraceLevelStr = "trace"
)

var levelMap = map[string]logrus.Level{
	PanicLevelStr: PanicLevel,
	FatalLevelStr: FatalLevel,
	ErrorLevelStr: ErrorLevel,
	WarnLevelStr:  WarnLevel,
	InfoLevelStr:  InfoLevel,
	DebugLevelStr: DebugLevel,
	TraceLevelStr: TraceLevel,
}

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

func WithHook(hook logrus.Hook) Option {
	return func(logger *Logger) {
		logger.AddHook(hook)
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

// NewDefaultLogger 创建一个通用日志对象
// filePath 日志文件名(最终文件名会是 filePath_20060102150405.log)(filePath为空且开启标准输出的情况下默认输出到stdout,否则无任何输出)
// maxAge 最大存放时间(过期会自动删除)
// rotationTime 自动切分间隔(到期日志自动切换新文件)
// level 日志级别(小于设置级别的信息都会被记录打印,设置级别如果超出限制,默认日志等级为error,取值为panic fatal error warn info debug trace)
// withCaller 是否需要调用者信息
// fullCaller 如果需要打印调用者信息,那么这个参数可以设置调用者信息的详细程度
// withColors 是否需要信息的颜色(基本上只能用于linux的前台打印)
// openStdout 是否开启标准输出(如果filePath为空,且openStdout未开启,那么将不会有任何日志信息被记录)
// TODO 支持远程日志钩子函数,可以将日志发送到远程的日志记录器上(这个函数需要go出去执行,不能阻塞)
func NewDefaultLogger(filePath string, maxAge, rotationTime time.Duration, level string, withCaller, fullCaller, withColors, openStdout bool) (*Logger, error) {
	var writers []io.Writer
	if len(filePath) > 0 {
		if rotationTime < time.Second*60 || rotationTime > time.Hour*24 {
			return nil, DefaultRotationTimeErr
		}
		pattern := "_%Y%m%d.log"
		if rotationTime < time.Minute*60 {
			pattern = "_%Y%m%d%H%M.log"
		} else if rotationTime < time.Hour*24 {
			pattern = "_%Y%m%d%H.log"
		}

		w, err := rotateNew(
			filePath,
			WithMaxAge(maxAge),
			WithRotationTime(rotationTime),
			WithPattern(pattern),
		)
		if err != nil {
			w.Close()
			return nil, err
		} else {
			writers = append(writers, w)
		}
	}

	if openStdout {
		writers = append(writers, os.Stdout)
	} else {
		writers = append(writers, io.Discard)
	}

	level = strings.ToLower(level)
	if _, ok := levelMap[level]; !ok {
		level = "error"
	}

	logger := New(
		WithLevel(levelMap[level]),
		WithCaller(withCaller),
		WithColor(withColors),
		WithOut(io.MultiWriter(writers...)),
		WithFullCaller(fullCaller),
	)

	// 由于是追加模式,所以默认为无锁(gpt认为这里在多线程环境中可能会产生一些问题)
	logger.SetNoLock()

	return logger, nil
}

func Release(logger *Logger) error {
	if logger == nil || logger.Writer() == nil {
		return nil
	}

	logger.SetOutput(os.Stdout)

	return logger.Writer().Close()
}
