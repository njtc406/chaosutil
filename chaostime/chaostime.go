/*
 * Copyright (c) 2023. YR. All rights reserved
 */

// Package chaostime
// 模块名: 时间模块
// 功能描述: 用于游戏时间的管理
// 作者:  yr  2023/5/20 0020 23:06
// 最后更新:  yr  2023/5/20 0020 23:06
package chaostime

import "time"

var offset time.Duration = 0

// SetTimeOffset 设置时间偏移
func SetTimeOffset(oset time.Duration) {
	offset = oset
}

// ResetTimeOffset 重置时间偏移
func ResetTimeOffset() {
	offset = 0
}

// Now 获取时间偏移后的当前时间
func Now() time.Time {
	return time.Now().Add(offset)
}
