package tcp

import (
	"context"
	"e-highway-collector/flux"
	"net"
)

// Handler TCP服务处理器接口定义
type Handler interface {
	Handle(ctx context.Context, conn net.Conn, lineChan chan flux.Line)
	Close() error
}
