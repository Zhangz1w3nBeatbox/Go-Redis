package database

import (
	"Go-Redis/interface/database"
	"Go-Redis/interface/resp"
	"Go-Redis/resp/reply"
)

func (db *DB) getAsString(key string) ([]byte, reply.ErrorReply) {
	entity, ok := db.GetEntity(key)
	if !ok {
		return nil, nil
	}

	bytes, ok := entity.Data.([]byte)

	if !ok {
		return nil, &reply.WrongTypeErrReply{}
	}

	return bytes, nil
}

//GET k1
func execGet(db *DB, args [][]byte) resp.Reply {
	key := string(args[0])

	bytes, err := db.getAsString(key)

	if err != nil {
		return err
	}

	if bytes == nil {
		return &reply.NullBulkReply{}
	}

	return reply.MakeBulkReply(bytes)
}

//SET k1 v1
func execSet(db *DB, args [][]byte) resp.Reply {
	key := string(args[0])
	val := args[1]

	entity := &database.DataEntity{
		Data: val,
	}

	db.PutEntity(key, entity)

	return reply.MakeOkReply()
}

//SETNX
func execSetNx(db *DB, args [][]byte) resp.Reply {
	key := string(args[0])
	val := string(args[1])

	entity := &database.DataEntity{
		Data: val,
	}

	code := db.PutIfAbsent(key, entity)

	return reply.MakeIntReply(int64(code))
}

//GET SET
func execGetSet(db *DB, args [][]byte) resp.Reply {
	key := string(args[0])
	val := args[1]
	entity, exists := db.GetEntity(key)
	db.PutEntity(key, &database.DataEntity{Data: val})

	if !exists {
		return reply.MakeNullBulkReply()
	}

	bytes := entity.Data.([]byte)

	return reply.MakeBulkReply(bytes)
}

//STRLEN
func execStrLen(db *DB, args [][]byte) resp.Reply {
	key := string(args[0])
	entity, exists := db.GetEntity(key)
	if !exists {
		return reply.MakeNullBulkReply()
	}

	bytes := entity.Data.([]byte)

	return reply.MakeIntReply(int64(len(bytes)))
}

func init() {
	RegisterCommand("get", execGet, 2)
	RegisterCommand("set", execSet, 3)
	RegisterCommand("getset", execGetSet, 3)
	RegisterCommand("setnx", execSetNx, 3)
	RegisterCommand("strlen", execStrLen, 2)
}
