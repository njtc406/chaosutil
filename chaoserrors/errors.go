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
	"sync"
)

// 错误码缓存池,减少GC消耗
var codePool = &sync.Pool{
	New: func() interface{} {
		return &ErrCode{}
	},
}

type caller struct {
	file string
	line int
}

func (c caller) string() string {
	return fmt.Sprintf("%s:%d", c.file, c.line)
}

func (c caller) reset() {
	c.line = 0
	c.file = ""
}

// ErrCode 错误码对象
type ErrCode struct {
	Code    int      // 错误码(这个主要用来给一些地方做判断使用,避免直接判断字符串)
	Msg     string   // 错误信息
	caller           // 调用者信息
	preCode *ErrCode // 前一个错误
}

func (e *ErrCode) reset() {
	e.Code = 0
	e.Msg = ""
	e.preCode = nil
	e.caller.reset()
}

// String 当前错误对象的错误信息
func (e *ErrCode) String() string {
	return fmt.Sprintf("%s %s", e.caller.string(), e.Msg)
}

// Error 返回错误信息(目前暂时修改为手动释放,后面根据需求看需不需要使用之后自动释放)
func (e *ErrCode) Error() string {
	errList := e.getAllErr()

	errStr := make([]string, 0, len(errList))
	for i := len(errList) - 1; i >= 0; i-- {
		errStr = append(errStr, errList[i].String())
		//errList[i].Release()
	}

	return strings.Join(errStr, "\n")
}

// Release 释放ErrCode对象
func (e *ErrCode) Release() {
	e.reset()
	codePool.Put(e)
}

// ReleaseAll 循环释放所有对象
func (e *ErrCode) ReleaseAll() {
	parentCode := e.preCode
	e.preCode = nil
	e.Release()
	var tmpParentCode *ErrCode
	for parentCode != nil {
		tmpParentCode = parentCode
		parentCode = parentCode.preCode
		tmpParentCode.preCode = nil
		tmpParentCode.Release()
	}
}

func (e *ErrCode) getAllErr() []*ErrCode {
	errList := []*ErrCode{e}
	parentCode := e.preCode
	e.preCode = nil
	var tmpParentCode *ErrCode
	for parentCode != nil {
		errList = append(errList, parentCode)
		tmpParentCode = parentCode
		parentCode = parentCode.preCode
		tmpParentCode.preCode = nil
	}

	return errList
}

// NewErrCode 新建错误码(禁止同一个错误码作为多个的父错误!!意思就是是错误只能被单线的传递,最终被一个地方捕获后输出)
// 例如: 错误A返回后,被B捕获,B又添加了一个新的错误,A作为父错误,那么此时A就不能再被其他错误捕获了,只能由B继续上传,直至最终被捕获输出
// 如果要多线传递,请自己理清逻辑,保证最后所有的对象都能正确的使用和释放
func NewErrCode(code int, msg string, parentCode *ErrCode) *ErrCode {
	errCode := codePool.Get().(*ErrCode)
	errCode.Code = code
	errCode.Msg = msg
	errCode.preCode = parentCode

	_, file, line, ok := runtime.Caller(1)
	if ok {
		errCode.line = line
		errCode.file = file
	} else {
		fmt.Printf("code:%d can not get caller info", code)
	}
	return errCode
}
