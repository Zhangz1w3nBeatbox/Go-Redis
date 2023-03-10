package tcp

import (
	"context"
	"net"
)

// Handler 抽象方法
//
type Handler interface {
	/*
		处理tcp连接
	*/

	Handler(ctx context.Context, conn net.Conn)

	close() error
}
