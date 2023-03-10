package main

import (
	"Go-Redis/lib/sync/atomic"
	"Go-Redis/lib/sync/wait"
	"context"
	"net"
	"sync"
	"time"
)

//客户端
type EchoClient struct {
	Conn    net.Conn
	waiting wait.Wait //自己封装的wait 有超时间的功能
}

//客户端关闭
func (e *EchoClient) Close() error {
	e.waiting.WaitWithTimeout(10 * time.Second)
	_ = e.Conn.Close()
	return nil
}

//回复处理器
type EchoHandler struct {
	activeCon sync.Map
	closing   atomic.Boolean
}

func (handler *EchoHandler) Handler(ctx context.Context, conn net.Conn) {

	if handler.closing.Get() {
		_ = conn.Close()
	}

	//新来的连接 封装成 client
	client := &EchoClient{
		Conn: conn,
	}

	//EchoHandler的map中 存放这个client
	handler.activeCon.Store(client, struct{}{})

}

func (handler *EchoHandler) Close() error {
	return nil
}

func main() {

}
