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
	"context"
	"fmt"
	"github.com/smallnest/rpcx/client"
)

// TODO: 每个service需要一个自己的channel用于rpc通信,所以util里面的service就需要加上rpc相关的功能,将功能完全从graceful中解耦出去

type RpcxClient struct {
	client client.XClient
}

type RpcxService struct {
	xclient client.XClient
	method  string
}

func NewRpcxClient(servicePath string, serviceMethod string, protocol string, address string) (*RpcxClient, error) {
	d, err := client.NewMultipleServersDiscovery([]*client.KVPair{
		{Key: servicePath, Value: address},
	})
	if err != nil {
		return nil, err
	}
	xclient := client.NewXClient(serviceMethod, client.Failover, client.RandomSelect, d, client.DefaultOption)

	return &RpcxClient{client: xclient}, nil
}

func (c *RpcxClient) GetService(method string) *RpcxService {
	return &RpcxService{
		xclient: c.client,
		method:  method,
	}
}

func (s *RpcxService) Call(args interface{}, reply interface{}) error {
	err := s.xclient.Call(context.Background(), s.method, args, reply)
	if err != nil {
		return fmt.Errorf("RPCX call error: %s", err.Error())
	}
	return nil
}

func (s *RpcxService) CallWithRetry(args interface{}, reply interface{}, retryCount int) error {
	var lastErr error
	for i := 0; i < retryCount; i++ {
		err := s.Call(args, reply)
		if err == nil {
			return nil
		}
		lastErr = err
	}
	return fmt.Errorf("RPCX call with retry error: %s", lastErr.Error())
}

func (c *RpcxClient) Close() error {
	return c.client.Close()
}
