/*
 * Copyright (c) 2024. YR. All rights reserved
 */

// Package log
// 模块名: 模块名
// 功能描述: 描述
// 作者:  yr  2024/3/2 0002 18:57
// 最后更新:  yr  2024/3/2 0002 18:57
package log

import (
	"bytes"
	"github.com/sirupsen/logrus"
	"io"
	"os"
	"path"
	"strings"
	"sync"
	"time"
)

type CLogger struct {
	logger *logrus.Logger
}

func (c *CLogger) GetOutput() io.Writer {
	return c.logger.Out
}

func (c *CLogger) Writer() *io.PipeWriter {
	return c.logger.Writer()
}

func (c *CLogger) WriterLevel(level Level) *io.PipeWriter {
	return c.logger.WriterLevel(level)
}

func (c *CLogger) Exit(code int) {
	c.logger.Exit(code)
}

func (c *CLogger) Traceln(args ...interface{}) {
	c.logger.Traceln(args...)
}

func (c *CLogger) Debugln(args ...interface{}) {
	c.logger.Debugln(args...)
}

func (c *CLogger) Infoln(args ...interface{}) {
	c.logger.Infoln(args...)
}

func (c *CLogger) Println(args ...interface{}) {
	c.logger.Println(args...)
}

func (c *CLogger) Warnln(args ...interface{}) {
	c.logger.Warnln(args...)
}

func (c *CLogger) Warningln(args ...interface{}) {
	c.logger.Warningln(args...)
}

func (c *CLogger) Errorln(args ...interface{}) {
	c.logger.Errorln(args...)
}

func (c *CLogger) Fatalln(args ...interface{}) {
	c.logger.Fatalln(args...)
}

func (c *CLogger) Panicln(args ...interface{}) {
	c.logger.Panicln(args...)
}

func (c *CLogger) TraceFn(fn FormatFunction) {
	c.logger.TraceFn(fn)
}

func (c *CLogger) DebugFn(fn FormatFunction) {
	c.logger.DebugFn(fn)
}

func (c *CLogger) InfoFn(fn FormatFunction) {
	c.logger.InfoFn(fn)
}

func (c *CLogger) PrintFn(fn FormatFunction) {
	c.logger.PrintFn(fn)
}

func (c *CLogger) WarnFn(fn FormatFunction) {
	c.logger.WarnFn(fn)
}

func (c *CLogger) WarningFn(fn FormatFunction) {
	c.logger.WarningFn(fn)
}

func (c *CLogger) ErrorFn(fn FormatFunction) {
	c.logger.ErrorFn(fn)
}

func (c *CLogger) FatalFn(fn FormatFunction) {
	c.logger.FatalFn(fn)
}

func (c *CLogger) PanicFn(fn FormatFunction) {
	c.logger.PanicFn(fn)
}

func (c *CLogger) Trace(args ...interface{}) {
	c.logger.Trace(args...)
}

func (c *CLogger) Debug(args ...interface{}) {
	c.logger.Debug(args...)
}

func (c *CLogger) Info(args ...interface{}) {
	c.logger.Info(args...)
}

func (c *CLogger) Print(args ...interface{}) {
	c.logger.Print(args...)
}

func (c *CLogger) Warn(args ...interface{}) {
	c.logger.Warn(args...)
}

func (c *CLogger) Warning(args ...interface{}) {
	c.logger.Warning(args...)
}

func (c *CLogger) Error(args ...interface{}) {
	c.logger.Error(args...)
}

func (c *CLogger) Fatal(args ...interface{}) {
	c.logger.Fatal(args...)
}

func (c *CLogger) Panic(args ...interface{}) {
	c.logger.Panic(args...)
}

func (c *CLogger) Tracef(format string, args ...interface{}) {
	c.logger.Tracef(format, args...)
}

func (c *CLogger) Debugf(format string, args ...interface{}) {
	c.logger.Debugf(format, args...)
}

func (c *CLogger) Infof(format string, args ...interface{}) {
	c.logger.Infof(format, args...)
}

func (c *CLogger) Printf(format string, args ...interface{}) {
	c.logger.Printf(format, args...)
}

func (c *CLogger) Warnf(format string, args ...interface{}) {
	c.logger.Warnf(format, args...)
}

func (c *CLogger) Warningf(format string, args ...interface{}) {
	c.logger.Warningf(format, args...)
}

func (c *CLogger) Errorf(format string, args ...interface{}) {
	c.logger.Errorf(format, args...)
}

func (c *CLogger) Fatalf(format string, args ...interface{}) {
	c.logger.Fatalf(format, args...)
}

func (c *CLogger) Panicf(format string, args ...interface{}) {
	c.logger.Panicf(format, args...)
}

func (c *CLogger) SetNoLock() {
	c.logger.SetNoLock()
}

func (c *CLogger) Release(w io.Writer) error {
	c.logger.SetOutput(w)
	return c.logger.Writer().Close()
}

// New creates a new Logger object.
func New(opts ...Option) ILogger {
	l := &CLogger{
		logger: logrus.New(),
	}
	//AddHook(&Hook{})
	l.logger.SetFormatter(&Formatter{
		Mu:              new(sync.Mutex),
		TimestampFormat: "2006-01-02 15:04:05.000",
		//Colors:          false,
		//TrimMessages:    false,
		//NoCaller:        false,
		bufPool: &defaultPool{
			pool: &sync.Pool{
				New: func() interface{} {
					return new(bytes.Buffer)
				},
			},
		},
	})
	for _, opt := range opts {
		opt(l.logger)
	}

	return l
}

// NewDefaultLogger 创建一个通用日志对象
// filePath 日志输出目录
// fileName 日志文件名(最终文件名会是 filePath/fileName_20060102150405.log)(fileName为空且开启标准输出的情况下默认输出到stdout,否则无任何输出)
// maxAge 最大存放时间(过期会自动删除)
// rotationTime 自动切分间隔(到期日志自动切换新文件)
// level 日志级别(小于设置级别的信息都会被记录打印,设置级别如果超出限制,默认日志等级为error,取值为panic fatal error warn info debug trace)
// withCaller 是否需要调用者信息
// fullCaller 如果需要打印调用者信息,那么这个参数可以设置调用者信息的详细程度
// withColors 是否需要信息的颜色(基本上只能用于linux的前台打印)
// openStdout 是否开启标准输出(如果fileName为空,且openStdout未开启,那么将不会有任何日志信息被记录)
// TODO 支持远程日志钩子函数,可以将日志发送到远程的日志记录器上(这个函数需要go出去执行,不能阻塞)
func NewDefaultLogger(filePath, fileName string, maxAge, rotationTime time.Duration, level string, withCaller, fullCaller, withColors, openStdout bool) (ILogger, error) {
	var writers []io.Writer
	if len(fileName) > 0 {
		if len(filePath) == 0 {
			filePath = "./" // 默认当前目录
		}
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
			path.Join(filePath, fileName),
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

	// 由于是追加模式,所以默认为无锁(gpt认为这里在多线程环境中可能会产生一些问题,在使用中确实遇到过)
	//logger.SetNoLock()

	return logger, nil
}

func Release(logger ILogger) error {
	if logger == nil || logger.Writer() == nil {
		return nil
	}

	return logger.Release(os.Stdout)
}
