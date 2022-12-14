package tcp

import (
	"bufio"
	"context"
	"e-highway-collector/flux"
	"e-highway-collector/lib/logger"
	"e-highway-collector/lib/sync/atomic"
	"e-highway-collector/lib/sync/wait"
	"e-highway-collector/parser"
	"io"
	"net"
	"sync"
	"time"
)

// EchoClient 回响服务测试
type EchoClient struct {
	Conn    net.Conn
	Waiting wait.Wait
}

func (e *EchoClient) Close() error {
	e.Waiting.WaitWithTimeout(10 * time.Second)
	_ = e.Conn.Close()
	return nil
}

// EchoHandler 回响Handler
type EchoHandler struct {
	activeConn sync.Map
	closing    atomic.Boolean
}

func (handler *EchoHandler) Handle(_ context.Context, conn net.Conn, ch chan flux.Line) {
	if handler.closing.Get() {
		_ = conn.Close()
	}
	client := &EchoClient{
		Conn: conn,
	}
	reader := bufio.NewReader(conn)
	handler.activeConn.Store(client, struct{}{})
	for {
		msg, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				logger.Info("Connection close")
				handler.activeConn.Delete(client)
			} else {
				logger.Warn(err)
			}
			return
		}
		// 正在向客户端发送数据，任务量 + 1
		client.Waiting.Add(1)
		b := []byte(msg)
		parser.ParseLine(b, ch)
		_, _ = conn.Write(b)
		client.Waiting.Done()
	}
}

func (handler *EchoHandler) Close() error {
	logger.Info("handler shutting down")
	handler.closing.Set(true)
	handler.activeConn.Range(func(key, value interface{}) bool {
		client := key.(*EchoClient)
		_ = client.Conn.Close()
		return true
	})
	return nil
}

func MakeEchoHandler() *EchoHandler {
	return &EchoHandler{}
}
