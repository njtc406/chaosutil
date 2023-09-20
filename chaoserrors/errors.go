/*
 * Copyright (c) 2023. YR. All rights reserved
 */

// Package chaoserrors
// 模块名: 错误码
// 功能描述: 用于错误传递,上层错误捕获可以看到这个错误最初是哪里来的
// 作者:  yr  2023/6/7 0007 0:24
// 最后更新:  yr  2023/6/7 0007 0:24
package chaoserrors

import (
	"fmt"
	"runtime"
	"strings"
)

type cError interface {
	error
	EqualErrCode(int) bool
}

// ErrCode 错误码对象
type ErrCode struct {
	Code   int     // 错误码(这个主要用来给一些地方做判断使用,避免直接判断字符串)
	Msg    string  // 错误信息
	caller *caller // 调用者信息
	preMsg string  // 收集的之前的错误
}

// String 当前错误对象的错误信息
func (e *ErrCode) String() string {
	// TODO 后面再做一个code转string的操作,这样可以预定义一些错误
	if e.caller != nil {
		return fmt.Sprintf("%s ---> code: %d, msg: %s", e.caller.string(), e.Code, e.Msg)
	} else {
		return fmt.Sprintf("---> code: %d, msg: %s", e.Code, e.Msg)
	}
}

// Error 返回错误信息
func (e *ErrCode) Error() string {
	return e.getAllErr()
}

func (e *ErrCode) getAllErr() string {
	builder := new(strings.Builder)
	builder.WriteString(e.String())
	if e.preMsg != "" {
		builder.WriteString("\n")
		builder.WriteString(e.preMsg)
	}
	return builder.String()
}

// EqualErrCode 比较错误码
func (e *ErrCode) EqualErrCode(code int) bool {
	return e.Code == code
}

// NewErrCode 新建错误码
func NewErrCode(code int, msg string, preMsg error) cError {
	errCode := &ErrCode{
		Code:   code,
		Msg:    msg,
		caller: nil,
	}

	if preMsg != nil {
		errCode.preMsg = preMsg.Error()
	}

	// 获取上一层调用者
	_, file, line, ok := runtime.Caller(1)
	if ok {
		callerInfo := &caller{
			line: line,
			file: file,
		}
		errCode.caller = callerInfo
	} else {
		fmt.Printf("code:%d can not get caller info\n", code)
	}
	return errCode
}
