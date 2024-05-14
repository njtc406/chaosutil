/*
 * Copyright (c) 2024. YR. All rights reserved
 */

// Package log
// 模块名: 模块名
// 功能描述: 描述
// 作者:  yr  2024/3/1 0001 11:12
// 最后更新:  yr  2024/3/1 0001 11:12
package log

import (
	"io"
)

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

	TraceFn(fn FormatFunction)
	DebugFn(fn FormatFunction)
	InfoFn(fn FormatFunction)
	PrintFn(fn FormatFunction)
	WarnFn(fn FormatFunction)
	WarningFn(fn FormatFunction)
	ErrorFn(fn FormatFunction)
	FatalFn(fn FormatFunction)
	PanicFn(fn FormatFunction)

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
	Writer() *io.PipeWriter
	WriterLevel(level Level) *io.PipeWriter
	GetOutput() io.Writer
	SetNoLock()
	Release(writer io.Writer) error
}
