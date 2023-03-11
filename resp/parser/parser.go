package parser

import (
	"Go-Redis/interface/resp"
	"io"
)

type PayLoad struct {
	Data  resp.Reply
	Error error
}

//解析器状态
type readState struct {
	//正在解析的数据是多行还是单行
	readingMultiLine bool

	//正在读取的指令应该有几个参数/期望的参数个数
	expectedArgsCount int

	//消息类型
	msgType byte

	//已经解析的参数的列表 set key val
	args [][]byte

	//字节组的长度
	bulkLen int64
}

//解析器是否解析完成
func (s *readState) finished() bool {
	return s.expectedArgsCount > 0 && len(s.args) == s.expectedArgsCount
}

//
func ParseStream(reader io.Reader) chan<- *PayLoad {
	ch := make(chan *PayLoad)
	//异步解析-减少tcp的阻塞等待
	go parse0(reader, ch)

	return ch
}

func parse0(reader io.Reader, ch chan<- *PayLoad) {

}
