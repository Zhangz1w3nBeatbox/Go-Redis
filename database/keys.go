package database

import (
	"Go-Redis/interface/resp"
	"Go-Redis/lib/wildcard"
	"Go-Redis/resp/reply"
)

// DEL K1 K2 K3
func execDel(db *DB, args [][]byte) resp.Reply {
	keys := make([]string, len(args))

	//
	for i, v := range keys {
		keys[i] = string(v)

	}

	deleted := db.Removes(keys...)

	return reply.MakeIntReply(int64(deleted))
}

// Exists K1 K2 K3
func execExists(db *DB, args [][]byte) resp.Reply {
	res := int64(0)

	for _, arg := range args {
		key := string(arg)
		_, exists := db.GetEntity(key)
		if exists {
			res++
		}
	}

	return reply.MakeIntReply(res)

}

// FlushDB
func execFlush(db *DB, args [][]byte) resp.Reply {
	db.Flush()
	return reply.MakeOkReply()
}

// Type
func execType(db *DB, args [][]byte) resp.Reply {
	key := string(args[0])
	entity, exists := db.GetEntity(key)

	if !exists {
		return reply.MakeStatusReply("none")
	}

	switch entity.Data.(type) {

	case []byte:
		reply.MakeStatusReply("string")
		//TODO: 自己实现 set list之类的
		//case []byte:
		//	reply.MakeStatusReply("set")

	}

	return reply.UnknownErrReply{}

}

//rename k1 k2
func execRename(db *DB, args [][]byte) resp.Reply {
	src := string(args[0])
	dest := string(args[1])

	entity, exists := db.GetEntity(src)

	if !exists {
		reply.MakeStandardErrReply("no such key")
	}

	db.PutEntity(dest, entity)
	db.Remove(src)

	return reply.MakeOkReply()
}

//renameNX k1 k2
func execRenameNx(db *DB, args [][]byte) resp.Reply {

	src := string(args[0])
	dest := string(args[1])

	_, ok := db.GetEntity(dest)

	if ok {
		return reply.MakeIntReply(0)
	}

	entity, exists := db.GetEntity(src)

	if !exists {
		reply.MakeStandardErrReply("no such key")
	}

	db.PutEntity(dest, entity)
	db.Remove(src)

	return reply.MakeIntReply(1)
}

//KEYS *
func execKeys(db *DB, args [][]byte) resp.Reply {
	pattern := wildcard.CompilePattern(string(args[0]))

	res := make([][]byte, 0)

	db.data.ForEach(func(key string, val interface{}) bool {

		isMatch := pattern.IsMatch(key)

		if isMatch {
			res = append(res, []byte(key))
		}

		return true
	})

	return reply.MakeMultiBulkReply(res)
}

func init() {
	RegisterCommand("DEL", execDel, -2)
	RegisterCommand("EXISTS", execExists, -2)
	RegisterCommand("flushdb", execFlush, 1)
	RegisterCommand("type", execType, 2)
	RegisterCommand("rename", execRename, 3)
	RegisterCommand("renamenx", execRenameNx, 3)
	RegisterCommand("keys", execKeys, 2) // keys *
}
