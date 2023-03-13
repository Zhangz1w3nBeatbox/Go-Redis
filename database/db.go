package database

import (
	"Go-Redis/datastruct/dict"
	"Go-Redis/interface/database"
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
	// 去command表中寻找 对应的指令
	cmd, ok := cmdTable[cmdName]

	if !ok {
		return reply.MakeStandardErrReply("ERR unknown command" + cmdName)
	}

	// 参数个数校验
	if !validateArity(cmd.arity, cmdLine) {
		return reply.MakeArgsNumErrReply(cmdName)
	}

	exectorFunc := cmd.exector

	// set k1 v1 => k1 v1
	return exectorFunc(db, cmdLine[1:])
}

func validateArity(arity int, cmdArgs [][]byte) bool {
	argNum := len(cmdArgs)
	if arity >= 0 {
		return argNum == arity
	}
	return argNum >= -arity
}

// GetEntity returns DataEntity bind to given key
func (db *DB) GetEntity(key string) (*database.DataEntity, bool) {

	raw, ok := db.data.Get(key)
	if !ok {
		return nil, false
	}
	entity, _ := raw.(*database.DataEntity)
	return entity, true
}

// PutEntity a DataEntity into DB
func (db *DB) PutEntity(key string, entity *database.DataEntity) int {
	return db.data.Put(key, entity)
}

// PutIfExists edit an existing DataEntity
func (db *DB) PutIfExists(key string, entity *database.DataEntity) int {
	return db.data.PutIfExists(key, entity)
}

// PutIfAbsent insert an DataEntity only if the key not exists
func (db *DB) PutIfAbsent(key string, entity *database.DataEntity) int {
	return db.data.PutIfAbsent(key, entity)
}

func (db *DB) Remove(key string) {
	db.data.Remove(key)
}

func (db *DB) Removes(keys ...string) (deleted int) {
	deleted = 0

	for _, key := range keys {
		_, exists := db.data.Get(key)
		if exists {
			db.Remove(key)
			deleted++
		}
	}

	return deleted
}

func (db *DB) Flush() {
	db.data.Clear()
}
