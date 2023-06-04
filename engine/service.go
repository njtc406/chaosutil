/*
 * Copyright (c) 2023. YR. All rights reserved
 */

// Package engine
// 模块名: 模块名
// 功能描述: 描述
// 作者:  yr  2023/4/17 0017 23:22
// 最后更新:  yr  2023/4/17 0017 23:22navicat
package engine

import (
	"context"
)

// IService service的基础功能定义
type IService interface {
	IEngine
	Call(ctx context.Context, method string, args interface{}, reply interface{}) error // 这个只是个示例,具体参数和返回再定
	Send(ctx context.Context, method string, args interface{}, reply interface{}) error // 这个只是个示例,具体参数和返回再定
}

type Service struct {
	msgQue chan interface{} // 通信管道
}

func (s *Service) Init() error {
	return nil
}

func (s *Service) Start() error {
	return nil
}

func (s *Service) Stop() error {
	return nil
}

func (s *Service) Send() error {
	return nil
}

func (s *Service) Call() error {
	return nil
}
