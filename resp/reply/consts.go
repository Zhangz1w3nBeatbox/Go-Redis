package reply

//通用回复指令

// pong reply
type PongReply struct {
}

var pongBytes = []byte("+PONG\r\n")

func (pongReply PongReply) ToBytes() []byte {
	return pongBytes
}

func MakePongReply() *PongReply {
	return &PongReply{}
}

// ok reply
type OkReply struct {
}

var OkBytes = []byte("+OK\r\n")

func (OkReply OkReply) ToBytes() []byte {
	return OkBytes
}

func MakeOkReply() *OkReply {
	return &OkReply{}
}

// 空 回复
type NullBulkReply struct {
}

var NullBulkBytes = []byte("$-1\r\n")

func (nullBulkReply NullBulkReply) ToBytes() []byte {
	return NullBulkBytes
}

func MakeNullBulkReply() *NullBulkReply {
	return &NullBulkReply{}
}

// 空 数组回复
type EmptyMultiBulkReply struct {
}

var EmptyMultiBulkBytes = []byte("*0\r\n")

func (EmptyMultiBulkReply EmptyMultiBulkReply) ToBytes() []byte {
	return EmptyMultiBulkBytes
}

func MakeEmptyMultiBulkReply() *EmptyMultiBulkReply {
	return &EmptyMultiBulkReply{}
}

//""回复
type NoReply struct {
}

var NoBytes = []byte("")

func (NoReply NoReply) ToBytes() []byte {
	return NoBytes
}

func MakeNoReply() *NoReply {
	return &NoReply{}
}
