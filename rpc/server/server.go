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
	"github.com/njtc406/chaosutil/log"
	"github.com/njtc406/chaosutil/rpc"
	"github.com/rpcxio/rpcx-etcd/serverplugin"
	"github.com/smallnest/rpcx/server"
)

// RPCServer rpc服务器
type RPCServer struct {
	svr         *server.Server         // rpc服务器对象
	Addr        string                 // rpc服务器监听地址
	etcdConf    *rpc.EtcdConf          // etcd配置
	serviceList map[string]interface{} // 注册的服务
	initialized bool                   // 初始化状态
	running     bool                   // 运行状态
	closeCh     chan struct{}          // 关闭通道
	logger      log.ILogger            // 日志记录器
}

// NewRPCService 新建一个rpc服务器对象
func NewRPCService(logger log.ILogger) *RPCServer {
	return &RPCServer{
		svr:         server.NewServer(),
		serviceList: make(map[string]interface{}),
		logger:      logger,
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
	}
	if err := p.Start(); err != nil {
		return err
	}
	r.svr.Plugins.Add(r)

	// 注册服务
	for svcName, svc := range serviceList {
		if err := r.svr.RegisterName(svcName, svc, ""); err != nil {
			return err
		}
	}

	r.initialized = true

	return nil
}

// Start rpc服务器启动(保持单线程中调用)
func (r *RPCServer) Start() {
	if r == nil || r.svr == nil {
		r.logger.Panic(rpc.ErrRPCServerNotInit)
	}

	if r.running {
		r.logger.Warnf("the rpc server is running")
		return
	}

	r.running = true

	go r.start()
}

func (r *RPCServer) start() {
	go r.svr.Serve("tcp", r.Addr)

	for r.running {
		select {
		case _, ok := <-r.closeCh:
			if !ok {
				r.running = false
				r.svr.Close()
			}
		}
	}
}

// Stop rpc服务器关闭
func (r *RPCServer) Stop() error {

	return nil
}

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
		BasePath:       etcdConf.BasePath,
	}
	if err := r.Start(); err != nil {
		return nil, err
	}
	s.Plugins.Add(r)

	for svcName, svc := range serviceList {
		s.RegisterName(svcName, svc, "")
	}

	return s, nil
}

func CheckServiceHealth() {

}
