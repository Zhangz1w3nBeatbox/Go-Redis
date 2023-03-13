package database

import (
	"Go-Redis/interface/resp"
	"Go-Redis/resp/reply"
)

// ping命令

func Ping(db *DB, args [][]byte) resp.Reply {
	return reply.MakePongReply()
}

func init() {
	RegisterCommand("ping", Ping, 1)
}
