package reply

//未知错误
type UnknownErrReply struct {
}

var unknownErrBytes = []byte("-Err unknown\r\n")

func (u UnknownErrReply) Error() string {
	return "-Err unknown"
}

func (u UnknownErrReply) ToBytes() []byte {
	return unknownErrBytes
}

//参数错误
type ArgsNumErrReply struct {
	Cmd string
}

var ArgsNumErrBytes = []byte("-Err unknown\r\n")

func (r *ArgsNumErrReply) Error() string {
	return "-ERR wrong number of arguments for '" + r.Cmd + "' command\r\n"
}

func (r *ArgsNumErrReply) ToBytes() []byte {
	return []byte("-ERR wrong number of arguments for '" + r.Cmd + "' command\r\n")
}

func MakeArgsNumErrReply(cmd string) *ArgsNumErrReply {
	return &ArgsNumErrReply{
		Cmd: cmd,
	}
}

//语法错误
type SyntaxErrReply struct {
}

var SyntaxErrBytes = []byte("-Err syntax error\r\n")

func (r *SyntaxErrReply) Error() string {
	return "Err syntax error"
}

func (r *SyntaxErrReply) ToBytes() []byte {
	return SyntaxErrBytes
}

func MakeSyntaxErrReply() *SyntaxErrReply {
	return &SyntaxErrReply{}
}

//Wrong_type 数据类型错误
type WrongTypeErrReply struct {
}

var WrongTypeErrBytes = []byte("-WRONGTYPE operation against a key holding the wrong kind of value\r\n")

func (r *WrongTypeErrReply) Error() string {
	return "WRONGTYPE operation against a key holding the wrong kind of value"
}

func (r *WrongTypeErrReply) ToBytes() []byte {
	return WrongTypeErrBytes
}

func MakeWrongTypeErrReply() *WrongTypeErrReply {
	return &WrongTypeErrReply{}
}

// ProtocolErrReply 接口协议错误
type ProtocolErrReply struct {
	Msg string
}

func (r *ProtocolErrReply) Error() string {
	return "ERR Protocol error: '" + r.Msg
}

func (r *ProtocolErrReply) ToBytes() []byte {
	return []byte("-ERR Protocol error: '" + r.Msg + "'\r\n")
}

func MakeProtocolErrReply() *ProtocolErrReply {
	return &ProtocolErrReply{}
}
