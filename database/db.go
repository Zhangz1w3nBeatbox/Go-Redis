package database

import (
	"Go-Redis/datastruct/dict"
	"Go-Redis/interface/resp"
	"Go-Redis/resp/reply"
	"strings"
)

type DB struct {
	index int
	data  dict.Dict
}

type ExecFunc func(db *DB, args [][]byte) resp.Reply

type CmdLine = [][]byte

func MakeDB() *DB {
	db := &DB{
		data: dict.MakeSyncDict(),
	}
	return db
}

func (db *DB) Exec(c resp.Connection, cmdLine CmdLine) resp.Reply {
	//用户发送的指令 ping or set or setnx
	cmdName := strings.ToLower(string(cmdLine[0]))
	cmd, ok := cmdTable[cmdName]

	if !ok {
		return reply.MakeStandardErrReply("ERR unknown command" + cmdName)
	}

	// 参数个数校验
	if !validateArity(cmd.arity, cmdLine) {
		return reply.MakeArgsNumErrReply(cmdName)
	}

	exectorFunc := cmd.exector

	return exectorFunc(db, cmdLine[1:])
}

func validateArity(arity int, cmdArgs [][]byte) bool {
	return true
}
