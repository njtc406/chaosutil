/*
 * Copyright (c) 2023. YR. All rights reserved
 */

// Package server
// 模块名: rpc服务器使用示例
// 功能描述: 以etcd作为注册中心的rpc服务器的使用示例代码
// 作者:  yr  2023/6/2 0002 0:37
// 最后更新:  yr  2023/6/2 0002 0:37
package main

import (
	"fmt"
	"github.com/njtc406/chaosutil/rpc"
	"github.com/njtc406/chaosutil/rpc/example/handler"
	"github.com/njtc406/chaosutil/rpc/server"
)

func serverByetcd() {
	conf := &rpc.EtcdConf{
		BasePath:    "/test",
		ServiceAddr: "192.168.0.106:4001",
		Addr:        []string{"192.168.0.103:2379"},
	}
	rpcServer := server.NewRPCService("0.0.0.0:4001", nil, conf)
	svcName := "aaa"
	err := rpcServer.RegisterService(&svcName, new(handler.AAA), false)
	if err != nil {
		fmt.Println("register", err)
		return
	}
	fmt.Println("server begin init")
	err = rpcServer.Init()
	if err != nil {
		fmt.Println("init", err)
		return
	}
	fmt.Println("server start" +
		"")
	err = rpcServer.Start()
	if err != nil {
		fmt.Println("start", err)
		return
	}

	fmt.Println(fmt.Sprintf("receive %v", <-closeCh))

	fmt.Println("server begin stop")
	err = rpcServer.Stop()
	if err != nil {
		fmt.Println("stop", err)
		return
	}

	fmt.Println("exit")
}
