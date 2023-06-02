/*
 * Copyright (c) 2023. YR. All rights reserved
 */

// Package client
// 模块名: 模块名
// 功能描述: 描述
// 作者:  yr  2023/6/2 0002 0:46
// 最后更新:  yr  2023/6/2 0002 0:46
package main

import (
	"context"
	"fmt"
	"github.com/njtc406/chaosutil/rpc"
	"github.com/njtc406/chaosutil/rpc/client"
	"github.com/njtc406/chaosutil/rpc/example/handler"
)

func startClient() {
	conf := &rpc.EtcdConf{
		BasePath: "/test",
		Addr:     []string{"192.168.0.103:2379"},
	}
	rpcClient, err := client.NewRpcXClientUseEtcd("aaa", conf, nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	defer rpcClient.Close()

	//wg := new(sync.WaitGroup)

	args := []*handler.Args{
		{Name: "张三"},
		{Name: "李四"},
		{Name: "王五"},
		{Name: "孙六"},
	}

	replies := []*handler.Reply{
		{},
		{},
		{},
		{},
	}

	for i := 0; i < len(args); i++ {
		if err = rpcClient.Call(context.Background(), "SayHi", args[i], replies[i]); err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(replies[i].Say)
		}
	}

}
