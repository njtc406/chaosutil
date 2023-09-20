/*
 * Copyright (c) 2023. YR. All rights reserved
 */

// Package handler
// 模块名: 模块名
// 功能描述: 描述
// 作者:  yr  2023/6/2 0002 0:49
// 最后更新:  yr  2023/6/2 0002 0:49
package handler

import "context"

type AAA struct {
}

type Args struct {
	Name string
}

type Reply struct {
	Say string
}

func (a *AAA) SayHi(ctx context.Context, args *Args, reply *Reply) error {
	reply.Say = "say hi to " + args.Name
	return nil
}
