package database

import (
	"Go-Redis/config"
	"Go-Redis/interface/resp"
	"Go-Redis/lib/logger"
	"Go-Redis/resp/reply"
	"strconv"
	"strings"
)

type Database struct {
	dbSet []*DB
}

func NewDataBase() *Database {

	database := &Database{}

	databasesNum := config.Properties.Databases

	if databasesNum == 0 {
		databasesNum = 16
	}

	dbs := make([]*DB, databasesNum)

	//新建16个数据库
	database.dbSet = dbs

	// 初始化16个数据库
	for i := range database.dbSet {
		db := MakeDB()
		db.index = i
		database.dbSet[i] = db
	}

	return database
}

//执行用户指令

func (database *Database) Exec(client resp.Connection, args [][]byte) resp.Reply {

	defer func() {
		if err := recover(); err != nil {
			logger.Error(err)
		}
	}()

	//处理select指令
	//select 转小写
	cmdName := strings.ToLower(string(args[0]))

	if cmdName == "select" {
		if len(args) != 2 {
			return reply.MakeArgsNumErrReply(cmdName)
		}
		return execSelect(client, database, args[1:])
	}

	//其他指令

	dbIndex := client.GetDBIndex()

	db := database.dbSet[dbIndex]
	// 分数据库的exec方法
	return db.Exec(client, args)
}

func (database *Database) AfterClientClose(c resp.Connection) {

}

func (database *Database) Close() {

}

// select 1
func execSelect(c resp.Connection, db *Database, args [][]byte) resp.Reply {
	dbIdx, err := strconv.Atoi(string(args[0]))

	if err != nil {
		return reply.MakeStandardErrReply("ERR invalid DB index!")
	}

	if dbIdx >= len(db.dbSet) {
		return reply.MakeStandardErrReply("ERR DB index is out of range!")
	}

	c.SelectDB(dbIdx)

	return reply.MakeOkReply()
}
