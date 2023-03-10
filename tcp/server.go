package main

import (
	"Go-Redis/interface/tcp"
	"Go-Redis/lib/logger"
	"context"
	"net"
)

// Config tcp配置
type Config struct {
	//监听地址 地址:端口
	Address string
}

func ListenAndServeWithSignal(config *Config, handler tcp.Handler) error {

	listen, err := net.Listen("tcp", config.Address)

	if err != nil {
		return err
	}

	logger.Info("start listen")

	channel := make(chan struct{})

	//监听新连接
	ListenAndServe(listen, handler, channel)

	return nil
}

func ListenAndServe(listener net.Listener,
	handler tcp.Handler,
	chanel chan struct{}) error {

	//上下文
	ctx := context.Background()

	for true {
		con, err := listener.Accept()

		if err != nil {
			break
		}

		logger.Info("接收新连接")

		//一个携程处理一个连接
		go func() {
			handler.Handler(ctx, con)
		}()
	}

	return nil
}
