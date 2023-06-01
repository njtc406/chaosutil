/*
 * Copyright (c) 2023. YR. All rights reserved
 */

package rpc

type EtcdConf struct {
	BasePath    string   // 服务的根路径
	ServiceAddr string   // ip:prot
	Addr        []string // etcd的集群地址(每个实例的ip:port)
}
