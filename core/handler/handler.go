package handler

import (
	"bufio"
	"context"
	"e-highway-collector/core/reply"
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

// SensorClient 传感器客户端
type SensorClient struct {
	Conn    net.Conn
	Waiting wait.Wait
}

func (e *SensorClient) Close() error {
	e.Waiting.WaitWithTimeout(10 * time.Second)
	_ = e.Conn.Close()
	return nil
}

// SensorHandler 回响Handler
type SensorHandler struct {
	activeConn sync.Map
	closing    atomic.Boolean
}

func (handler *SensorHandler) Handle(_ context.Context, conn net.Conn, ch chan flux.Line) {
	if handler.closing.Get() {
		_ = conn.Close()
	}
	client := &SensorClient{
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
		_, _ = conn.Write(reply.MakeOkReply().ToBytes())
		client.Waiting.Done()
	}
}

func (handler *SensorHandler) Close() error {
	logger.Info("handler shutting down")
	handler.closing.Set(true)
	handler.activeConn.Range(func(key, value interface{}) bool {
		client := key.(*SensorClient)
		_ = client.Conn.Close()
		return true
	})
	return nil
}

func MakeSensorHandler() *SensorHandler {
	return &SensorHandler{}
}
