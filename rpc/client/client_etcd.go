/*
 * Copyright (c) 2023. YR. All rights reserved
 */

// Package client
// 模块名: rpc客户端
// 功能描述: 这是一个使用etcd作为注册中心的rpc客户端
// 作者:  yr  2023/5/24 0024 20:22
// 最后更新:  yr  2023/5/24 0024 20:22
package client

import (
	"github.com/njtc406/chaosutil/rpc"
	"github.com/rpcxio/libkv/store"
	etcdClientV3 "github.com/rpcxio/rpcx-etcd/client"
	"github.com/smallnest/rpcx/client"
	"time"
)

// NewRpcXClientUseEtcd 创建一个以etcd作为服务发现方式的rpc客户端
func NewRpcXClientUseEtcd(servicePath string, etcdConf *rpc.EtcdConf, options *store.Config) (client.XClient, error) {
	// 新建etcd服务发现插件
	d, err := etcdClientV3.NewEtcdV3Discovery(etcdConf.BasePath, servicePath, etcdConf.Addr, false, options)
	if err != nil {
		return nil, err
	}

	defaultOption := client.DefaultOption
	defaultOption.Heartbeat = true
	defaultOption.HeartbeatInterval = time.Minute
	defaultOption.ConnectTimeout = time.Second

	// 新建rpc客户端
	cli := client.NewXClient(servicePath, client.Failtry, client.RoundRobin, d, client.DefaultOption)

	return cli, nil
}
