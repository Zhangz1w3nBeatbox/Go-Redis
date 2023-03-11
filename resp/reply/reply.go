package reply

import (
	"Go-Redis/interface/resp"
	"bytes"
	"strconv"
)

/*
	定义 服务器 回复 客户端的信息
*/

var (
	nullBulkReplyBytes = []byte("$-1")
	CRLF               = "\r\n"
)

//字符串
type BulkReply struct {
	Arg []byte
}

// test -> $4\r\ntest\r\n
func (b *BulkReply) ToBytes() []byte {
	if len(b.Arg) == 0 {
		return NullBulkBytes
	}
	return []byte("$" + strconv.Itoa(len(b.Arg)) + CRLF + string(b.Arg) + CRLF)
}

func MakeBulkReply(arg []byte) *BulkReply {
	return &BulkReply{
		Arg: arg,
	}
}

//多字符串
type MultiBulkReply struct {
	Args [][]byte
}

//
func (b *MultiBulkReply) ToBytes() []byte {

	argLen := len(b.Args)

	var buf bytes.Buffer

	buf.WriteString("*" + strconv.Itoa(argLen) + CRLF) // *2

	for _, arg := range b.Args {
		if arg == nil {
			buf.WriteString(string(nullBulkReplyBytes) + CRLF)
		} else {
			buf.WriteString("$" + strconv.Itoa(len(arg)) + CRLF + string(arg) + CRLF)
		}
	}

	return buf.Bytes()
}

func MakeMultiBulkReply(args [][]byte) *MultiBulkReply {
	return &MultiBulkReply{
		Args: args,
	}
}

//状态消息回复
// +ok换行

type StatusReply struct {
	Status string
}

func (b *StatusReply) ToBytes() []byte {

	return []byte("+" + b.Status + CRLF)
}

func MakeStatusReply(Status string) *StatusReply {
	return &StatusReply{
		Status: Status,
	}
}

//数字消息回复
// :100 换行

type IntReply struct {
	Code int64
}

func (b *IntReply) ToBytes() []byte {

	return []byte(":" + strconv.FormatInt(b.Code, 10) + CRLF)
}

func MakeIntReply(code int64) *IntReply {
	return &IntReply{
		Code: code,
	}
}

//通用自定义回复
//

type StandardErrReply struct {
	Status string
}

func (b *StandardErrReply) ToBytes() []byte {

	return []byte("-" + b.Status + CRLF)
}

func (b *StandardErrReply) Error() string {

	return b.Status
}

func MakeStandardErrReply(Status string) *StandardErrReply {
	return &StandardErrReply{
		Status: Status,
	}
}

//判断是否为错误回复

func IsErrReply(reply resp.Reply) bool {
	return reply.ToBytes()[0] == '-'
}

type ErrorReply interface {
	Error() string
	ToBytes() []byte
}
