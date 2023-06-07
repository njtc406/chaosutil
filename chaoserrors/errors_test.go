/*
 * Copyright (c) 2023. YR. All rights reserved
 */

// Package chaoserrors
// 模块名: 模块名
// 功能描述: 描述
// 作者:  yr  2023/6/7 0007 1:36
// 最后更新:  yr  2023/6/7 0007 1:36
package chaoserrors

import (
	"fmt"
	"testing"
)

func TestNewErrCode(t *testing.T) {
	err := NewErrCode(1, nil)
	fmt.Println(err.String())
	fmt.Println(err.StringWithCaller())
}

func TestConvertCodeToError(t *testing.T) {
	err1 := NewErrCode(1, nil)
	err2 := NewErrCode(2, err1)

	fmt.Println(ConvertCodeToError(err2, false))

	err3 := NewErrCode(3, nil)
	err4 := NewErrCode(4, err3)

	fmt.Println(ConvertCodeToError(err4, true))
}
