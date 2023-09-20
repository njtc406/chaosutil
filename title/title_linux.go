//go:build linux || ignore || !windows
// +build linux ignore !windows

/*
 * Copyright (c) 2023. YR. All rights reserved
 */

// Package title
// 模块名:
// 模块功能简介:
package title

import (
	"fmt"
)

var title1 = `
%s:
               [1;31m ██████╗[0m[1;33m██╗  ██╗[0m[1;32m █████╗ [0m[1;34m ██████╗ [0m[1;35m███████╗[0m        
               [1;31m██╔════╝[0m[1;33m██║  ██║[0m[1;32m██╔══██╗[0m[1;34m██╔═══██╗[0m[1;35m██╔════╝[0m        
               [1;31m██║     [0m[1;33m███████║[0m[1;32m███████║[0m[1;34m██║   ██║[0m[1;35m███████╗[0m        
               [1;31m██║     [0m[1;33m██╔══██║[0m[1;32m██╔══██║[0m[1;34m██║   ██║[0m[1;35m╚════██║[0m        
               [1;31m╚██████╗[0m[1;33m██║  ██║[0m[1;32m██║  ██║[0m[1;34m╚██████╔╝[0m[1;35m███████║[0m        
               [1;31m ╚═════╝[0m[1;33m╚═╝  ╚═╝[0m[1;32m╚═╝  ╚═╝[0m[1;34m ╚═════╝ [0m[1;35m╚══════╝[0m        
                              [1;36m%s: %s[0m
`

var bakUrl = "https://patorjk.com/software/taag/#p=display&f=3D%20Diagonal" // 可以在这里做新的title

func EchoByeBye() {
	fmt.Println("exit")
}
