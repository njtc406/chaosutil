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
	err := NewErrCode(1, "验证操作", nil)
	fmt.Println(err)
	err1 := NewErrCode(2, "递归信息", err)
	fmt.Println(err1)
	err.Release()
	err1.Release()
}
