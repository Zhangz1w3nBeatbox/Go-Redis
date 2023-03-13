package handler

import (
	"Go-Redis/database"
	databaseface "Go-Redis/interface/database"
	"Go-Redis/lib/logger"
	"Go-Redis/lib/sync/atomic"
	"Go-Redis/resp/connection"
	"Go-Redis/resp/parser"
	"Go-Redis/resp/reply"
	"context"
	"io"
	"net"
	"strings"
	"sync"
)

/*
 * A tcp.RespHandler implements redis protocol
 */

var (
	unknownErrReplyBytes = []byte("-ERR unknown\r\n")
)

// RespHandler implements tcp.Handler and serves as a redis handler
type RespHandler struct {
	activeConn sync.Map // *client -> placeholder
	db         databaseface.Database
	closing    atomic.Boolean // refusing new client and new request
}

func (h *RespHandler) Handler(ctx context.Context, conn net.Conn) {
	if h.closing.Get() {
		// closing handler refuse new connection
		_ = conn.Close()
	}

	client := connection.NewConn(conn)
	h.activeConn.Store(client, 1)

	//ch := parser.ParseStream(conn)
	channel := parser.ParseStream(conn)

	for payload := range channel {
		if payload.Error != nil {
			if payload.Error == io.EOF ||
				payload.Error == io.ErrUnexpectedEOF ||
				strings.Contains(payload.Error.Error(), "use of closed network connection") {
				// connection closed
				h.closeClient(client)
				logger.Info("connection closed: " + client.RemoteAddr().String())
				return
			}
			// protocol err
			errReply := reply.MakeStandardErrReply(payload.Error.Error())
			err := client.Write(errReply.ToBytes())
			if err != nil {
				h.closeClient(client)
				logger.Info("connection closed: " + client.RemoteAddr().String())
				return
			}
			continue
		}
		if payload.Data == nil {
			logger.Error("empty payload")
			continue
		}
		r, ok := payload.Data.(*reply.MultiBulkReply)
		if !ok {
			logger.Error("require multi bulk reply")
			continue
		}

		//
		result := h.db.Exec(client, r.Args)

		if result != nil {
			_ = client.Write(result.ToBytes())
		} else {
			_ = client.Write(unknownErrReplyBytes)
		}
	}
}

// MakeHandler creates a RespHandler instance
func MakeHandler() *RespHandler {
	var db databaseface.Database
	db = database.NewDataBase() // 使用正在的database
	return &RespHandler{
		db: db,
	}
}

func (h *RespHandler) closeClient(client *connection.Connection) {
	_ = client.Close()
	h.db.AfterClientClose(client)
	h.activeConn.Delete(client)
}

// Close stops handler
func (h *RespHandler) Close() error {
	logger.Info("handler shutting down...")
	h.closing.Set(true)
	// TODO: concurrent wait
	h.activeConn.Range(func(key interface{}, val interface{}) bool {
		client := key.(*connection.Connection)
		_ = client.Close()
		return true
	})
	h.db.Close()
	return nil
}
