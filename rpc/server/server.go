/*
 * Copyright (c) 2023. YR. All rights reserved
 */

// Package server
// 模块名: 模块名
// 功能描述: 描述
// 作者:  yr  2023/5/25 0025 20:00
// 最后更新:  yr  2023/5/25 0025 20:00
package server

import (
	"fmt"
	"github.com/njtc406/chaosutil/rpc"
	"github.com/rpcxio/rpcx-etcd/serverplugin"
	"github.com/smallnest/rpcx/server"
	"path"
)

// serviceList 服务列表
var serviceList = make(map[string]interface{})

// RegisterService 注册服务
func RegisterService(svcName string, svc interface{}, force bool) error {
	_, ok := serviceList[svcName]
	if ok && !force {
		return rpc.ErrServiceHasRegistered
	}

	serviceList[svcName] = svc

	return nil
}

func NewETCDRpcxServer(addr string, etcdConf *rpc.EtcdConf) (*server.Server, error) {
	s := server.NewServer()
	r := &serverplugin.EtcdV3RegisterPlugin{
		ServiceAddress: addr,
		EtcdServers:    etcdConf.Addr,
		BasePath:       path.Join(etcdConf.BasePath, etcdConf.ServicePath),
		//Metrics:        metrics.NewRegistry(),
		//UpdateInterval: time.Minute,
	}
	if err := r.Start(); err != nil {
		return nil, err
	}
	s.Plugins.Add(r)

	for svcName, svc := range serviceList {
		fmt.Println(svcName)
		s.RegisterName(svcName, svc, "")
	}

	return s, nil
}
