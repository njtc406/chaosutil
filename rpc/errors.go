/*
 * Copyright (c) 2023. YR. All rights reserved
 */

package rpc

import (
	"errors"
)

var (
	ErrServiceHasRegistered = errors.New("the service has already been registered") // 服务已经注册过了
	ErrRPCServerNotInit     = errors.New("rpc server not instantiated")             // 未实例化
	ErrInvalidServiceName   = errors.New("invalid service name")                    // 无效的服务名称
	ErrInvalidService       = errors.New("invalid service")                         // 无效的服务对象
	ErrRPCServerInitialized = errors.New("the service has already been initialized")
)
