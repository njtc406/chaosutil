/*
 * Copyright (c) 2023. YR. All rights reserved
 */

// Package engine
// 模块名: 模块名
// 功能描述: 描述
// 作者:  yr  2023/4/17 0017 23:22
// 最后更新:  yr  2023/4/17 0017 23:22
package engine

type Reply interface {
}

type IService interface {
	IEngine
	Call(interface{}, Reply) error // 这个只是个示例,具体参数和返回再定
	Send(interface{})              // 这个只是个示例,具体参数和返回再定
}

type Service struct {
	msgQue chan interface{} // 通信管道
}

func (s *Service) Init() {

}

func (s *Service) Start() {

}

func (s *Service) Stop() {

}

func (s *Service) Send() {

}

func (s *Service) Call() {

}
