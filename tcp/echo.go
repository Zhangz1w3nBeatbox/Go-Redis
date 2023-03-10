package tcp

import (
	"Go-Redis/lib/logger"
	"Go-Redis/lib/sync/atomic"
	"Go-Redis/lib/sync/wait"
	"bufio"
	"context"
	"fmt"
	"io"
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

func MakeHandler() *EchoHandler {
	return &EchoHandler{}
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

	reader := bufio.NewReader(conn)

	//循环读取客户端发送过来的字符

	for {
		// 接收消息 如果是 \n 回车就写回去
		msg, err := reader.ReadString('\n')

		fmt.Println("接收到数据:", msg)

		//有异常
		if err != nil {
			if err == io.EOF {
				logger.Info("connection close...")
				handler.activeCon.Delete(client)
			} else {
				logger.Warn(err)
			}

			return
		}

		//等客户端读完 之后才能做其他事情 否则 阻塞
		client.waiting.Add(1)

		b := []byte(msg)

		_, _ = conn.Write(b)

		client.waiting.Done()

	}

}

//关闭handler方法
func (handler *EchoHandler) Close() error {
	logger.Info("handler shout down")

	handler.closing.Set(true)

	handler.activeCon.Range(func(key, value interface{}) bool {
		client := key.(*EchoClient)
		_ = client.Conn.Close()
		return true
	})

	return nil
}
