/*
 * Copyright (c) 2023. YR. All rights reserved
 */

// Package translate
// 模块名: 模块名
// 功能描述: 描述
// 作者:  yr  2023/4/26 0026 23:01
// 最后更新:  yr  2023/4/26 0026 23:01
package translate

type LanguageType uint8

const (
	ZH_CN        LanguageType = iota // 中文
	EN_US                            // 英文
	LANGUAGE_MAX                     // 最大语言数
)

var languageConf = ZH_CN                     // 当前的语言(默认中文)
var transMap [LANGUAGE_MAX]map[string]string // 这里的下标和上面对应的语言类型对应

func init() {
	transMap = [LANGUAGE_MAX]map[string]string{}
}

func register(languageType LanguageType, content map[string]string) {
	if transMap[languageType] == nil {
		transMap[languageType] = make(map[string]string, len(content))
	}

	for k, v := range content {
		transMap[languageType][k] = v
	}
}

// SetLanguage 设置当前语言类型
func SetLanguage(languageType LanguageType) {
	languageConf = languageType
}

// Translate 翻译字符串为对应语言类型
func Translate(str string) string {
	if transMap[languageConf] == nil {
		return str
	}

	trans, ok := transMap[languageConf][str]
	if !ok {
		return str
	}

	return trans
}
