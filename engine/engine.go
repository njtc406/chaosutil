/*
 * Copyright (c) 2023. YR. All rights reserved
 */

// Package engine
// 模块名: 模块名
// 功能描述: 描述
// 作者:  yr  2023/4/17 0017 23:21
// 最后更新:  yr  2023/4/17 0017 23:21
package engine

// IEngine 这个engine其实应该是service,等下再改
type IEngine interface {
	Init() error
	Start() error
	Stop() error
}
