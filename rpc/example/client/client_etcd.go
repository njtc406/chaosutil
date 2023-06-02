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
	myClient "github.com/njtc406/chaosutil/rpc/client"
	"github.com/njtc406/chaosutil/rpc/example/handler"
	"github.com/smallnest/rpcx/client"
	"sync"
	"time"
)

func startClient() {
	conf := &rpc.EtcdConf{
		BasePath: "/test",
		Addr:     []string{"192.168.0.103:2379"},
	}
	rpcClient, err := myClient.NewRpcXClientUseEtcd("aaa", conf, nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	defer rpcClient.Close()

	//args := []*handler.Args{
	//	{Name: "张三"},
	//	{Name: "李四"},
	//	{Name: "王五"},
	//	{Name: "孙六"},
	//}
	//
	//replies := []*handler.Reply{
	//	{},
	//	{},
	//	{},
	//	{},
	//}
	wg := new(sync.WaitGroup)
	//for i := 0; i < len(args); i++ {
	//	wg.Add(1)
	//	go startCall(rpcClient, args[i], replies[i], wg)
	//}

	startTime := time.Now()
	for i := 0; i < 10000; i++ {
		//wg.Add(1)
		//go startCall(rpcClient, &handler.Args{Name: "张三"}, &handler.Reply{}, wg)
		//go startGo(rpcClient, &handler.Args{Name: "张三"}, &handler.Reply{}, wg)
		startCall(rpcClient, &handler.Args{Name: "张三"}, &handler.Reply{}, wg)
	}

	//wg.Wait()
	fmt.Println("cost time:%d", time.Now().Sub(startTime))
}

func startCall(cli client.XClient, args interface{}, reply *handler.Reply, wg *sync.WaitGroup) {
	//defer wg.Done()
	if err := cli.Call(context.Background(), "SayHi", args, reply); err != nil {
		fmt.Println(err)
	} else {
		//fmt.Println(reply.Say)
	}
}

func startGo(cli client.XClient, args interface{}, reply *handler.Reply, wg *sync.WaitGroup) {
	defer wg.Done()
	done := make(chan *client.Call, 1)
	call, err := cli.Go(context.Background(), "SayHi", args, reply, done)
	if err != nil {
		fmt.Println(err)
	} else {
		<-call.Done
		//fmt.Println(reply.Say)
	}
}
