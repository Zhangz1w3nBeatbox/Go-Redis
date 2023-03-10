package tcp

import (
	"Go-Redis/interface/tcp"
	"Go-Redis/lib/logger"
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

// Config tcp配置
type Config struct {
	//监听地址 地址:端口
	Address string
}

func ListenAndServeWithSignal(config *Config, handler tcp.Handler) error {

	listen, err := net.Listen("tcp", config.Address)
	closeChannel := make(chan struct{})

	//信号管道
	signChannel := make(chan os.Signal)

	//监控系统信号-比如关闭tcp服务器 写到 signChannel
	signal.Notify(signChannel, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)

	//取出channel
	go func() {
		sg := <-signChannel
		fmt.Println(sg)

		//如果signChannel的型号是关闭类型的 那就给关闭channel发送型号
		switch sg {
		case syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			closeChannel <- struct{}{}
		}
	}()

	if err != nil {
		return err
	}

	logger.Info("start listen")

	//监听新连接
	ListenAndServe(listen, handler, closeChannel)

	return nil
}

// ListenAndServe :处理连接
func ListenAndServe(
	listener net.Listener,
	handler tcp.Handler,
	closeChanel chan struct{}) {

	go func() {
		<-closeChanel
		logger.Info("正在关闭...")
		_ = listener.Close()
		_ = handler.Close()
	}()

	//结束服务 关闭 listener 和handler
	defer func() {
		_ = listener.Close()
		_ = handler.Close()
	}()

	//上下文
	ctx := context.Background()

	var waitDone sync.WaitGroup

	for true {

		con, err := listener.Accept()

		if err != nil {
			break
		}

		logger.Info("接收新连接")

		//新服务 新客户端 进入队列
		waitDone.Add(1)

		//一个协程处理一个新连接
		go func() {
			//处理完一个请求就减少一个任务
			defer func() {
				waitDone.Done()
			}()
			handler.Handler(ctx, con)
		}()

	}

	//等待任务处理完才能结束
	waitDone.Wait()
}
