/*
 * Copyright (c) 2023. YR. All rights reserved
 */

package network

import (
	"net"
)

type Conn interface {
	ReadMsg() ([]byte, error)
	WriteMsg(args ...[]byte) error
	LocalAddr() net.Addr
	RemoteAddr() net.Addr
	Close()
	Destroy()
	ReleaseReadMsg(byteBuff []byte)
}
