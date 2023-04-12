/*
 * Copyright (c) 2023. YR. All rights reserved
 */

package network

import (
	"github.com/njtc406/chaosutil/log"
	"net"
	"sync"
	"time"
)

// TODO 需要实现连接池

// TCPClient 用于tcp连接的客户端结构
type TCPClient struct {
	sync.Mutex
	Addr            string        // 连接地址
	ConnNum         int           // 连接数
	ConnectInterval time.Duration // 重连间隔
	PendingWriteNum int           // 待写数量
	ReadDeadline    time.Duration // 读超时
	WriteDeadline   time.Duration // 写超时
	AutoReconnect   bool          // 是否自动重连
	NewAgent        func(*TCPConn) Agent
	cons            ConnSet
	wg              sync.WaitGroup
	closeFlag       bool
	logger          log.ILogger
	// msg parser
	MsgParser
}

func (client *TCPClient) Start() {
	client.init()

	for i := 0; i < client.ConnNum; i++ {
		client.wg.Add(1)
		go client.connect()
	}
}

func (client *TCPClient) init() {
	client.Lock()
	defer client.Unlock()
	if client.logger == nil {
		panic("tcp client requires a logger")
	}

	if client.ConnNum <= 0 {
		client.ConnNum = 1
		client.logger.Infof("invalid ConnNum, reset to ", client.ConnNum)
	}
	if client.ConnectInterval <= 0 {
		client.ConnectInterval = 3 * time.Second
		client.logger.Infof("invalid ConnectInterval, reset to ", client.ConnectInterval)
	}
	if client.PendingWriteNum <= 0 {
		client.PendingWriteNum = 1000
		client.logger.Infof("invalid WriteBuffSize, reset to ", client.PendingWriteNum)
	}
	if client.ReadDeadline == 0 {
		client.ReadDeadline = 15 * time.Second
		client.logger.Infof("invalid ReadTimeOut, reset to ", int64(client.ReadDeadline.Seconds()), "s")
	}
	if client.WriteDeadline == 0 {
		client.WriteDeadline = 15 * time.Second
		client.logger.Infof("invalid WriteTimeOut, reset to ", int64(client.WriteDeadline.Seconds()), "s")
	}
	if client.NewAgent == nil {
		client.logger.Fatal("NewAgent must not be nil")
	}
	if client.cons != nil {
		client.logger.Fatal("client is running")
	}
	if client.LenMsgLen == 0 {
		client.LenMsgLen = DefaultLenMsgLen
	}
	if client.MinMsgLen == 0 {
		client.MinMsgLen = DefaultMinMsgLen
	}
	if client.MaxMsgLen == 0 {
		client.MaxMsgLen = DefaultMaxMsgLen
	}

	client.cons = make(ConnSet)
	client.closeFlag = false

	// msg parser
	client.MsgParser.init()
}

func (client *TCPClient) GetCloseFlag() bool {
	client.Lock()
	defer client.Unlock()

	return client.closeFlag
}

func (client *TCPClient) dial() net.Conn {
	for {
		conn, err := net.Dial("tcp", client.Addr)
		if client.closeFlag {
			return conn
		} else if err == nil && conn != nil {
			conn.(*net.TCPConn).SetNoDelay(true)
			return conn
		}

		client.logger.Warnf("connect to %s error: %s", client.Addr, err.Error())
		time.Sleep(client.ConnectInterval)
		continue
	}
}

func (client *TCPClient) connect() {
	defer client.wg.Done()

reconnect:
	conn := client.dial()
	if conn == nil {
		return
	}

	client.Lock()
	if client.closeFlag {
		client.Unlock()
		conn.Close()
		return
	}
	client.cons[conn] = struct{}{}
	client.Unlock()

	tcpConn := newTCPConn(conn, client.PendingWriteNum, &client.MsgParser, client.WriteDeadline, client.logger)
	agent := client.NewAgent(tcpConn)
	agent.Run()

	// cleanup
	tcpConn.Close()
	client.Lock()
	delete(client.cons, conn)
	client.Unlock()
	agent.OnClose()

	if client.AutoReconnect {
		time.Sleep(client.ConnectInterval)
		goto reconnect
	}
}

func (client *TCPClient) Close(waitDone bool) {
	client.Lock()
	client.closeFlag = true
	for conn := range client.cons {
		conn.Close()
	}
	client.cons = nil
	client.Unlock()

	if waitDone == true {
		client.wg.Wait()
	}
}
