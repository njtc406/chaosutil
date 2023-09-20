/*
 * Copyright (c) 2023. YR. All rights reserved
 */

// Package translate
// 模块名: 英文翻译
// 功能描述: 对应字符串转换为英文
// 作者:  yr  2023/4/26 0026 23:00
// 最后更新:  yr  2023/4/26 0026 23:00
package translate

func init() {
	Register(EN_US, enUsMap)
}

var enUsMap = map[string]string{
	"Press enter key to exit...": "Press enter key to exit...", // 回车键退出
	"Version":                    "Version",                    // 版本
	"Supported by":               "Supported by",               // 由xxx提供支持
}
