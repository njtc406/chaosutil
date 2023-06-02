/*
 * Copyright (c) 2023. YR. All rights reserved
 */

// Package server
// 模块名: rpc服务器
// 功能描述: 这是一个以etcd作为注册中心的rpc服务器模块
// 作者:  yr  2023/5/25 0025 20:00
// 最后更新:  yr  2023/5/25 0025 20:00
package server

import (
	"github.com/njtc406/chaosutil/log"
	"github.com/njtc406/chaosutil/rpc"
	"github.com/rcrowley/go-metrics"
	"github.com/rpcxio/rpcx-etcd/serverplugin"
	"github.com/smallnest/rpcx/server"
	"time"
)

// TODO 之后需要整理一下，把一些东西提出来作为配置
// TODO 目前只是最简单的使用版本,后面逐渐修改

// RPCServer rpc服务器
type RPCServer struct {
	svr         *server.Server         // rpc服务器对象
	addr        string                 // rpc服务器监听地址
	etcdConf    *rpc.EtcdConf          // etcd配置
	serviceList map[string]interface{} // 注册的服务
	initialized bool                   // 初始化状态
	running     bool                   // 运行状态
	closeCh     chan struct{}          // 关闭通道
	logger      log.ILogger            // 日志记录器
}

// NewRPCService 新建一个rpc服务器对象
// @params addr ip:port
func NewRPCService(addr string, logger log.ILogger, etcdConf *rpc.EtcdConf) *RPCServer {
	return &RPCServer{
		svr:         server.NewServer(),
		serviceList: make(map[string]interface{}),
		logger:      logger,
		closeCh:     make(chan struct{}),
		etcdConf:    etcdConf,
		addr:        addr,
	}
}

// RegisterService 向rpc服务器对象中注册方法(保持单线程中调用)
func (r *RPCServer) RegisterService(svcName *string, svc interface{}, force bool) error {
	if r == nil || r.serviceList == nil {
		return rpc.ErrRPCServerNotInit
	}

	if svcName == nil || len(*svcName) == 0 {
		return rpc.ErrInvalidServiceName
	}

	if svc == nil {
		// 这个判断需要测试
		return rpc.ErrInvalidService
	}

	_, ok := r.serviceList[*svcName]
	if ok && !force {
		return rpc.ErrServiceHasRegistered
	}

	r.serviceList[*svcName] = svc

	return nil
}

// Init rpc服务器初始化(保持单线程中调用)
func (r *RPCServer) Init() error {
	if r == nil || r.svr == nil {
		return rpc.ErrRPCServerNotInit
	}

	if r.initialized {
		return rpc.ErrRPCServerInitialized
	}

	// 初始化etcd插件
	p := &serverplugin.EtcdV3RegisterPlugin{
		ServiceAddress: "tcp@" + r.etcdConf.ServiceAddr,
		EtcdServers:    r.etcdConf.Addr,
		BasePath:       r.etcdConf.BasePath,
		Metrics:        metrics.NewRegistry(),
		UpdateInterval: time.Minute,
	}
	if err := p.Start(); err != nil {
		return err
	}
	r.svr.Plugins.Add(p)

	// 注册服务
	for svcName, svc := range r.serviceList {
		if err := r.svr.RegisterName(svcName, svc, ""); err != nil {
			return err
		}
	}

	r.initialized = true

	return nil
}

// Start rpc服务器启动(保持单线程中调用)
func (r *RPCServer) Start() error {
	if r == nil || r.svr == nil {
		return rpc.ErrRPCServerNotInit
	}

	if r.running {
		return rpc.ErrRPCServerIsRunning
	}

	r.running = true

	go r.start()

	return nil
}

func (r *RPCServer) start() {

	go r.svr.Serve("tcp", r.addr)

	for r.running {
		select {
		case <-r.closeCh:
			r.running = false
			r.svr.Close()
		}
	}
}

// Stop rpc服务器关闭
func (r *RPCServer) Stop() error {
	if r == nil || r.svr == nil {
		return rpc.ErrRPCServerNotInit
	}

	r.closeCh <- struct{}{}

	return nil
}

func (r *RPCServer) Release() {
	// TODO 释放资源
}
