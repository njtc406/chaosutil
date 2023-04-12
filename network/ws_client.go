package network

import (
	"github.com/gorilla/websocket"
	"github.com/njtc406/chaosutil/log"
	"sync"
	"time"
)

type WSClient struct {
	sync.Mutex
	Addr             string
	ConnNum          int
	ConnectInterval  time.Duration
	PendingWriteNum  int
	MaxMsgLen        uint32
	MessageType      int
	HandshakeTimeout time.Duration
	AutoReconnect    bool
	NewAgent         func(*WSConn) Agent
	dialer           websocket.Dialer
	cons             WebsocketConnSet
	wg               sync.WaitGroup
	closeFlag        bool
	logger           log.ILogger
}

func (client *WSClient) Start() {
	client.init()

	for i := 0; i < client.ConnNum; i++ {
		client.wg.Add(1)
		go client.connect()
	}
}

func (client *WSClient) init() {
	client.Lock()
	defer client.Unlock()

	if client.logger == nil {
		panic("websocket client required a logger")
	}

	if client.ConnNum <= 0 {
		client.ConnNum = 1
		client.logger.Infof("invalid ConnNum, reset to %d", client.ConnNum)
	}
	if client.ConnectInterval <= 0 {
		client.ConnectInterval = 3 * time.Second
		client.logger.Infof("invalid ConnectInterval, reset to %d", client.ConnectInterval)
	}
	if client.PendingWriteNum <= 0 {
		client.PendingWriteNum = 100
		client.logger.Infof("invalid WriteBuffSize, reset to %d", client.PendingWriteNum)
	}
	if client.MaxMsgLen <= 0 {
		client.MaxMsgLen = 4096
		client.logger.Infof("invalid MaxMsgLen, reset to %d", client.MaxMsgLen)
	}
	if client.HandshakeTimeout <= 0 {
		client.HandshakeTimeout = 10 * time.Second
		client.logger.Infof("invalid HandshakeTimeout, reset to %d", client.HandshakeTimeout)
	}
	if client.NewAgent == nil {
		client.logger.Fatal("NewAgent must not be nil")
	}
	if client.cons != nil {
		client.logger.Fatal("client is running")
	}

	if client.MessageType == 0 {
		client.MessageType = websocket.TextMessage
	}

	client.cons = make(WebsocketConnSet)
	client.closeFlag = false
	client.dialer = websocket.Dialer{
		HandshakeTimeout: client.HandshakeTimeout,
	}
}

func (client *WSClient) dial() *websocket.Conn {
	for {
		conn, _, err := client.dialer.Dial(client.Addr, nil)
		if err == nil || client.closeFlag {
			return conn
		}

		client.logger.Infof("connect to ", client.Addr, " error: ", err.Error())
		time.Sleep(client.ConnectInterval)
		continue
	}
}

func (client *WSClient) connect() {
	defer client.wg.Done()

reconnect:
	conn := client.dial()
	if conn == nil {
		return
	}
	conn.SetReadLimit(int64(client.MaxMsgLen))

	client.Lock()
	if client.closeFlag {
		client.Unlock()
		conn.Close()
		return
	}
	client.cons[conn] = struct{}{}
	client.Unlock()

	wsConn := newWSConn(conn, client.PendingWriteNum, client.MaxMsgLen, client.MessageType, client.logger)
	agent := client.NewAgent(wsConn)
	agent.Run()

	// cleanup
	wsConn.Close()
	client.Lock()
	delete(client.cons, conn)
	client.Unlock()
	agent.OnClose()

	if client.AutoReconnect {
		time.Sleep(client.ConnectInterval)
		goto reconnect
	}
}

func (client *WSClient) Close() {
	client.Lock()
	client.closeFlag = true
	for conn := range client.cons {
		conn.Close()
	}
	client.cons = nil
	client.Unlock()

	client.wg.Wait()
}
