/*
 * Copyright (c) 2023. YR. All rights reserved
 */

// Package title
// 模块名: 模块名
// 功能描述: 描述
// 作者:  yr  2023/4/26 0026 22:51
// 最后更新:  yr  2023/4/26 0026 22:51
package title

import (
	"fmt"
	"github.com/njtc406/chaosutil/translate"
)

var version = "1.0"

var str = `
[1;31m ________[0m[1;33m  ___  ___[0m[1;32m  ________[0m[1;34m  ________[0m[1;35m  ________[0m
[1;31m|\   ____\[0m[1;33m|\  \|\  \[0m[1;32m|\   __  \[0m[1;34m|\   __  \[0m[1;35m|\   ____\[0m
[1;31m\ \  \___|[0m[1;33m\ \  \\\  \[0m[1;32m \  \|\  \[0m[1;34m \  \|\  \[0m[1;35m \  \___|_[0m
[1;31m \ \  \    [0m[1;33m\ \   __  \[0m[1;32m \   __  \[0m[1;34m \  \\\  \[0m[1;35m \_____  \[0m
[1;31m  \ \  \____[0m[1;33m\ \  \ \  \[0m[1;32m \  \ \  \[0m[1;34m \  \\\  \[0m[1;35m|____|\  \[0m
[1;31m   \ \_______[0m[1;33m\ \__\ \__\[0m[1;32m \__\ \__\[0m[1;34m \_______\[0m[1;35m____\_\  \[0m
[1;31m    \|_______|[0m[1;33m\|__|\|__|[0m[1;32m\|__|\|__|[0m[1;34m\|_______[0m[1;35m|\_________\[0m
[1;31m              [0m[1;33m           [0m[1;32m          [0m[1;34m        [0m[1;35m\|_________|[0m
`

var title1Base = `
Supported by:
                 ██████╗██╗  ██╗ █████╗  ██████╗ ███████╗
                ██╔════╝██║  ██║██╔══██╗██╔═══██╗██╔════╝
                ██║     ███████║███████║██║   ██║███████╗
                ██║     ██╔══██║██╔══██║██║   ██║╚════██║
                ╚██████╗██║  ██║██║  ██║╚██████╔╝███████║
                 ╚═════╝╚═╝  ╚═╝╚═╝  ╚═╝ ╚═════╝ ╚══════╝
`

var title2 = `
%s:
                ▄▄·  ▄ .▄ ▄▄▄·       .▄▄ · 
                ▐█ ▌▪██▪▐█▐█ ▀█ ▪     ▐█ ▀. 
                ██ ▄▄██▀▐█▄█▀▀█  ▄█▀▄ ▄▀▀▀█▄
                ▐███▌██▌▐▀▐█ ▪▐▌▐█▌.▐▌▐█▄▪▐█
                ·▀▀▀ ▀▀▀ · ▀  ▀  ▀█▄▀▪ ▀▀▀▀ 
						%s: %s
`

var title3 = `
Entropy increase // 熵增
`

func EchoTitle() {
	fmt.Print(fmt.Sprintf(title1, translate.Translate("Supported by"), translate.Translate("Version"), version))
}

func EchoTitleByType(tp int) {
	switch tp {
	case 0:
		fmt.Print(fmt.Sprintf(title1, translate.Translate("Supported by"), translate.Translate("Version"), version))
	case 1:
		fmt.Print(fmt.Sprintf(title2, translate.Translate("Supported by"), translate.Translate("Version"), version))
	}
}
