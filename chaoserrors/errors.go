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
	"errors"
	"fmt"
	"runtime"
	"strings"
	"sync"
)

// 错误码缓存池,减少GC消耗
var codePool *sync.Pool

// 错误码对应的翻译文字(这个后面看看有没有必要,是直接给错误信息好还是给错误码好,给错误码的话主要上层可以根据不同的错误码来处理,错误信息的话只是看起来不那么舒服)
// 还有一个,这个东西应该可以使用sync.Map,少写多读场景
var errStrMap map[int]string
var mapMu *sync.RWMutex

func init() {
	codePool = &sync.Pool{
		New: func() interface{} {
			return &ErrCode{}
		},
	}

	errStrMap = make(map[int]string)
}

// RegisterCodeString 注册错误码对应的错误字符串
func RegisterCodeString(code int, codeStr string) {
	mapMu.Lock()
	defer mapMu.Unlock()
	errStrMap[code] = codeStr
}

type caller struct {
	file string
	line int
}

func (c caller) String() string {
	return fmt.Sprintf("%s:%d", c.file, c.line)
}

func (c caller) Reset() {
	c.line = 0
	c.file = ""
}

// TODO 这个东西我这么考虑的,这里设置一个code，然后获取一个行号等信息，
// 上层在转换这个code的时候就连行号一起返回，这样即使只在最上层处理错误，也能知道这个错误是出自哪里的返回

// ErrCode 错误码对象
type ErrCode struct {
	Code int
	caller
	parentCode *ErrCode
	// desc gpt建议我在这里加一个错误的描述信息(可能不是错误翻译后的信息,只是用来调试的时候可以知道这个错误码的一些描述,能更快定位错误本身),可以考虑
}

func (e *ErrCode) Reset() {
	e.Code = 0
	e.parentCode = nil
	e.caller.Reset()
}

func (e *ErrCode) String() string {
	return TransformErrCode(e.Code)
}

func (e *ErrCode) StringWithCaller() string {
	return fmt.Sprintf("%s %s", e.caller.String(), TransformErrCode(e.Code))
}

// ConvertCodeToError 转换为标准错误(转换完成后会释放错误对象)
func ConvertCodeToError(errCode *ErrCode, withCaller bool) error {
	errList := []*ErrCode{errCode}
	parentCode := errCode.parentCode
	errCode.parentCode = nil
	var tmpParentCode *ErrCode
	for parentCode != nil {
		errList = append(errList, parentCode)
		tmpParentCode = parentCode
		parentCode = parentCode.parentCode
		tmpParentCode.parentCode = nil
	}

	var errStr []string
	if withCaller {
		for i := len(errList) - 1; i >= 0; i-- {
			errStr = append(errStr, errList[i].StringWithCaller())
			Release(errList[i])
		}
	} else {
		for i := len(errList) - 1; i >= 0; i-- {
			errStr = append(errStr, errList[i].String())
			Release(errList[i])
		}
	}

	return errors.New(strings.Join(errStr, "\n"))
}

// NewErrCode 新建错误码(禁止同一个错误码作为多个的父错误!!意思就是是错误只能被单线的传递,最终被一个地方捕获后输出)
// 例如: 错误A返回后,被B捕获,B又添加了一个新的错误,A作为父错误,那么此时A就不能再被其他错误捕获了,只能由B继续上传,直至最终被捕获输出
// 由于转换时会自动释放掉,所以被多个不同错误捕获会发生一个错误被转换释放后,所有的父错误已经被释放了
func NewErrCode(code int, parentCode *ErrCode) *ErrCode {
	errCode := codePool.Get().(*ErrCode)
	errCode.Code = code
	errCode.parentCode = parentCode
	_, file, line, ok := runtime.Caller(1)
	if ok {
		errCode.line = line
		errCode.file = file
	} else {
		fmt.Printf("code:%d can not get caller info", code)
	}
	return errCode
}

// Release 释放ErrCode对象
func Release(errCode *ErrCode) {
	// TODO 考虑在释放对象前先检查该对象是否已经被释放，避免出现重复释放的情况。
	// 这个需要怎么判断?
	errCode.Reset()
	codePool.Put(errCode)
}

// TransformErrCode 翻译错误码
func TransformErrCode(code int) string {
	mapMu.RLock()
	defer mapMu.RUnlock()
	str, ok := errStrMap[code]
	if ok {
		return str
	} else {
		return fmt.Sprintf("unhandled error code:%d", code)
	}
}

// TODO 1. https://github.com/go-kratos/kratos这个里面有个相同的实现，需要找时间看看别人怎么实现的,参考一下有啥改进
// TODO 2. 需要把之前的模块都替换成新的错误处理
