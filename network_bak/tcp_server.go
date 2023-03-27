/*
 * Copyright (c) 2023. YR. All rights reserved
 */

package network_bak

import (
	"github.com/njtc406/chaosutil/log"
	"net"
	"sync"
	"time"
)

const (
	DefaultReadTimeOut     = time.Second * 30 // 默认读超时时间,30s
	DefaultWriteTimeOut    = time.Second * 30 // 默认写超时时间,30s
	DefaultMaxConnNum      = 9000             // 默认最大连接数
	DefaultPendingWriteNum = 10000            // 默认最大待写入数量
	DefaultLittleEndian    = false            // 默认大端字节序
	DefaultMinMsgLen       = 2                // 默认最小消息长度
	DefaultMaxMsgLen       = 65535            // 默认最大消息长度
	DefaultLenMsgLen       = 2                // 默认消息头长度,2字节
)

type GetAgentFun func(*TCPConn) Agent

type TCPServer struct {
	Addr          string        // 服务器监听地址
	MaxConnNum    int           // 支持的最大连接数量
	WriteBuffSize int           // 写队列大小
	ReadTimeOut   time.Duration // 读超时
	WriteTimeOut  time.Duration // 写超时

	NewAgent      GetAgentFun
	ln            net.Listener // 监听器
	connPool      ConnSet      // 连接池
	mutexConnPool sync.Mutex   // 连接锁
	wgLn          sync.WaitGroup
	wgConnPool    sync.WaitGroup

	logger log.ILogger // 日志的记录器
	// msgParser 消息的解析器
	MsgParser
}

func (server *TCPServer) Start() {
	server.init()
	go server.run()
}

func (server *TCPServer) init() {
	if server.logger == nil {
		panic("tcp server required a logger")
	}
	ln, err := net.Listen("tcp", server.Addr)
	if err != nil {
		server.logger.Fatalf("Listen tcp error: %s", err.Error())
	}

	if server.MaxConnNum <= 0 {
		server.MaxConnNum = DefaultMaxConnNum
		server.logger.Infof("invalid MaxConnNum, reset to %d", server.MaxConnNum)
	}
	if server.WriteBuffSize <= 0 {
		server.WriteBuffSize = DefaultPendingWriteNum
		server.logger.Infof("invalid WriteBuffSize, reset to %d", server.WriteBuffSize)
	}

	if server.MinMsgLen <= 0 {
		server.MinMsgLen = DefaultMinMsgLen
		server.logger.Infof("invalid MinMsgLen, reset to %d", server.MinMsgLen)
	}

	if server.MaxMsgLen <= 0 {
		server.MaxMsgLen = DefaultMaxMsgLen
		server.logger.Infof("invalid MaxMsgLen, reset to %d", server.MaxMsgLen)
	}

	if server.WriteTimeOut == 0 {
		server.WriteTimeOut = DefaultWriteTimeOut
		server.logger.Infof("invalid WriteTimeOut, reset to %ds", server.WriteTimeOut.Seconds())
	}

	if server.ReadTimeOut == 0 {
		server.ReadTimeOut = DefaultReadTimeOut
		server.logger.Infof("invalid ReadTimeOut, reset to %ds", server.ReadTimeOut.Seconds())
	}

	if server.LenMsgLen == 0 {
		server.LenMsgLen = DefaultLenMsgLen
	}

	if server.NewAgent == nil {
		server.logger.Fatal("NewAgent must not be nil")
	}

	server.ln = ln
	server.connPool = make(ConnSet)
	server.INetMempool = NewMemAreaPool()

	server.MsgParser.init()
}

func (server *TCPServer) SetNetMempool(mempool INetMempool) {
	server.INetMempool = mempool
}

func (server *TCPServer) GetNetMempool() INetMempool {
	return server.INetMempool
}

func (server *TCPServer) run() {
	server.wgLn.Add(1)
	defer server.wgLn.Done()

	var tempDelay time.Duration
	for {
		conn, err := server.ln.Accept()
		if err != nil {
			if ne, ok := err.(net.Error); ok && ne.Temporary() {
				if tempDelay == 0 {
					tempDelay = 5 * time.Millisecond
				} else {
					tempDelay *= 2
				}
				if max := 1 * time.Second; tempDelay > max {
					tempDelay = max
				}
				server.logger.Infof("accept error: %s; retrying in ", err.Error(), tempDelay)
				time.Sleep(tempDelay)
				continue
			}
			return
		}

		conn.(*net.TCPConn).SetNoDelay(true)
		tempDelay = 0

		server.mutexConnPool.Lock()
		if len(server.connPool) >= server.MaxConnNum {
			server.mutexConnPool.Unlock()
			conn.Close()
			server.logger.Warn("too many connections")
			continue
		}

		server.connPool[conn] = struct{}{}
		server.mutexConnPool.Unlock()
		server.wgConnPool.Add(1)

		tcpConn := newTCPConn(conn, server.WriteBuffSize, &server.MsgParser, server.WriteTimeOut, server.logger)
		agent := server.NewAgent(tcpConn)

		go func() {
			agent.Run()
			// cleanup
			tcpConn.Close()
			server.mutexConnPool.Lock()
			delete(server.connPool, conn)
			server.mutexConnPool.Unlock()
			agent.OnClose()

			server.wgConnPool.Done()
		}()
	}
}

func (server *TCPServer) Close() {
	server.ln.Close()
	server.wgLn.Wait()

	server.mutexConnPool.Lock()
	for conn := range server.connPool {
		conn.Close()
	}
	server.connPool = nil
	server.mutexConnPool.Unlock()
	server.wgConnPool.Wait()
}
