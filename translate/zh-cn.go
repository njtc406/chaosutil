/*
 * Copyright (c) 2023. YR. All rights reserved
 */

// Package translate
// 模块名: 模块名
// 功能描述: 描述
// 作者:  yr  2023/4/26 0026 23:01
// 最后更新:  yr  2023/4/26 0026 23:01
package translate

func init() {
	Register(ZH_CN, zhCnMap)
}

var zhCnMap = map[string]string{
	"Press enter key to exit...": "按回车键退出...",  // 回车键退出
	"Version":                    "版本号",        // 版本
	"Supported by":               "由chaos提供支持", // 提供支持
}
