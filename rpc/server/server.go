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
	"github.com/njtc406/chaosutil/rpc"
	"github.com/rpcxio/rpcx-etcd/serverplugin"
	"github.com/smallnest/rpcx/server"
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

func NewETCDRpcxServer(etcdConf *rpc.EtcdConf) {
	s := server.NewServer()

	for svcName, svc := range serviceList {
		s.RegisterName(svcName, svc, "")
	}

	r := &serverplugin.EtcdV3RegisterPlugin{
		//ServiceAddress: "tcp@127.0.0.1:8080",
		//EtcdServers:    []string{"localhost:2379"},
		//BasePath:       etcdConf.BasePath,
		//Services:       []string{},
		//UpdateInterval: time.Second * 10,
		ServiceAddress: "",
		EtcdServers:    nil,
		BasePath:       "",
		Metrics:        nil,
		Services:       nil,
		UpdateInterval: 0,
		Expired:        0,
		Options:        nil,
	}
	s.Plugins.Add(r)
}
