/*
 * Copyright (c) 2023. YR. All rights reserved
 */

package rpc

type EtcdConf struct {
	BasePath    string
	ServiceAddr string // ip:prot
	Addr        []string
}
