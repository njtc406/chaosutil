/*
 * Copyright (c) 2023. YR. All rights reserved
 */

// Package client
// 模块名: 模块名
// 功能描述: 描述
// 作者:  yr  2023/5/24 0024 20:22
// 最后更新:  yr  2023/5/24 0024 20:22
package client

import (
	"github.com/njtc406/chaosutil/rpc"
	"github.com/rpcxio/libkv/store"
	etcdClientV3 "github.com/rpcxio/rpcx-etcd/client"
	"github.com/smallnest/rpcx/client"
)

// TODO: 每个service需要一个自己的channel用于rpc通信,所以util里面的service就需要加上rpc相关的功能,将功能完全从graceful中解耦出去

// NewRpcXClientUseEtcd 创建一个以etcd作为服务发现方式的rpc客户端
func NewRpcXClientUseEtcd(etcdConf *rpc.EtcdConf, options *store.Config) (client.XClient, error) {
	// 新建etcd服务发现插件
	d, err := etcdClientV3.NewEtcdV3Discovery(etcdConf.BasePath, etcdConf.ServicePath, etcdConf.Addr, false, options)
	if err != nil {
		return nil, err
	}

	// 新建rpc客户端
	cli := client.NewXClient(etcdConf.ServicePath, client.Failtry, client.RoundRobin, d, client.DefaultOption)

	return cli, nil
}
